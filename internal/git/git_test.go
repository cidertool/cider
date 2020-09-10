package git

import (
	"os/exec"
	"strconv"
	"testing"

	"github.com/cidertool/cider/internal/shell"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func newMockGit(commands ...command) *Git {
	ctx := context.New(config.Project{})
	return newMockGitWithContext(ctx, commands...)
}

func newMockGitWithContext(ctx *context.Context, commands ...command) *Git {
	return &Git{
		Shell: &mockShell{
			Context:      ctx,
			commandIndex: 0,
			commands:     commands,
		},
	}
}

func TestNew(t *testing.T) {
	ctx := context.New(config.Project{})
	client := New(ctx)
	ok := client.IsRepo()
	assert.True(t, ok)
}

func TestSanitizeProcess(t *testing.T) {
	runFunc := func(client *Git) (*shell.CompletedProcess, error) {
		return client.RunInEnv(map[string]string{
			"TEST": "TEST",
		}, "test")
	}

	ctx := context.New(config.Project{})
	ctx.CurrentDirectory = "test"
	client := newMockGitWithContext(
		ctx,
		command{Stdout: "true", Stderr: "false"},
		command{ReturnCode: 1, Stdout: "true", Stderr: "false"},
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
}

func TestExtractRemoteFromConfig(t *testing.T) {
	expected := config.Repo{
		Name:  "cider",
		Owner: "cidertool",
	}

	// happy path
	client := newMockGit(
		command{Stdout: "true"},
		command{Stdout: "git@github.com:cidertool/cider.git"},
	)
	repo, err := client.ExtractRepoFromConfig()
	assert.NoError(t, err)
	assert.Equal(t, expected, repo)

	client = newMockGit(
		command{Stdout: "false"},
	)
	repo, err = client.ExtractRepoFromConfig()
	assert.Error(t, err)
	assert.Empty(t, repo)

	client = newMockGit(
		command{Stdout: "true"},
		command{ReturnCode: 1, Stderr: "no repo"},
	)
	repo, err = client.ExtractRepoFromConfig()
	assert.Error(t, err)
	assert.Empty(t, repo)
}

func TestExtractRepoFromURL(t *testing.T) {
	var repo config.Repo
	expected := config.Repo{
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

type mockShell struct {
	Context      *context.Context
	commandIndex int
	commands     []command
}

type command struct {
	ReturnCode int
	Stdout     string
	Stderr     string
}

func (sh *mockShell) NewCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...) // #nosec
}

func (sh *mockShell) Exec(cmd *exec.Cmd) (*shell.CompletedProcess, error) {
	currentCommand := sh.commands[sh.commandIndex]
	ps := shell.CompletedProcess{
		Name:       cmd.Path,
		Args:       cmd.Args,
		ReturnCode: currentCommand.ReturnCode,
		Stdout:     currentCommand.Stdout,
		Stderr:     currentCommand.Stderr,
	}
	sh.commandIndex++
	var err error
	if currentCommand.ReturnCode != 0 {
		err = &mockShellError{
			Process: ps,
		}
	}
	return &ps, err
}

func (sh *mockShell) Exists(program string) bool {
	return true
}

func (sh *mockShell) CurrentDirectory() string {
	return sh.Context.CurrentDirectory
}

type mockShellError struct {
	Process shell.CompletedProcess
}

func (err *mockShellError) Error() string {
	return strconv.Itoa(err.Process.ReturnCode)
}
