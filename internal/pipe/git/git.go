package git

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aaronsky/applereleaser/internal/git"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/apex/log"
)

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "getting and validating git state"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	if _, err := exec.LookPath("git"); err != nil {
		return ErrNoGit
	}
	info, err := getInfo(ctx)
	if err != nil {
		return err
	}
	ctx.Git = info
	log.Infof("releasing %s, commit %s", info.CurrentTag, info.Commit)
	ctx.Version = strings.TrimPrefix(ctx.Git.CurrentTag, "v")
	return validate(ctx)
}

func getInfo(ctx *context.Context) (context.GitInfo, error) {
	if !git.IsRepo(ctx) {
		return context.GitInfo{}, ErrNotRepository
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
	tag, err := getTag(ctx)
	if err != nil {
		return context.GitInfo{
			Commit:      full,
			FullCommit:  full,
			ShortCommit: short,
			CommitDate:  date,
			URL:         url,
			CurrentTag:  "v0.0.0",
		}, ErrNoTag
	}
	return context.GitInfo{
		CurrentTag:  tag,
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
		return ErrDirty{status: out}
	}
	_, err = git.Clean(git.Run(ctx, "describe", "--exact-match", "--tags", "--match", ctx.Git.CurrentTag))
	if err != nil {
		return ErrWrongRef{
			commit: ctx.Git.Commit,
			tag:    ctx.Git.CurrentTag,
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
