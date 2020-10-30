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

// errUnsupportedPublishMode happens when an unsupported publish mode is provided to the pipe.
type errUnsupportedPublishMode struct {
	mode context.PublishMode
}

func (e errUnsupportedPublishMode) Error() string {
	return fmt.Sprintf("failed to publish: unsupported publish mode %s", e.mode)
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
		return errUnsupportedPublishMode{ctx.PublishMode}
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
