package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRequestFormAndResult(t *testing.T) {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost || req.URL.String() != "https://example.com/api?api=" {
			t.Fatalf("request = %s %s", req.Method, req.URL)
		}
		body, _ := io.ReadAll(req.Body)
		if string(body) != "api_key=secret" {
			t.Fatalf("body = %q", body)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"ok":true}`))}, nil
	})
	rest, err := NewClient(&Capability{BaseURL: "https://example.com/api", HTTPClient: transport})
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		OK bool `json:"ok"`
	}
	if err := rest.Post().WithContext(context.Background()).WithQuery(testRequest{API: ""}).Form(testRequest{APIKey: "secret"}).Do(context.Background()).Into(&result); err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Fatal("response was not decoded")
	}
}

type testRequest struct {
	API    string `query:"api"`
	APIKey string `form:"api_key"`
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) Do(req *http.Request) (*http.Response, error) { return fn(req) }
