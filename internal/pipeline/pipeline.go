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

// Package pipeline stores the top-level pipeline and Piper interface used by most pipes
package pipeline

import (
	"fmt"

	"github.com/cidertool/cider/internal/pipe/defaults"
	"github.com/cidertool/cider/internal/pipe/env"
	"github.com/cidertool/cider/internal/pipe/git"
	"github.com/cidertool/cider/internal/pipe/publish"
	"github.com/cidertool/cider/internal/pipe/semver"
	"github.com/cidertool/cider/internal/pipe/template"
	"github.com/cidertool/cider/pkg/context"
)

// Piper defines a pipe, which can be part of a pipeline (a serie of pipes).
type Piper interface {
	fmt.Stringer

	// Run the pipe
	Run(ctx *context.Context) error
}

// Pipeline contains all pipe implementations in order
// nolint: gochecknoglobals
var Pipeline = []Piper{
	env.Pipe{},
	git.Pipe{},
	semver.Pipe{},
	template.Pipe{},
	defaults.Pipe{},
	publish.Pipe{},
}
