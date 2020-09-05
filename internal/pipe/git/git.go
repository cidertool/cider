// Package git is a pipe that reads and validates git environment
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aaronsky/applereleaser/internal/git"
	"github.com/aaronsky/applereleaser/internal/pipe"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/apex/log"
)

// NoTag is a constant value representing an absent tag. This is used in the case of a hardcoded version string, or an invalid Git tag value.
const NoTag = "v0.0.0"

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "getting and validating git state"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	if ctx.SkipGit && ctx.Version != "" {
		ctx.Git = context.GitInfo{
			CurrentTag:  NoTag,
			Commit:      "none",
			ShortCommit: "none",
			FullCommit:  "none",
		}
		return pipe.ErrSkipGitEnabled
	}
	if _, err := exec.LookPath("git"); err != nil {
		return git.ErrNoGit
	}
	info, err := getInfo(ctx)
	if err != nil {
		return err
	}
	ctx.Git = info
	if ctx.Version == "" {
		tag, err := getTag(ctx)
		if err != nil {
			return git.ErrNoTag
		}
		ctx.Git.CurrentTag = tag
		ctx.Version = strings.TrimPrefix(tag, "v")
	}
	log.WithFields(log.Fields{
		"version": ctx.Version,
		"commit":  info.Commit,
	}).Infof("releasing")
	return validate(ctx)
}

func getInfo(ctx *context.Context) (context.GitInfo, error) {
	if !git.IsRepo(ctx) {
		return context.GitInfo{}, git.ErrNotRepository{Dir: ctx.CurrentDirectory}
	}
	return getGitInfo(ctx)
}

func getGitInfo(ctx *context.Context) (context.GitInfo, error) {
	short, err := getShortCommit(ctx)
	if err != nil {
		return context.GitInfo{}, fmt.Errorf("couldn't get current commit: %w", err)
	}
	full, err := getFullCommit(ctx)
	if err != nil {
		return context.GitInfo{}, fmt.Errorf("couldn't get current commit: %w", err)
	}
	date, err := getCommitDate(ctx)
	if err != nil {
		return context.GitInfo{}, fmt.Errorf("couldn't get commit date: %w", err)
	}
	url, err := getURL(ctx)
	if err != nil {
		return context.GitInfo{}, fmt.Errorf("couldn't get remote URL: %w", err)
	}
	return context.GitInfo{
		CurrentTag:  NoTag,
		Commit:      full,
		FullCommit:  full,
		ShortCommit: short,
		CommitDate:  date,
		URL:         url,
	}, nil
}

func validate(ctx *context.Context) error {
	out, err := git.Run(ctx, "status", "--porcelain")
	if strings.TrimSpace(out) != "" || err != nil {
		return git.ErrDirty{Status: out}
	}
	if ctx.Git.CurrentTag != NoTag {
		_, err = git.Clean(git.Run(ctx, "describe", "--exact-match", "--tags", "--match", ctx.Git.CurrentTag))
		if err != nil {
			return git.ErrWrongRef{
				Commit: ctx.Git.Commit,
				Tag:    ctx.Git.CurrentTag,
			}
		}
	}
	return nil
}

func getCommitDate(ctx *context.Context) (time.Time, error) {
	ct, err := git.Clean(git.Run(ctx, "show", "--format='%ct'", "HEAD", "--quiet"))
	if err != nil {
		return time.Time{}, err
	}
	if ct == "" {
		return time.Time{}, nil
	}
	i, err := strconv.ParseInt(ct, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(i, 0).UTC()
	return t, nil
}

func getShortCommit(ctx *context.Context) (string, error) {
	return git.Clean(git.Run(ctx, "show", "--format='%h'", "HEAD", "--quiet"))
}

func getFullCommit(ctx *context.Context) (string, error) {
	return git.Clean(git.Run(ctx, "show", "--format='%H'", "HEAD", "--quiet"))
}

func getTag(ctx *context.Context) (string, error) {
	if tag := os.Getenv("APPLERELEASER_CURRENT_TAG"); tag != "" {
		return tag, nil
	}

	return git.Clean(git.Run(ctx, "describe", "--tags", "--abbrev=0"))
}

func getURL(ctx *context.Context) (string, error) {
	return git.Clean(git.Run(ctx, "ls-remote", "--get-url"))
}
