// Package git is a pipe that reads and validates git environment
package git

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/cidertool/cider/internal/git"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/context"
)

// NoTag is a constant value representing an absent tag. This is used in the case of a hardcoded version string, or an invalid Git tag value.
const NoTag = "v0.0.0"

// Pipe is a global hook pipe.
type Pipe struct {
	client *git.Git
}

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

	client := p.client
	if client == nil {
		client = git.New(ctx)
	}

	if ok := client.Exists("git"); !ok {
		return git.ErrNoGit
	}

	info, err := getInfo(client)
	if err != nil {
		return err
	}
	ctx.Git = info
	log.WithFields(log.Fields{
		"commit": info.Commit,
		"tag":    info.CurrentTag,
		"date":   info.CommitDate.String(),
		"url":    info.URL,
	}).Debug("git info")

	if ctx.Version == "" {
		tag, err := getTag(client)
		if err != nil {
			return git.ErrNoTag
		}
		ctx.Git.CurrentTag = tag
		ctx.Version = strings.TrimPrefix(tag, "v")
	}

	log.WithFields(log.Fields{
		"version": ctx.Version,
		"commit":  info.Commit,
		"workdir": ctx.CurrentDirectory,
	}).Info("releasing")

	return validate(ctx, client)
}

func getInfo(client *git.Git) (context.GitInfo, error) {
	if !client.IsRepo() {
		return context.GitInfo{}, git.ErrNotRepository{
			Dir: client.CurrentDirectory(),
		}
	}
	return getGitInfo(client)
}

func getGitInfo(client *git.Git) (context.GitInfo, error) {
	short, err := getShortCommit(client)
	if err != nil {
		return context.GitInfo{}, fmt.Errorf("couldn't get current commit: %w", err)
	}
	full, err := getFullCommit(client)
	if err != nil {
		return context.GitInfo{}, fmt.Errorf("couldn't get current commit: %w", err)
	}
	date, err := getCommitDate(client)
	if err != nil {
		return context.GitInfo{}, fmt.Errorf("couldn't get commit date: %w", err)
	}
	url, err := getURL(client)
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

func validate(ctx *context.Context, client *git.Git) error {
	proc, err := client.Run("status", "--porcelain")
	if strings.TrimSpace(proc.Stdout) != "" || err != nil {
		return git.ErrDirty{Status: proc.Stdout}
	}
	if ctx.Git.CurrentTag != NoTag {
		_, err := client.SanitizeProcess(client.Run("describe", "--exact-match", "--tags", "--match", ctx.Git.CurrentTag))
		if err != nil {
			return git.ErrWrongRef{
				Commit: ctx.Git.Commit,
				Tag:    ctx.Git.CurrentTag,
			}
		}
	}
	return nil
}

func getCommitDate(client *git.Git) (time.Time, error) {
	commit, err := client.Show("%ct")
	if err != nil {
		return time.Time{}, err
	}
	if commit == "" {
		return time.Time{}, nil
	}
	i, err := strconv.ParseInt(commit, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(i, 0).UTC()
	return t, nil
}

func getShortCommit(client *git.Git) (string, error) {
	return client.Show("%h")
}

func getFullCommit(client *git.Git) (string, error) {
	return client.Show("%H")
}

func getTag(client *git.Git) (string, error) {
	if tag := os.Getenv("CIDER_CURRENT_TAG"); tag != "" {
		return tag, nil
	}
	return client.SanitizeProcess(client.Run("describe", "--tags", "--abbrev=0"))
}

func getURL(client *git.Git) (string, error) {
	return client.SanitizeProcess(client.Run("ls-remote", "--get-url"))
}
