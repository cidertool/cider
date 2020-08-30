package client

import (
	"fmt"
	"time"

	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

type appleClient struct {
	client *asc.Client
}

// NewApple returns a new Client implemented by asc.Client.
func NewApple(ctx *context.Context) (Client, error) {
	token, err := asc.NewTokenConfig(ctx.Credentials.KeyID, ctx.Credentials.IssuerID, time.Minute*20, []byte(ctx.Credentials.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("could not interpret the p8 private key: %w", err)
	}
	c := appleClient{asc.NewClient(token.Client())}
	return c, nil
}
