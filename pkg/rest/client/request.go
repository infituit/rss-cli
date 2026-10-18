package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	requestContextNilError = "request context is nil"
	httpSpanNamePrefix     = "HTTP "
	httpRequestMethodAttr  = "http.request.method"
	urlPathAttr            = "url.path"
	httpResponseStatusAttr = "http.response.status_code"
	httpTracerName         = "rss-cli.http"
	retryDelay             = 20 * time.Millisecond
	defaultAccept          = "application/json"
	defaultContentType     = "application/json"
	formContentType        = "application/x-www-form-urlencoded"
	httpScheme             = "http"
	httpsScheme            = "https"
)

type Request struct {
	client       *Client
	method       string
	path         string
	targetURL    *url.URL
	query        url.Values
	headers      http.Header
	body         []byte
	bodyReader   io.Reader
	context      context.Context
	timeout      time.Duration
	logBody      bool
	logResponse  bool
	queryMaskers map[string]func(string) string
	err          error
}

func (r *Request) WithContext(ctx context.Context) *Request   { r.context = ctx; return r }
func (r *Request) WithTimeout(timeout time.Duration) *Request { r.timeout = timeout; return r }
func (r *Request) WithParam(key, value string) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Add(key, value)
	return r
}

// WithQuery encodes the request's query fields from a declared request value.
func (r *Request) WithQuery(value any) *Request {
	values, err := encodeValues(value, "query")
	if err != nil {
		r.err = err
		return r
	}
	for key, entries := range values {
		for _, entry := range entries {
			r.WithParam(key, entry)
		}
	}
	return r
}
func (r *Request) WithParams(values map[string]string) *Request {
	for k, v := range values {
		r.WithParam(k, v)
	}
	return r
}

func (r *Request) WithHeaders(headers http.Header) *Request {
	for key, values := range headers {
		for _, value := range values {
			r.headers.Add(key, value)
		}
	}
	return r
}

func (r *Request) SubResourcef(path string, args ...any) *Request {
	if r.targetURL != nil {
		r.err = errors.New("target URL cannot be used with sub resource")
		return r
	}
	r.path = fmt.Sprintf(strings.TrimLeft(path, "/"), args...)
	return r
}

func (r *Request) WithURL(target *url.URL) *Request {
	if target == nil {
		r.err = errors.New("target URL is nil")
		return r
	}
	if target.Scheme != httpScheme && target.Scheme != httpsScheme {
		r.err = fmt.Errorf("unsupported target URL scheme %s", target.Scheme)
		return r
	}
	if target.Host == "" {
		r.err = errors.New("target URL host is empty")
		return r
	}
	copyURL := *target
	r.targetURL = &copyURL
	return r
}

func (r *Request) Body(body any) *Request {
	if body == nil {
		r.body = []byte{}
		return r
	}
	data, err := marshalBody(body)
	if err != nil {
		r.err = err
		return r
	}
	r.body = data
	return r
}

func (r *Request) Form(value any) *Request {
	values, err := encodeValues(value, "form")
	if err != nil {
		r.err = err
		return r
	}
	r.body = []byte(values.Encode())
	r.headers.Set("Content-Type", formContentType)
	return r
}

func (r *Request) BodyReader(reader io.Reader) *Request { r.bodyReader = reader; return r }
func (r *Request) EnableLogBody() *Request              { r.logBody = true; return r }
func (r *Request) EnableLogResponse() *Request          { r.logResponse = true; return r }

func (r *Request) Do(ctx context.Context) *Result {
	if ctx == nil {
		return &Result{Err: errors.New(requestContextNilError)}
	}
	if r.err != nil {
		return &Result{Err: r.err}
	}
	requestURL, err := r.buildURL()
	if err != nil {
		return &Result{Err: err}
	}
	var lastErr error
	for attempt := 0; attempt < r.client.capability.MaxRetryCycle; attempt++ {
		result, retry := r.doOnce(ctx, requestURL)
		if !retry {
			return result
		}
		lastErr = result.Err
		time.Sleep(retryDelay)
	}
	return &Result{Err: lastErr}
}

