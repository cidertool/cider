package client

import "github.com/aaronsky/applereleaser/pkg/context"

// Client interface.
type Client interface{}

// New creates a new client depending on the token type.
func New(ctx *context.Context) (Client, error) {
	return NewApple(ctx)
}
