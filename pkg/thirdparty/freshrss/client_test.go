package freshrss

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAPIKey(t *testing.T) {
	got := APIKey("kevin", "freshrss")
	want := "4a6911fb47a87a77f4de285f4fac856d"
	if got != want {
		t.Fatalf("APIKey() = %q, want %q", got, want)
	}
}

func TestNewClientBuildsFreshRSSEndpoint(t *testing.T) {
	client, err := NewClient(Config{BaseURL: "https://example.com", APIKey: "key"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint != "https://example.com/api/fever.php" {
		t.Fatalf("endpoint = %q", client.endpoint)
	}
}

func TestAuthCheckPostsAPIKey(t *testing.T) {
	httpClient := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %s", req.Method)
		}
		if req.URL.String() != "https://example.com/api/fever.php?api=" {
			t.Fatalf("url = %s", req.URL.String())
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "api_key=key" {
			t.Fatalf("body = %q", body)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"api_version":3,"auth":1,"last_refreshed_on_time":"1520013061"}`)),
		}, nil
	})

	handler, err := NewHandler(Config{BaseURL: "https://example.com", APIKey: "key"}, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := handler.AuthCheck(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.APIVersion != 3 || !resp.Authenticated || resp.LastRefreshedOnTime != "1520013061" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return fn(req)
}
