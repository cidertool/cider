// Package defaults runs all defaulter pipelines
package defaults

import (
	"github.com/aaronsky/applereleaser/internal/defaults"
	"github.com/aaronsky/applereleaser/internal/middleware"
	"github.com/aaronsky/applereleaser/pkg/context"
)

// Pipe that sets the defaults.
type Pipe struct{}

func (Pipe) String() string {
	return "setting defaults"
}

// Run the pipe.
func (Pipe) Run(ctx *context.Context) error {
	for _, defaulter := range defaults.Defaulters {
		if err := middleware.Logging(
			defaulter.String(),
			middleware.ErrHandler(defaulter.Default),
			middleware.ExtraPadding,
		)(ctx); err != nil {
			return err
		}
	}
	return nil
}
