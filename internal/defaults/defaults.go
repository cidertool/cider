// Package defaults is used by pipes to provide default values to a context configuration
package defaults

import (
	"fmt"

	"github.com/aaronsky/applereleaser/pkg/context"
)

// Defaulter can be implemented by a Piper to set default values for its
// configuration.
type Defaulter interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Default(ctx *context.Context) error
}

// Defaulters is the list of defaulters
// nolint: gochecknoglobals
var Defaulters = []Defaulter{}
