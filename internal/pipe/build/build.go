package build

import (
	"fmt"

	"github.com/aaronsky/applereleaser/internal/client"
	"github.com/aaronsky/applereleaser/pkg/context"
)

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "choosing processed build"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	c, err := client.New(ctx)
	if err != nil {
		return err
	}
	return doBuild(ctx, c)
}

func doBuild(ctx *context.Context, client client.Client) error {
	fmt.Println("build?")
	return nil
}
