package middleware

import (
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	ctx := context.New(config.Project{})
	wrapped := Logging("TEST", func(ctx *context.Context) error {
		return nil
	}, DefaultInitialPadding)
	err := wrapped(ctx)
	assert.NoError(t, err)
}
