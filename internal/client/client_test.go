package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
)

type mockCredentials struct {
	server *httptest.Server
}

func newMockContext(status int, rawResponse string) *context.Context {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintln(w, rawResponse)
	}))
	ctx := context.New(config.Project{
		Name: "TEST",
	})
	ctx.Credentials = &mockCredentials{
		server: server,
	}
	return ctx
}

func (c *mockCredentials) Close() {
	c.server.Close()
}

func (c *mockCredentials) Client() *http.Client {
	return c.server.Client()
}

// func TestGetApp(t *testing.T) {
// 	ctx := newMockContext(http.StatusOK, ``)
// 	client := New(ctx)
// 	app, err := client.GetAppForBundleID(ctx, "com.app.bundleid")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, app)
// }
