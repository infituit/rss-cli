//go:build integration

package testsuite

import (
	"context"
	"testing"

	"github.com/infituit/rss-cli/pkg/thirdparty/freshrss"
	"github.com/infituit/rss-cli/testsuite/support"
)

func TestFreshRSSFeverAuthCheck(t *testing.T) {
	instance := support.RequireFreshRSS(t)
	handler, err := freshrss.NewHandler(freshrss.Config{
		BaseURL:     instance.BaseURL,
		Username:    instance.Username,
		APIPassword: instance.APIPassword,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := handler.AuthCheck(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !response.Authenticated {
		t.Fatalf("auth = %d, want 1", response.Auth)
	}
}
