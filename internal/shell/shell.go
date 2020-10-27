// Package shell wraps shell execution
package shell

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/apex/log"
	"github.com/cidertool/cider/pkg/context"
)

// Shell is an abstraction for shell-program execution meant to make
// testing and client design easier.
type Shell interface {
	NewCommand(program string, args ...string) *exec.Cmd
	Exec(cmd *exec.Cmd) (*CompletedProcess, error)
	Exists(program string) bool
	CurrentDirectory() string
}

// New returns a new shell bound to the provided context.
func New(ctx *context.Context) Shell {
	return &loginShell{ctx}
}

// CompletedProcess represents a subshell execution that finished, and includes its arguments, return code, and
// standard buffers as strings.
type CompletedProcess struct {
	Name       string
	Args       []string
	ReturnCode int
	Stdout     string
	Stderr     string
}

func newCompletedProcess(cmd *exec.Cmd) CompletedProcess {
	stdout := cmd.Stdout.(*bytes.Buffer)
	stderr := cmd.Stderr.(*bytes.Buffer)

	var stdoutString, stderrString string

	if stdout != nil {
		stdoutString = strings.TrimSpace(stdout.String())
	}

	if stderr != nil {
		stderrString = strings.TrimSpace(stderr.String())
	}

	return CompletedProcess{
		Name:       cmd.Path,
		Args:       cmd.Args,
		ReturnCode: cmd.ProcessState.ExitCode(),
		Stdout:     stdoutString,
		Stderr:     stderrString,
	}
}

type shellError struct {
	process CompletedProcess
}

func (e *shellError) Error() string {
	return fmt.Sprintf(
		"`%s %s` returned a %d code: \nstdout: %s\nstderr: %s",
		e.process.Name,
		strings.Join(e.process.Args, " "),
		e.process.ReturnCode,
		e.process.Stdout,
		e.process.Stderr,
	)
}

// loginShell is an empty struct that implements shell.Shell with default
// os.Exec subshell execution logic.
type loginShell struct {
	*context.Context
}

// NewCommand takes a program name and series of arguments and constructs an
// exec.Cmd object that can be manipulated and fed to Exec().
func (sh *loginShell) NewCommand(program string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(sh.Context, program, escapeArgs(arg)...) // #nosec

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = sh.Context.CurrentDirectory

	return cmd
}

func escapeArgs(args []string) []string {
	copy := make([]string, len(args))
	for i, arg := range args {
		copy[i] = shellescape.Quote(arg)
	}

	return copy
}

// Exec executes a command.
func (sh *loginShell) Exec(cmd *exec.Cmd) (*CompletedProcess, error) {
	log.WithField("args", cmd.Args).Debug(cmd.Path)

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	process := newCompletedProcess(cmd)

	log.WithFields(log.Fields{
		"code":   process.ReturnCode,
		"stdout": strings.TrimSpace(process.Stdout),
		"stderr": strings.TrimSpace(process.Stderr),
	}).Debugf("%s result", process.Name)

	var shellErr error
	if process.ReturnCode != 0 {
		shellErr = &shellError{process: process}
	}

	return &process, shellErr
}

// Exists returns whether or not a given program is installed.
func (sh *loginShell) Exists(program string) bool {
	path, err := exec.LookPath(program)
	if err != nil {
		return false
	}

	return path != ""
}

func (sh *loginShell) CurrentDirectory() string {
	return sh.Context.CurrentDirectory
}
