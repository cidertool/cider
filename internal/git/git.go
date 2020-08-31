// Package git provides an integration with the git command
package git

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/apex/log"
)

// IsRepo returns true if current folder is a git repository.
func IsRepo(ctx *context.Context) bool {
	out, err := Run(ctx, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// RunEnv runs a git command with the specified env vars and returns its output or errors.
func RunEnv(ctx *context.Context, env map[string]string, args ...string) (string, error) {
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}
	if ctx.CurrentDirectory != "" && ctx.CurrentDirectory != "." {
		extraArgs = append(extraArgs, "-C", ctx.CurrentDirectory)
	}
	args = append(extraArgs, args...)
	var cmd = exec.CommandContext(ctx, "git", args...)

	if env != nil {
		cmd.Env = []string{}
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.WithField("args", args).Debug("running git")
	err := cmd.Run()

	log.WithField("stdout", stdout.String()).
		WithField("stderr", stderr.String()).
		Debug("git result")

	if err != nil {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}

// Run runs a git command and returns its output or errors.
func Run(ctx *context.Context, args ...string) (string, error) {
	return RunEnv(ctx, nil, args...)
}

// Clean the output.
func Clean(output string, err error) (string, error) {
	output = strings.Replace(strings.Split(output, "\n")[0], "'", "", -1)
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
