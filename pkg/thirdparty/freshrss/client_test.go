package freshrss

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAPIKeyFromCredentials(t *testing.T) {
	got := apiKeyFromCredentials("kevin", "freshrss")
	want := "4a6911fb47a87a77f4de285f4fac856d"
	if got != want {
		t.Fatalf("apiKeyFromCredentials() = %q, want %q", got, want)
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

func TestFeverReadRequests(t *testing.T) {
	tests := []struct {
		name     string
		call     func(*cli) error
		query    string
		body     string
		response string
	}{
		{
			name: "groups",
			call: func(c *cli) error {
				_, err := c.getGroups(context.Background(), &groupsRequest{
					API:    feverAPIQueryValue,
					Groups: feverAPIQueryValue,
				})
				return err
			},
			query:    "groups",
			response: `{"groups":[]}`,
		},
		{
			name: "feeds",
			call: func(c *cli) error {
				_, err := c.getFeeds(context.Background(), &feedsRequest{
					API:   feverAPIQueryValue,
					Feeds: feverAPIQueryValue,
				})
				return err
			},
			query:    "feeds",
			response: `{"feeds":[]}`,
		},
		{
			name: "favicons",
			call: func(c *cli) error {
				_, err := c.getFavicons(context.Background(), &faviconsRequest{
					API:      feverAPIQueryValue,
					Favicons: feverAPIQueryValue,
				})
				return err
			},
			query:    "favicons",
			response: `{"favicons":[]}`,
		},
		{
			name: "items",
			call: func(c *cli) error {
				_, err := c.getItems(context.Background(), &itemsRequest{API: feverAPIQueryValue, Items: feverAPIQueryValue, SinceID: "10", MaxID: "2", WithIDs: "1,2"})
				return err
			},
			query:    "items",
			response: `{"items":[],"total_items":0}`,
		},
		{
			name: "unread item IDs",
			call: func(c *cli) error {
				_, err := c.getUnreadItemIDs(context.Background(), &itemIDsRequest{
					API:           feverAPIQueryValue,
					UnreadItemIDs: feverAPIQueryValue,
				})
				return err
			},
			query:    "unread_item_ids",
			response: `{"unread_item_ids":"1,2"}`,
		},
		{
			name: "saved item IDs",
			call: func(c *cli) error {
				_, err := c.getSavedItemIDs(context.Background(), &itemIDsRequest{
					API:          feverAPIQueryValue,
					SavedItemIDs: feverAPIQueryValue,
				})
				return err
			},
			query:    "saved_item_ids",
			response: `{"saved_item_ids":"1,2"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpClient := roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodPost {
					t.Fatalf("method = %s", req.Method)
				}
				values := req.URL.Query()
				if _, ok := values[test.query]; !ok {
					t.Fatalf("query = %s, missing %s", req.URL.RawQuery, test.query)
				}
				if values.Get("api") != "" {
					t.Fatalf("api = %q", values.Get("api"))
				}
				if test.name == "items" {
					for key, want := range map[string]string{"since_id": "10", "max_id": "2", "with_ids": "1,2"} {
						if values.Get(key) != want {
							t.Fatalf("%s = %q, want %q", key, values.Get(key), want)
						}
					}
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				form, err := url.ParseQuery(string(body))
				if err != nil {
					t.Fatal(err)
				}
				if form.Get("api_key") != "key" {
					t.Fatalf("body = %q", body)
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(test.response))}, nil
			})
			freshRSS, err := newClient(Config{BaseURL: "https://example.com", APIKey: "key"}, httpClient)
			if err != nil {
				t.Fatal(err)
			}
			if err := test.call(freshRSS); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMarkItemPostsFieldsAndParsesResponse(t *testing.T) {
	httpClient := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatal(err)
		}
		fields := map[string]string{
			"api_key": "key",
			"mark":    "item",
			"as":      "read",
			"id":      "42",
		}
		for key, want := range fields {
			if form.Get(key) != want {
				t.Fatalf("%s = %q, want %q", key, form.Get(key), want)
			}
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"unread_item_ids":"1","saved_item_ids":"2"}`))}, nil
	})
	freshRSS, err := newClient(Config{BaseURL: "https://example.com", APIKey: "key"}, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	response, err := freshRSS.markItem(context.Background(), &markItemRequest{As: readValue, ID: "42"})
	if err != nil {
		t.Fatal(err)
	}
	if response.UnreadItemIDs != "1" || response.SavedItemIDs != "2" {
		t.Fatalf("response = %+v", response)
	}
}

func TestMarkFeedAndGroupPostsFields(t *testing.T) {
	for _, test := range []struct {
		name string
		call func(*cli) error
		mark string
	}{
		{
			name: "feed",
			call: func(c *cli) error {
				_, err := c.markFeed(
					context.Background(),
					&markFeedRequest{ID: "7", Before: "123"},
				)
				return err
			},
			mark: markFeedValue,
		},
		{
			name: "group",
			call: func(c *cli) error {
				_, err := c.markGroup(
					context.Background(),
					&markGroupRequest{ID: "8", Before: "123"},
				)
				return err
			},
			mark: markGroupValue,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			httpClient := roundTripFunc(func(req *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				form, err := url.ParseQuery(string(body))
				if err != nil {
					t.Fatal(err)
				}
				fields := map[string]string{
					"api_key": "key",
					"mark":    test.mark,
					"as":      readValue,
					"before":  "123",
				}
				for key, want := range fields {
					if form.Get(key) != want {
						t.Fatalf("%s = %q, want %q", key, form.Get(key), want)
					}
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
			})
			freshRSS, err := newClient(Config{BaseURL: "https://example.com", APIKey: "key"}, httpClient)
			if err != nil {
				t.Fatal(err)
			}
			if err := test.call(freshRSS); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUnreadRecentlyReadPostsAction(t *testing.T) {
	httpClient := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatal(err)
		}
		if form.Get("api_key") != "key" || form.Get("unread_recently_read") != "1" {
			t.Fatalf("body = %q", body)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})
	freshRSS, err := newClient(Config{BaseURL: "https://example.com", APIKey: "key"}, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := freshRSS.unreadRecentlyRead(
		context.Background(),
		&unreadRecentlyReadRequest{
			API:                feverAPIQueryValue,
			UnreadRecentlyRead: "1",
		},
	); err != nil {
		t.Fatal(err)
	}
}

func TestFeverRequestErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		statusCode int
		response   string
		want       string
	}{
		{
			name:       "non-2xx",
			statusCode: http.StatusBadGateway,
			response:   "bad gateway",
			want:       "fever groups request failed",
		},
		{
			name:       "invalid JSON",
			statusCode: http.StatusOK,
			response:   "{",
			want:       "fever groups request failed",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			httpClient := roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: test.statusCode, Body: io.NopCloser(strings.NewReader(test.response))}, nil
			})
			freshRSS, err := newClient(Config{BaseURL: "https://example.com", APIKey: "key"}, httpClient)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := freshRSS.getGroups(
				context.Background(),
				&groupsRequest{
					API:    feverAPIQueryValue,
					Groups: feverAPIQueryValue,
				},
			); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want substring %q", err, test.want)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return fn(req)
}
