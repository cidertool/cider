/**
Copyright (C) 2020 Aaron Sky.

This file is part of Cider, a tool for automating submission
of apps to Apple's App Stores.

Cider is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cider is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cider.  If not, see <http://www.gnu.org/licenses/>.
*/

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"

	"github.com/apex/log"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

type response struct {
	StatusCode  int
	RawResponse string
	Response    interface{}
}

type testContext struct {
	Context              *context.Context
	Responses            []response
	CurrentResponseIndex int
	server               *httptest.Server
}

type mockCredentials struct {
	url    string
	client *http.Client
}

type mockTransport struct {
	URL       *url.URL
	Transport http.RoundTripper
}

func newTestContext(resp ...response) (*testContext, Client) {
	ctx := testContext{}
	ctx.Context = context.New(config.Project{})
	server := httptest.NewServer(&ctx)
	ctx.Context.Credentials = &mockCredentials{
		url:    server.URL,
		client: server.Client(),
	}
	ctx.Responses = resp
	ctx.server = server
	client := New(ctx.Context)

	return &ctx, client
}

func (c *testContext) Close() {
	c.server.Close()
}

func (c *testContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.CurrentResponseIndex >= len(c.Responses) {
		log.WithFields(log.Fields{
			"currentResponseIndex": c.CurrentResponseIndex,
			"responsesCount":       len(c.Responses),
			"route":                r.URL.Path,
		}).Fatal("index out of bounds")
	}

	resp := c.Responses[c.CurrentResponseIndex]

	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}

	w.WriteHeader(resp.StatusCode)

	var body = resp.RawResponse

	if resp.Response != nil {
		b, err := json.Marshal(resp.Response)
		if err == nil {
			body = string(b)
		}
	}

	if body == "" {
		body = `{}`
	}

	fmt.Fprintln(w, body)
	c.CurrentResponseIndex++
}

func (c *testContext) SetResponses(resp ...response) {
	c.Responses = resp
	c.CurrentResponseIndex = 0
}

func (c *mockCredentials) Client() *http.Client {
	url, _ := url.Parse(c.url)
	c.client.Transport = &mockTransport{URL: url, Transport: c.client.Transport}

	return c.client
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newURL := *t.URL
	newURL.Path = path.Join(newURL.Path, req.URL.Path)
	req.URL = &newURL

	var transport http.RoundTripper
	if t.Transport == nil {
		transport = http.DefaultTransport
	} else {
		transport = t.Transport
	}

	return transport.RoundTrip(req)
}
