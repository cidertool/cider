package git

import "github.com/aaronsky/applereleaser/pkg/context"

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return ""
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	return nil
}
