// Package publish is a pipe that runs the testflight or store pipes depending on publish mode
package publish

import (
	"fmt"

	"github.com/cidertool/cider/internal/middleware"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/internal/pipe/store"
	"github.com/cidertool/cider/internal/pipe/testflight"
	"github.com/cidertool/cider/pkg/context"
)

// Pipe that publishes artifacts.
type Pipe struct{}

func (Pipe) String() string {
	return "publishing from app store connect"
}

// Publisher should be implemented by pipes that want to publish artifacts.
type Publisher interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Publish(ctx *context.Context) error
}

// nolint: gochecknoglobals
var publishers = []Publisher{
	testflight.Pipe{},
	store.Pipe{},
}

// Run the pipe.
func (Pipe) Run(ctx *context.Context) error {
	if len(ctx.AppsToRelease) == 0 {
		return pipe.Skip("no apps selected to publish")
	}
	for _, publisher := range publishers {
		if err := middleware.Logging(
			publisher.String(),
			middleware.ErrHandler(publisher.Publish),
			middleware.ExtraPadding,
		)(ctx); err != nil {
			return fmt.Errorf("%s: failed to publish: %w", publisher.String(), err)
		}
	}
	return nil
}
