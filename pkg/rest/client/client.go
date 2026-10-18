package client

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const defaultMaxRetryCycle = 3

type IClient interface {
	Get() *Request
	Post() *Request
	Put() *Request
	Delete() *Request
	Patch() *Request
	Head() *Request
}

type Client struct {
	capability *Capability
	baseURL    string
}

func NewClient(capability *Capability) (IClient, error) {
	if capability == nil {
		return nil, fmt.Errorf("REST client capability is nil")
	}
	if strings.TrimSpace(capability.BaseURL) == "" {
		return nil, fmt.Errorf("REST client base URL is required")
	}
	if capability.HTTPClient == nil {
		capability.HTTPClient = http.DefaultClient
	}
	if capability.Logger == nil {
		capability.Logger = slog.New(slog.NewTextHandler(nilWriter{}, nil))
	}
	if capability.ToleranceLatencyTime <= 0 {
		capability.ToleranceLatencyTime = ToleranceLatencyTimeDefault
	}
	if capability.MaxRetryCycle <= 0 {
		capability.MaxRetryCycle = defaultMaxRetryCycle
	}
	return &Client{capability: capability, baseURL: strings.TrimRight(capability.BaseURL, "/")}, nil
}

func (c *Client) Get() *Request    { return c.request(http.MethodGet) }
func (c *Client) Post() *Request   { return c.request(http.MethodPost) }
func (c *Client) Put() *Request    { return c.request(http.MethodPut) }
func (c *Client) Delete() *Request { return c.request(http.MethodDelete) }
func (c *Client) Patch() *Request  { return c.request(http.MethodPatch) }
func (c *Client) Head() *Request   { return c.request(http.MethodHead) }

func (c *Client) request(method string) *Request {
	return &Request{
		client:       c,
		method:       method,
		headers:      make(http.Header),
		queryMaskers: c.capability.URLQueryMaskers,
	}
}

type nilWriter struct{}

func (nilWriter) Write(p []byte) (int, error) { return len(p), nil }
