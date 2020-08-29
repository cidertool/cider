package semver

import "github.com/aaronsky/applereleaser/pkg/context"

type Pipe struct{}

func (Pipe) String() string {
	return ""
}

func (p Pipe) Run(ctx *context.Context) error {
	return nil
}
