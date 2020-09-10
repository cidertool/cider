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
