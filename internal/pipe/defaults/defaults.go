// Package defaults runs all defaulter pipelines
package defaults

import (
	"github.com/cidertool/cider/internal/defaults"
	"github.com/cidertool/cider/internal/middleware"
	"github.com/cidertool/cider/pkg/context"
)

// Pipe that sets the defaults.
type Pipe struct {
	defaulters []defaults.Defaulter
}

func (Pipe) String() string {
	return "setting defaults"
}

// Run the pipe.
func (p Pipe) Run(ctx *context.Context) error {
	if len(p.defaulters) == 0 {
		p.defaulters = defaults.Defaulters
	}

	for _, defaulter := range p.defaulters {
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
