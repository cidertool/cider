package shell

import (
	"testing"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	sh := New(context.New(config.Project{}))
	cmd := sh.NewCommand("echo", "dogs")
	ps, err := sh.Exec(cmd)
	if err != nil {
		t.Error(err)
	}
	if ps.Stdout != "dogs" {
		t.Error("expected output \"dogs\" does not equal", ps.Stdout)
	}
}

func TestExec_Error(t *testing.T) {
	sh := New(context.New(config.Project{}))
	cmd := sh.NewCommand("exit", "1")
	ps, err := sh.Exec(cmd)
	if err == nil {
		t.Error("expected err return to not be nil")
	}
	if ps != nil {
		t.Error("expected process return to be nil")
	}
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
	expected := "`dan went home` returned a 2 code: \nstdout: home is missing\nstderr: failed to go home"
	if err.Error() != expected {
		t.Error("error message didn't match expected:", err)
	}
}

func TestEscapeArgs(t *testing.T) {
	original := []string{"dan", "wears", "big jorts"}
	expected := []string{"dan", "wears", "'big jorts'"}
	actual := escapeArgs(original)
	if len(expected) != len(actual) {
		t.Error("expected length", len(expected), "does not equal actual length", len(actual))
	}
	for i, exp := range expected {
		if actual[i] != exp {
			t.Error(actual[i], "does not equal", exp)
		}
	}
}

func TestExists(t *testing.T) {
	sh := New(context.New(config.Project{}))
	assert.True(t, sh.Exists("git"))
	assert.False(t, sh.Exists("nonexistent_program.exe"))
}
