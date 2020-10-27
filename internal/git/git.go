// Package git provides an integration with the git command
package git

import (
	"errors"
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

// SanitizeProcess cleans up the output.
func (git *Git) SanitizeProcess(proc *shell.CompletedProcess, err error) (string, error) {
	var out string

	if proc != nil {
		firstline := strings.Split(proc.Stdout, "\n")[0]
		out = strings.ReplaceAll(firstline, "'", "")

		if err != nil {
			err = errors.New(strings.TrimSuffix(proc.Stderr, "\n"))
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
		return result, errors.New("repository doesn't have an `origin` remote")
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
