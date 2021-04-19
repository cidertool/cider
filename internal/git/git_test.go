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
	"testing"

	"github.com/cidertool/cider/internal/shell"
	"github.com/cidertool/cider/internal/shell/shelltest"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func newMockGit(t *testing.T, commands ...shelltest.Command) *Git {
	ctx := context.New(config.Project{})

	return newMockGitWithContext(ctx, t, commands...)
}

func newMockGitWithContext(ctx *context.Context, t *testing.T, commands ...shelltest.Command) *Git {
	return &Git{
		Shell: &shelltest.Shell{
			T:        t,
			Context:  ctx,
			Commands: commands,
		},
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	ctx := context.New(config.Project{})
	client := New(ctx)
	ok := client.IsRepo()
	assert.True(t, ok)
}

func TestSanitizeProcess(t *testing.T) {
	t.Parallel()

	runFunc := func(client *Git) (*shell.CompletedProcess, error) {
		return client.RunInEnv(map[string]string{
			"TEST": "TEST",
		}, "test")
	}

	ctx := context.New(config.Project{})
	ctx.CurrentDirectory = "test"
	client := newMockGitWithContext(
		ctx,
		t,
		shelltest.Command{Stdout: "true", Stderr: "false"},
		shelltest.Command{ReturnCode: 1, Stdout: "true", Stderr: "false"},
	)

	// Test out sanitize
	proc, err := runFunc(client)
	assert.NoError(t, err)
	assert.Equal(t, []string{"git", "-c", "log.showSignature=false", "-C", "test", "test"}, proc.Args)
	out, err := client.SanitizeProcess(proc, err)
	assert.Equal(t, "true", out)
	assert.NoError(t, err)

	// Test error sanitize
	proc, err = runFunc(client)
	assert.Error(t, err)
	out, err = client.SanitizeProcess(proc, err)
	assert.Equal(t, "true", out)
	assert.Error(t, err)
	assert.EqualError(t, err, "false")
}

func TestShowRef(t *testing.T) {
	t.Parallel()

	// Selected the initial commit of this repo, because I needed a sha1 hash.
	expected := "eac16d260ebf8af83873c9704169cf40a5501f84"
	client := newMockGit(
		t,
		shelltest.Command{Stdout: expected},
	)
	got, err := client.Show("%H")
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestExtractRemoteFromConfig_Happy(t *testing.T) {
	t.Parallel()

	expected := Repo{
		Name:  "cider",
		Owner: "cidertool",
	}

	client := newMockGit(
		t,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{Stdout: "git@github.com:cidertool/cider.git"},
	)
	repo, err := client.ExtractRepoFromConfig()
	assert.NoError(t, err)
	assert.Equal(t, expected, repo)
}

func TestExtractRemoteFromConfig_ErrIsNotRepo(t *testing.T) {
	t.Parallel()

	client := newMockGit(
		t,
		shelltest.Command{Stdout: "false"},
	)
	repo, err := client.ExtractRepoFromConfig()
	assert.Error(t, err)
	assert.Empty(t, repo)
}

func TestExtractRemoteFromConfig_ErrNoRemoteNamedOrigin(t *testing.T) {
	t.Parallel()

	client := newMockGit(
		t,
		shelltest.Command{Stdout: "true"},
		shelltest.Command{ReturnCode: 1, Stderr: "no repo"},
	)
	repo, err := client.ExtractRepoFromConfig()
	assert.Error(t, err)
	assert.Empty(t, repo)
}

func TestExtractRepoFromURL(t *testing.T) {
	t.Parallel()

	var repo Repo

	expected := Repo{
		Name:  "cider",
		Owner: "cidertool",
	}
	repo = ExtractRepoFromURL("https://github.com/cidertool/cider")
	assert.Equal(t, expected, repo)
	repo = ExtractRepoFromURL("https://github.com/cidertool/cider.git")
	assert.Equal(t, expected, repo)
	repo = ExtractRepoFromURL("ssh://github.com/cidertool/cider.git")
	assert.Equal(t, expected, repo)
	repo = ExtractRepoFromURL("ssh://git@github.com/cidertool/cider.git")
	assert.Equal(t, expected, repo)
	repo = ExtractRepoFromURL("git@github.com:cidertool/cider.git")
	assert.Equal(t, expected, repo)
}
