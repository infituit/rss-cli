package freshrss

import (
	"context"

	"github.com/infituit/rss-cli/pkg/rest/client"
	"github.com/infituit/rss-cli/pkg/types"
)

// Handler exposes FreshRSS capabilities using rss-cli application types.
type Handler struct {
	cli *cli
}

// IHandler is the provider boundary consumed by command code.
type IHandler interface {
	AuthCheck(context.Context) (types.AuthResult, error)
}

var _ IHandler = (*Handler)(nil)

// NewHandler creates a FreshRSS provider handler.
func NewHandler(cfg Config, httpClient client.HTTPClient) (*Handler, error) {
	cli, err := newClient(cfg, httpClient)
	if err != nil {
		return nil, err
	}
	return &Handler{cli: cli}, nil
}

// AuthCheck checks authentication and converts the Fever response into the
// application-level result type.
func (h *Handler) AuthCheck(ctx context.Context) (types.AuthResult, error) {
	response, err := h.cli.authCheck(ctx, &authCheckRequest{
		API:    feverAPIQueryValue,
		APIKey: h.cli.apiKey,
	})
	if err != nil {
		return types.AuthResult{}, err
	}
	return types.AuthResult{
		APIVersion:          response.APIVersion,
		Authenticated:       response.Auth == authenticatedValue,
		LastRefreshedOnTime: response.LastRefreshedOnTime,
	}, nil
}
