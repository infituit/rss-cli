package freshrss

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/infituit/rss-cli/pkg/rest/client"
)

type Config struct {
	BaseURL     string
	Username    string
	APIPassword string
	APIKey      string
}

type cli struct {
	endpoint   string
	apiKey     string
	restClient client.IClient
}

type authCheckRequest struct {
	API    string `query:"api"`
	APIKey string `form:"api_key"`
}

type authCheckResponse struct {
	APIVersion          int    `json:"api_version"`
	Auth                int    `json:"auth"`
	LastRefreshedOnTime string `json:"last_refreshed_on_time"`
}

func newClient(cfg Config, httpClient client.HTTPClient) (*cli, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, fmt.Errorf("FreshRSS URL is required")
	}

	apiKey := strings.TrimSpace(cfg.APIKey)
	if apiKey == "" {
		if strings.TrimSpace(cfg.Username) == "" || cfg.APIPassword == "" {
			return nil, fmt.Errorf("fever API key or username and API password are required")
		}
		apiKey = apiKeyFromCredentials(cfg.Username, cfg.APIPassword)
	}

	endpoint, err := feverEndpoint(cfg.BaseURL)
	if err != nil {
		return nil, err
	}
	restClient, err := client.NewClient(&client.Capability{
		Name:       providerName,
		BaseURL:    endpoint,
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	return &cli{endpoint: endpoint, apiKey: apiKey, restClient: restClient}, nil
}

func apiKeyFromCredentials(username, apiPassword string) string {
	sum := md5.Sum([]byte(username + credentialSeparator + apiPassword))
	return hex.EncodeToString(sum[:])
}

func (c *cli) authCheck(ctx context.Context, req *authCheckRequest) (*authCheckResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever auth check request is nil")
	}
	response := new(authCheckResponse)
	if err := c.restClient.Post().
		WithContext(ctx).
		WithQuery(req).
		Form(req).
		Do(ctx).
		Into(response); err != nil {
		return nil, fmt.Errorf("fever auth check failed: %w", err)
	}
	return response, nil
}

func feverEndpoint(rawURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("FreshRSS URL must include scheme and host")
	}
	if !strings.HasSuffix(parsed.Path, feverAPIPath) {
		parsed.Path = strings.TrimRight(parsed.Path, "/") + feverAPIPath
	}
	return parsed.String(), nil
}
