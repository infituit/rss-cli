package freshrss

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/infituit/rss-cli/pkg/rest/client"
)

type Config struct {
	BaseURL     string
	Username    string
	APIPassword string
	APIKey      string
}

type cli struct {
	apiKey     string
	restClient client.IClient
}

func newClient(cfg Config, httpClient client.HTTPClient) (*cli, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("FreshRSS URL is required")
	}

	apiKey := cfg.APIKey
	if apiKey == "" {
		if cfg.Username == "" || cfg.APIPassword == "" {
			return nil, fmt.Errorf("fever API key or username and API password are required")
		}
		apiKey = apiKeyFromCredentials(cfg.Username, cfg.APIPassword)
	}

	restClient, err := client.NewClient(&client.Capability{
		Name:       providerName,
		BaseURL:    cfg.BaseURL,
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	return &cli{apiKey: apiKey, restClient: restClient}, nil
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
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req).
		Form(req).
		Do(ctx).
		Into(response); err != nil {
		return nil, fmt.Errorf("fever auth check failed: %w", err)
	}
	return response, nil
}

func (c *cli) getGroups(ctx context.Context, req *groupsRequest) (*groupsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever groups request is nil")
	}
	response := new(groupsResponse)
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req).
		Form(&authRequest{APIKey: c.apiKey}).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever groups request failed: %w", err)
	}
	return response, nil
}

func (c *cli) getFeeds(ctx context.Context, req *feedsRequest) (*feedsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever feeds request is nil")
	}
	response := new(feedsResponse)
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req).
		Form(&authRequest{APIKey: c.apiKey}).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever feeds request failed: %w", err)
	}
	return response, nil
}

func (c *cli) getFavicons(ctx context.Context, req *faviconsRequest) (*faviconsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever favicons request is nil")
	}
	response := new(faviconsResponse)
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req).
		Form(&authRequest{APIKey: c.apiKey}).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever favicons request failed: %w", err)
	}
	return response, nil
}

func (c *cli) getItems(ctx context.Context, req *itemsRequest) (*itemsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever items request is nil")
	}
	response := new(itemsResponse)
	request := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req)
	if req.SinceID != "" {
		request.WithParam("since_id", req.SinceID)
	}
	if req.MaxID != "" {
		request.WithParam("max_id", req.MaxID)
	}
	if req.WithIDs != "" {
		request.WithParam("with_ids", req.WithIDs)
	}
	err := request.
		Form(&authRequest{APIKey: c.apiKey}).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever items request failed: %w", err)
	}
	return response, nil
}

func (c *cli) getUnreadItemIDs(
	ctx context.Context,
	req *itemIDsRequest,
) (*itemIDsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever unread item IDs request is nil")
	}
	response := new(itemIDsResponse)
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req).
		Form(&authRequest{APIKey: c.apiKey}).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever unread item IDs request failed: %w", err)
	}
	return response, nil
}

func (c *cli) getSavedItemIDs(
	ctx context.Context,
	req *itemIDsRequest,
) (*itemIDsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever saved item IDs request is nil")
	}
	response := new(itemIDsResponse)
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req).
		Form(&authRequest{APIKey: c.apiKey}).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever saved item IDs request failed: %w", err)
	}
	return response, nil
}

func (c *cli) markItem(ctx context.Context, req *markItemRequest) (*itemIDsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever mark item request is nil")
	}
	response := new(itemIDsResponse)
	request := &apiRequest{API: feverAPIQueryValue}
	form := &markItemForm{
		APIKey: c.apiKey,
		Mark:   markItemValue,
		As:     req.As,
		ID:     req.ID,
	}
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(request).
		Form(form).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever mark item request failed: %w", err)
	}
	return response, nil
}

func (c *cli) markFeed(ctx context.Context, req *markFeedRequest) (*itemIDsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever mark feed request is nil")
	}
	response := new(itemIDsResponse)
	request := &apiRequest{API: feverAPIQueryValue}
	form := &markFeedForm{
		APIKey: c.apiKey,
		Mark:   markFeedValue,
		As:     readValue,
		ID:     req.ID,
		Before: req.Before,
	}
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(request).
		Form(form).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever mark feed request failed: %w", err)
	}
	return response, nil
}

func (c *cli) markGroup(ctx context.Context, req *markGroupRequest) (*itemIDsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("fever mark group request is nil")
	}
	response := new(itemIDsResponse)
	request := &apiRequest{API: feverAPIQueryValue}
	form := &markGroupForm{
		APIKey: c.apiKey,
		Mark:   markGroupValue,
		As:     readValue,
		ID:     req.ID,
		Before: req.Before,
	}
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(request).
		Form(form).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever mark group request failed: %w", err)
	}
	return response, nil
}

func (c *cli) unreadRecentlyRead(
	ctx context.Context,
	req *unreadRecentlyReadRequest,
) (*itemIDsResponse, error) {
	response := new(itemIDsResponse)
	form := &unreadRecentlyReadForm{
		APIKey:             c.apiKey,
		UnreadRecentlyRead: req.UnreadRecentlyRead,
	}
	err := c.restClient.Post().
		SubResourcef(feverAPIPath).
		WithContext(ctx).
		WithQuery(req).
		Form(form).
		Do(ctx).
		Into(response)
	if err != nil {
		return nil, fmt.Errorf("fever unread recently read request failed: %w", err)
	}
	return response, nil
}
