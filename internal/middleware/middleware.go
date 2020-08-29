package middleware

import "github.com/aaronsky/applereleaser/pkg/context"

type Action func(ctx *context.Context) error
