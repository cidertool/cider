package context

import (
	ctx "context"
	"os"
	"strings"
	"time"

	"github.com/aaronsky/applereleaser/pkg/config"
)

// Env is the environment variables.
type Env map[string]string

// Copy returns a copy of the environment.
func (e Env) Copy() Env {
	var out = Env{}
	for k, v := range e {
		out[k] = v
	}
	return out
}

type credentials struct {
	KeyID      string
	IssuerID   string
	PrivateKey string
}

// Context carries along some data through the pipes.
type Context struct {
	ctx.Context
	Config      config.Project
	Env         Env
	Date        time.Time
	Credentials credentials
	SkipPublish bool
}

// New context.
func New(config config.Project) *Context {
	return Wrap(ctx.Background(), config)
}

// NewWithTimeout new context with the given timeout.
func NewWithTimeout(config config.Project, timeout time.Duration) (*Context, ctx.CancelFunc) {
	ctx, cancel := ctx.WithTimeout(ctx.Background(), timeout)
	return Wrap(ctx, config), cancel
}

// Wrap wraps an existing context.
func Wrap(ctx ctx.Context, config config.Project) *Context {
	return &Context{
		Context: ctx,
		Config:  config,
		Env:     splitEnv(os.Environ()),
		Date:    time.Now(),
	}
}

func splitEnv(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		r[p[0]] = p[1]
	}
	return r
}
