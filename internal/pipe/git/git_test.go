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

package git

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cidertool/cider/internal/git"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/internal/shell/shelltest"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestGit_Happy(t *testing.T) {
	t.Parallel()

	expected := context.GitInfo{
		CurrentTag:  "1.0.0",
		Commit:      "abcdef1234567890abcdef1234567890abcdef12",
		ShortCommit: "abcdef12",
		FullCommit:  "abcdef1234567890abcdef1234567890abcdef12",
		CommitDate:  time.Unix(1600914830, 0).UTC(),
		URL:         "git@github.com:cidertool/cider.git",
	}

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{Stdout: strconv.FormatInt(expected.CommitDate.Unix(), 10)},
		shelltest.Command{Stdout: expected.URL},
		shelltest.Command{Stdout: expected.CurrentTag},
		shelltest.Command{Stdout: ""},
		shelltest.Command{Stdout: expected.CurrentTag},
	)

	assert.Equal(t, "getting and validating git state", p.String())

	err := p.Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, ctx.Git)
}

func TestGit_RealGitClient(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	ctx.CurrentDirectory = "TEST"

	p := Pipe{}
	err := p.Run(ctx)
	assert.EqualError(t, err, "the directory at TEST is not a git repository")
}

func TestGit_SkipGit(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	ctx.Version = "1.0"
	ctx.SkipGit = true

	p := Pipe{}

	err := p.Run(ctx)
	assert.EqualError(t, err, pipe.ErrSkipGitEnabled.Error())
}

func TestGit_Happy_EnvCurrentTag(t *testing.T) {
	t.Parallel()

	expected := context.GitInfo{
		CurrentTag:  "1.0.0",
		Commit:      "abcdef1234567890abcdef1234567890abcdef12",
		ShortCommit: "abcdef12",
		FullCommit:  "abcdef1234567890abcdef1234567890abcdef12",
		CommitDate:  time.Unix(1600914830, 0).UTC(),
		URL:         "git@github.com:cidertool/cider.git",
	}

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{Stdout: strconv.FormatInt(expected.CommitDate.Unix(), 10)},
		shelltest.Command{Stdout: expected.URL},
		shelltest.Command{},
		shelltest.Command{Stdout: expected.CurrentTag},
	)

	err := os.Setenv("CIDER_CURRENT_TAG", expected.CurrentTag)
	assert.NoError(t, err)

	err = p.Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, ctx.Git)

	err = os.Unsetenv("CIDER_CURRENT_TAG")
	assert.NoError(t, err)
}

func TestGit_Err_NoGit(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = &git.Git{
		Shell: &shelltest.Shell{
			Context: ctx,
			SupportedPrograms: map[string]bool{
				"git": false,
			},
		},
	}

	err := p.Run(ctx)
	assert.Error(t, err)
	assert.EqualError(t, err, git.ErrNoGit.Error())
}

func TestGit_Err_NotInRepo(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{ReturnCode: 128, Stdout: "fatal"},
	)

	err := p.Run(ctx)
	assert.Error(t, err)
}

func TestGit_Err_BadCommit(t *testing.T) {
	t.Parallel()

	expected := context.GitInfo{
		ShortCommit: "abcdef12",
	}

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{ReturnCode: 128},

		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{ReturnCode: 128},
	)

	err := p.Run(ctx)
	assert.Error(t, err)

	err = p.Run(ctx)
	assert.Error(t, err)
}

func TestGit_Err_BadTime(t *testing.T) {
	t.Parallel()

	expected := context.GitInfo{
		ShortCommit: "abcdef12",
		FullCommit:  "abcdef1234567890abcdef1234567890abcdef12",
	}

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{ReturnCode: 128},

		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{},
		shelltest.Command{ReturnCode: 128},

		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{Stdout: "bad output"},
	)

	err := p.Run(ctx)
	assert.Error(t, err)

	err = p.Run(ctx)
	assert.Error(t, err)

	err = p.Run(ctx)
	assert.Error(t, err)
}

func TestGit_Err_BadTag(t *testing.T) {
	t.Parallel()

	expected := context.GitInfo{
		CurrentTag:  "1.0.0",
		Commit:      "abcdef1234567890abcdef1234567890abcdef12",
		ShortCommit: "abcdef12",
		FullCommit:  "abcdef1234567890abcdef1234567890abcdef12",
		CommitDate:  time.Unix(1600914830, 0).UTC(),
		URL:         "git@github.com:cidertool/cider.git",
	}

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{Stdout: strconv.FormatInt(expected.CommitDate.Unix(), 10)},
		shelltest.Command{Stdout: expected.URL},
		shelltest.Command{ReturnCode: 128, Stderr: "no tags!"},
	)

	err := p.Run(ctx)
	assert.Error(t, err)
}

func TestGit_Err_DirtyWorkingCopy(t *testing.T) {
	t.Parallel()

	expected := context.GitInfo{
		CurrentTag:  "1.0.0",
		Commit:      "abcdef1234567890abcdef1234567890abcdef12",
		ShortCommit: "abcdef12",
		FullCommit:  "abcdef1234567890abcdef1234567890abcdef12",
		CommitDate:  time.Unix(1600914830, 0).UTC(),
		URL:         "git@github.com:cidertool/cider.git",
	}

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{Stdout: strconv.FormatInt(expected.CommitDate.Unix(), 10)},
		shelltest.Command{Stdout: expected.URL},
		shelltest.Command{Stdout: expected.CurrentTag},
		shelltest.Command{Stdout: "some stuff in the working copy"},
	)

	err := p.Run(ctx)
	assert.Error(t, err)
}

func TestGit_Err_InvalidTag(t *testing.T) {
	t.Parallel()

	expected := context.GitInfo{
		CurrentTag:  "1.0.0",
		Commit:      "abcdef1234567890abcdef1234567890abcdef12",
		ShortCommit: "abcdef12",
		FullCommit:  "abcdef1234567890abcdef1234567890abcdef12",
		CommitDate:  time.Unix(1600914830, 0).UTC(),
		URL:         "git@github.com:cidertool/cider.git",
	}

	ctx := context.New(config.Project{})

	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: expected.ShortCommit},
		shelltest.Command{Stdout: expected.FullCommit},
		shelltest.Command{Stdout: strconv.FormatInt(expected.CommitDate.Unix(), 10)},
		shelltest.Command{Stdout: expected.URL},
		shelltest.Command{Stdout: expected.CurrentTag},
		shelltest.Command{Stdout: ""},
		shelltest.Command{ReturnCode: 128, Stdout: "fatal"},
	)

	err := p.Run(ctx)
	assert.Error(t, err)
}

func newMockGitWithContext(ctx *context.Context, commands ...shelltest.Command) *git.Git {
	return &git.Git{
		Shell: &shelltest.Shell{
			Context:  ctx,
			Commands: commands,
		},
	}
}
