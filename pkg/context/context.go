// Package context manages the state of the pipeline
package context

import (
	ctx "context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cidertool/cider/pkg/config"
)

// PublishMode describes which review destination to publish to.
type PublishMode string

const (
	// PublishModeTestflight publishes to Testflight via beta app review.
	PublishModeTestflight PublishMode = "testflight"
	// PublishModeAppStore publishes for App Store review.
	PublishModeAppStore PublishMode = "appstore"
)

// Context carries along some data through the pipes.
type Context struct {
	ctx.Context
	Config                  config.Project
	RawConfig               config.Project
	Env                     Env
	Date                    time.Time
	Git                     GitInfo
	CurrentDirectory        string
	Credentials             Credentials
	AppsToRelease           []string
	PublishMode             PublishMode
	MaxProcesses            int
	SkipGit                 bool
	SkipUpdatePricing       bool
	SkipUpdateMetadata      bool
	SkipSubmit              bool
	OverrideBetaGroups      bool
	OverrideBetaTesters     bool
	VersionIsInitialRelease bool
	Version                 string
	Build                   string
	Semver                  Semver
}

// Env is the environment variables.
type Env map[string]string

// GitInfo includes tags and refs.
type GitInfo struct {
	CurrentTag  string
	Commit      string
	ShortCommit string
	FullCommit  string
	CommitDate  time.Time
	URL         string
}

// Semver represents a semantic version.
type Semver struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease string
	RawVersion string
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
		Context:      ctx,
		Config:       config,
		RawConfig:    config,
		Env:          splitEnv(os.Environ()),
		Date:         time.Now(),
		MaxProcesses: 1,
	}
}

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

func splitEnv(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		r[p[0]] = p[1]
	}
	return r
}

// String returns the string value of the mode.
func (m PublishMode) String() string {
	return string(m)
}

// Set the mode to an allowed value, or return an error.
func (m *PublishMode) Set(value string) error {
	switch value {
	case "appstore":
		*m = PublishModeAppStore
		return nil
	case "testflight":
		*m = PublishModeTestflight
		return nil
	}
	return fmt.Errorf("invalid value %s for publish mode", value)
}

// Type returns a representation of permissible values.
func (m PublishMode) Type() string {
	return "{appstore,testflight}"
}
