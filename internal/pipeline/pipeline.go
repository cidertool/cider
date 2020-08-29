package pipeline

import (
	"fmt"

	"github.com/aaronsky/applereleaser/pkg/context"

	"github.com/aaronsky/applereleaser/internal/pipe/changelog"
	"github.com/aaronsky/applereleaser/internal/pipe/env"
	"github.com/aaronsky/applereleaser/internal/pipe/git"
	"github.com/aaronsky/applereleaser/internal/pipe/publish"
	"github.com/aaronsky/applereleaser/internal/pipe/semver"
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
	changelog.Pipe{},
	publish.Pipe{},
}
