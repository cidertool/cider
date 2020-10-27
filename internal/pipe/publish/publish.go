// Package publish is a pipe that runs the testflight or store pipes depending on publish mode
package publish

import (
	"fmt"

	"github.com/cidertool/cider/internal/client"
	"github.com/cidertool/cider/internal/middleware"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/internal/pipe/store"
	"github.com/cidertool/cider/internal/pipe/testflight"
	"github.com/cidertool/cider/pkg/context"
)

// ErrUnsupportedPublishMode happens when an unsupported publish mode is provided to the pipe.
func ErrUnsupportedPublishMode(mode string) error {
	return fmt.Errorf("failed to publish: unsupported publish mode %s", mode)
}

// Pipe that publishes artifacts.
type Pipe struct {
	client client.Client
}

func (Pipe) String() string {
	return "publishing from app store connect"
}

// Publisher should be implemented by pipes that want to publish artifacts.
type Publisher interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Publish(ctx *context.Context) error
}

// Run the pipe.
func (p Pipe) Run(ctx *context.Context) error {
	if len(ctx.AppsToRelease) == 0 {
		return pipe.ErrSkipNoAppsToPublish
	}

	var publisher Publisher

	switch ctx.PublishMode {
	case context.PublishModeTestflight:
		publisher = &testflight.Pipe{Client: p.client}
	case context.PublishModeAppStore:
		publisher = &store.Pipe{Client: p.client}
	default:
		return ErrUnsupportedPublishMode(ctx.PublishMode.String())
	}

	if err := middleware.Logging(
		publisher.String(),
		middleware.ErrHandler(publisher.Publish),
		middleware.ExtraPadding,
	)(ctx); err != nil {
		return fmt.Errorf("%s: failed to publish: %w", publisher.String(), err)
	}

	return nil
}
