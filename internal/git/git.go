/**
Copyright (C) 2020 Aaron Sky.

This file is part of Cider, a tool for automating submission
of apps to Apple's App Stores.

Cider is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cider is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cider.  If not, see <http://www.gnu.org/licenses/>.
*/

// Package git provides an integration with the git command
package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cidertool/cider/internal/shell"
	"github.com/cidertool/cider/pkg/context"
)

// Repo represents any kind of repo (github, gitlab, etc).
type Repo struct {
	Owner string
	Name  string
}

// Git wraps a shell.Shell provider to provide an interface over
// the git program in the PATH.
type Git struct {
	shell.Shell
}

// New constructs a new Git instance with a default shell provider
// based on the login shell.
func New(ctx *context.Context) *Git {
	return &Git{
		Shell: shell.New(ctx),
	}
}

// Run runs a Git command and returns its output or errors.
func (git *Git) Run(args ...string) (*shell.CompletedProcess, error) {
	return git.RunInEnv(nil, args...)
}

// RunInEnv runs a Git command with the specified env vars and returns its output or errors.
func (git *Git) RunInEnv(env map[string]string, args ...string) (*shell.CompletedProcess, error) {
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}

	cd := git.CurrentDirectory()
	if cd != "" && cd != "." {
		extraArgs = append(extraArgs, "-C", filepath.Clean(cd))
	}

	args = append(extraArgs, args...)

	var cmd = git.Shell.NewCommand("git", args...)

	if env != nil {
		cmd.Env = []string{}
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	return git.Shell.Exec(cmd)
}

// IsRepo returns true if current folder is a Git repository.
func (git *Git) IsRepo() bool {
	proc, err := git.Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(proc.Stdout) == "true"
}

// SanitizedError wraps the contents of a sanitized Git error.
type SanitizedError struct {
	stderr string
}

func (e SanitizedError) Error() string {
	return e.stderr
}

// SanitizeProcess cleans up the output.
func (git *Git) SanitizeProcess(proc *shell.CompletedProcess, err error) (string, error) {
	var out string

	if proc != nil {
		firstline := strings.Split(proc.Stdout, "\n")[0]
		out = strings.ReplaceAll(firstline, "'", "")

		if err != nil {
			err = SanitizedError{strings.TrimSuffix(proc.Stderr, "\n")}
		}
	}

	return out, err
}

// MARK: Helpers

// Show returns the requested information for the commit pointed to by HEAD.
func (git *Git) Show(spec string) (string, error) {
	return git.ShowRef(spec, "HEAD")
}

// ShowRef returns the requested information for the given ref.
func (git *Git) ShowRef(spec, ref string) (string, error) {
	return git.SanitizeProcess(git.Run("show", fmt.Sprintf("--format=%s", spec), ref, "--quiet"))
}

// ExtractRepoFromConfig gets the repo name from the Git config.
func (git *Git) ExtractRepoFromConfig() (result Repo, err error) {
	if !git.IsRepo() {
		return result, ErrNotRepository{git.CurrentDirectory()}
	}

	proc, err := git.Run("config", "--get", "remote.origin.url")
	if err != nil {
		return result, ErrNoRemoteOrigin
	}

	return ExtractRepoFromURL(proc.Stdout), nil
}

// ExtractRepoFromURL gets the repo name from the remote URL.
func ExtractRepoFromURL(s string) Repo {
	// removes the .git suffix and any new lines
	s = strings.NewReplacer(
		".git", "",
		"\n", "",
	).Replace(s)
	// if the URL contains a :, indicating a SSH config,
	// remove all chars until it, including itself
	// on HTTP and HTTPS URLs it will remove the http(s): prefix,
	// which is ok. On SSH URLs the whole user@server will be removed,
	// which is required.
	s = s[strings.LastIndex(s, ":")+1:]
	// split by /, the last to parts should be the owner and name
	ss := strings.Split(s, "/")

	return Repo{
		Owner: ss[len(ss)-2],
		Name:  ss[len(ss)-1],
	}
}
