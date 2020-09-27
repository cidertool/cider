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

func TestGit_SkipGit(t *testing.T) {
	ctx := context.New(config.Project{})
	ctx.Version = "1.0"
	ctx.SkipGit = true
	p := Pipe{}

	err := p.Run(ctx)
	assert.EqualError(t, err, pipe.ErrSkipGitEnabled.Error())
}

func TestGit_Happy_EnvCurrentTag(t *testing.T) {
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

func TestGit_Err_LoginShell(t *testing.T) {
	ctx := context.New(config.Project{})
	p := Pipe{}

	err := p.Run(ctx)
	assert.Error(t, err)
}

func TestGit_Err_NoGit(t *testing.T) {
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
}

func TestGit_Err_NotInRepo(t *testing.T) {
	ctx := context.New(config.Project{})
	p := Pipe{}
	p.client = newMockGitWithContext(ctx,
		shelltest.Command{ReturnCode: 128, Stdout: "fatal"},
	)

	err := p.Run(ctx)
	assert.Error(t, err)
}

func TestGit_Err_BadCommit(t *testing.T) {
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