func (r *Request) buildURL() (*url.URL, error) {
	var target url.URL
	if r.targetURL != nil {
		target = *r.targetURL
	} else {
		base, err := url.Parse(r.client.baseURL)
		if err != nil {
			return nil, err
		}
		target = *base
		target.Path, _ = url.JoinPath(target.Path, r.path)
	}
	query := target.Query()
	for key, values := range r.query {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	if r.timeout > 0 {
		query.Set("timeout", r.timeout.String())
	}
	target.RawQuery = query.Encode()
	return &target, nil
}

func (r *Request) doOnce(ctx context.Context, target *url.URL) (*Result, bool) {
	ctx, span := otel.Tracer(httpTracerName).Start(ctx, httpSpanNamePrefix+r.method)
	defer span.End()
	span.SetAttributes(attribute.String(httpRequestMethodAttr, r.method), attribute.String(urlPathAttr, target.Path))
	reader := r.bodyReader
	if reader == nil {
		reader = bytes.NewReader(r.body)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, target.String(), reader)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return &Result{Err: err}, false
	}
	req.Header = r.headers.Clone()
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", defaultAccept)
	}
	if req.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", defaultContentType)
	}
	r.client.capability.Logger.Debug("HTTP request", "method", r.method, "url", maskURL(target.String(), r.queryMaskers), "headers", r.maskHeaders(req.Header), "body", r.loggedBody())
	start := time.Now()
	resp, err := r.client.capability.HTTPClient.Do(req)
	duration := time.Since(start)
	if r.client.capability.MetricObserver != nil {
		r.client.capability.MetricObserver.Observe(r.client.capability.Name, r.method, target.Path, responseStatus(resp), duration)
	}
	if duration >= r.client.capability.ToleranceLatencyTime {
		r.client.capability.Logger.Warn("HTTP request exceeded latency threshold", "method", r.method, "url", maskURL(target.String(), r.queryMaskers), "duration", duration)
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.client.capability.Logger.Error("HTTP request failed", "method", r.method, "url", maskURL(target.String(), r.queryMaskers), "error", err)
		if r.method == http.MethodGet && isConnectionReset(err) {
			return &Result{Err: err}, true
		}
		return &Result{Err: err}, false
	}
	span.SetAttributes(attribute.Int(httpResponseStatusAttr, resp.StatusCode))
	r.client.capability.Logger.Debug("HTTP response", "method", r.method, "url", maskURL(target.String(), r.queryMaskers), "status", resp.StatusCode, "body", r.logResponseBody())
	return &Result{FullURL: target.String(), Body: resp.Body, StatusCode: resp.StatusCode, Status: resp.Status, Header: resp.Header, logResponse: r.logResponse, queryMaskers: r.queryMaskers}, false
}

func (r *Request) loggedBody() string {
	if !r.logBody {
		return "hidden"
	}
	return string(r.body)
}

func (r *Request) logResponseBody() string {
	if !r.logResponse {
		return "hidden"
	}
	return "response body is available through Result"
}

func (r *Request) maskHeaders(headers http.Header) http.Header {
	masked := headers.Clone()
	for key, mask := range r.client.capability.HeaderMaskers {
		values, ok := masked[http.CanonicalHeaderKey(key)]
		if !ok {
			continue
		}
		for i, value := range values {
			values[i] = mask(value)
		}
	}
	return masked
}

func responseStatus(resp *http.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode
}

func marshalBody(body any) ([]byte, error) {
	if value, ok := body.(string); ok {
		return []byte(value), nil
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}
	return data, nil
}

func encodeValues(value any, tagName string) (url.Values, error) {
	if value == nil {
		return nil, errors.New("request value is nil")
	}
	reflected := reflect.ValueOf(value)
	for reflected.Kind() == reflect.Pointer {
		if reflected.IsNil() {
			return nil, errors.New("request value is nil")
		}
		reflected = reflected.Elem()
	}
	if reflected.Kind() != reflect.Struct {
		return nil, fmt.Errorf("request value must be a struct, got %s", reflected.Kind())
	}

	values := url.Values{}
	typeInfo := reflected.Type()
	for index := 0; index < reflected.NumField(); index++ {
		field := typeInfo.Field(index)
		name := field.Tag.Get(tagName)
		if name == "" || name == "-" {
			continue
		}
		values.Set(name, fmt.Sprint(reflected.Field(index).Interface()))
	}
	return values, nil
}

func isConnectionReset(err error) bool {
	return errors.Is(err, syscall.ECONNRESET) || errors.Is(err, net.ErrClosed)
}

func (r *Result) Into(value any) error {
	data, err := r.RawData()
	if err != nil {
		return err
	}
	if len(data) == 0 || value == nil {
		return nil
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}
	return nil
}

func (r *Result) RawData() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if r.Body == nil {
		return nil, errors.New("response body is nil")
	}
	data, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP request failed, status(%d), body(%s)", r.StatusCode, limitBody(data))
	}
	return data, nil
}

func (r *Result) RawStream() (io.ReadCloser, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if r.Body == nil {
		return nil, errors.New("response body is nil")
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		data, _ := io.ReadAll(io.LimitReader(r.Body, 1024))
		_ = r.Body.Close()
		return nil, fmt.Errorf("HTTP request failed, status(%d), body(%s)", r.StatusCode, data)
	}
	return r.Body, nil
}

func limitBody(data []byte) []byte {
	if len(data) > 1024 {
		return data[:1024]
	}
	return data
}

func maskURL(raw string, maskers map[string]func(string) string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	query := parsed.Query()
	for key, values := range query {
		if mask, ok := maskers[key]; ok {
			for i, value := range values {
				values[i] = mask(value)
			}
			query[key] = values
		}
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}
