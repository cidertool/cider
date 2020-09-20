package middleware

import (
	"errors"
	"testing"

	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestErrHandler_WrapsError(t *testing.T) {
	ctx := context.New(config.Project{})
	wrapped := ErrHandler(func(ctx *context.Context) error {
		return errors.New("TEST")
	})
	err := wrapped(ctx)
	assert.Error(t, err)
}

func TestErrHandler_IgnoresNoError(t *testing.T) {
	ctx := context.New(config.Project{})
	wrapped := ErrHandler(func(ctx *context.Context) error {
		return nil
	})
	err := wrapped(ctx)
	assert.NoError(t, err)
}

func TestErrHandler_HandlesSkip(t *testing.T) {
	ctx := context.New(config.Project{})
	wrapped := ErrHandler(func(ctx *context.Context) error {
		return pipe.Skip("TEST")
	})
	err := wrapped(ctx)
	assert.NoError(t, err)
}
