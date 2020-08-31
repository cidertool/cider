package context

import (
	ctx "context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/asc-go/asc"
)

// GitInfo includes tags and refs
type GitInfo struct {
	CurrentTag  string
	Commit      string
	ShortCommit string
	FullCommit  string
	CommitDate  time.Time
	URL         string
}

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

// Strings returns the current environment as a list of strings, suitable for
// os executions.
func (e Env) Strings() []string {
	var result = make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}
	return result
}

// Credentials stores credentials used by clients
type Credentials struct {
	*asc.AuthTransport
}

// NewCredentials returns a new store object for App Store Connect credentials
func NewCredentials(keyID, issuerID string, privateKey []byte) (Credentials, error) {
	token, err := asc.NewTokenConfig(keyID, issuerID, time.Minute*20, privateKey)
	if err != nil {
		err = fmt.Errorf("could not interpret the p8 private key: %w", err)
	}
	return Credentials{token}, err
}

// Context carries along some data through the pipes.
type Context struct {
	ctx.Context
	Config           config.Project
	Env              Env
	Date             time.Time
	CurrentDirectory string
	Credentials      Credentials
	SkipPublish      bool
	Git              GitInfo
	Version          string
	Semver           Semver
}

// Semver represents a semantic version.
type Semver struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	RawVersion string
	Prerelease string
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
