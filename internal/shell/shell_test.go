package shell

import (
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	sh := New(context.New(config.Project{}))
	cmd := sh.NewCommand("echo", "dogs")
	ps, err := sh.Exec(cmd)
	assert.NoError(t, err)
	assert.Equal(t, "dogs", ps.Stdout)
}

func TestExec_Error(t *testing.T) {
	sh := New(context.New(config.Project{}))
	cmd := sh.NewCommand("exit", "1")
	ps, err := sh.Exec(cmd)
	assert.Error(t, err)
	assert.Nil(t, ps)
}

func TestShellError_Error(t *testing.T) {
	err := shellError{
		process: CompletedProcess{
			Name:       "dan",
			Args:       []string{"went", "home"},
			ReturnCode: 2,
			Stdout:     "home is missing",
			Stderr:     "failed to go home",
		},
	}
	assert.EqualError(t, &err, "`dan went home` returned a 2 code: \nstdout: home is missing\nstderr: failed to go home")
}

func TestEscapeArgs(t *testing.T) {
	original := []string{"dan", "wears", "big jorts"}
	expected := []string{"dan", "wears", "'big jorts'"}
	actual := escapeArgs(original)
	assert.Equal(t, expected, actual)
}

func TestExists(t *testing.T) {
	sh := New(context.New(config.Project{}))
	assert.True(t, sh.Exists("git"))
	assert.False(t, sh.Exists("nonexistent_program.exe"))
}
