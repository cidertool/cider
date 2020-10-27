// Package semver is a pipe that parses a version string into semver components
package semver

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/cidertool/cider/pkg/context"
)

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "parsing version"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	sv, err := semver.NewVersion(ctx.Version)
	if err != nil {
		return fmt.Errorf("failed to parse tag %s as semver: %w", ctx.Version, err)
	}

	ctx.Semver = context.Semver{
		Major:      sv.Major(),
		Minor:      sv.Minor(),
		Patch:      sv.Patch(),
		Prerelease: sv.Prerelease(),
		RawVersion: sv.Original(),
	}

	return nil
}
