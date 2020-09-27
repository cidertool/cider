// Package shelltest provides utilities for mocking the login shell.
package shelltest

import (
	"os/exec"
	"strconv"

	"github.com/apex/log"
	"github.com/cidertool/cider/internal/shell"
	"github.com/cidertool/cider/pkg/context"
)

// Shell is a type that conforms to shell.Shell.
type Shell struct {
	Context           *context.Context
	SupportedPrograms map[string]bool
	Commands          []Command
	index             int
}

// Command represents the result of some executed command.
type Command struct {
	ReturnCode int
	Stdout     string
	Stderr     string
}

// NewCommand takes a program and arguments and constructs a new exec.Cmd instance.
func (sh *Shell) NewCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...) // #nosec
}

// Exec executes the command.
func (sh *Shell) Exec(cmd *exec.Cmd) (*shell.CompletedProcess, error) {
	if sh.index >= len(sh.Commands) {
		log.WithFields(log.Fields{
			"index":   sh.index,
			"count":   len(sh.Commands),
			"command": cmd.String(),
		}).Fatal("index out of bounds")
	}
	currentCommand := sh.Commands[sh.index]
	ps := shell.CompletedProcess{
		Name:       cmd.Path,
		Args:       cmd.Args,
		ReturnCode: currentCommand.ReturnCode,
		Stdout:     currentCommand.Stdout,
		Stderr:     currentCommand.Stderr,
	}
	sh.index++
	var err error
	if currentCommand.ReturnCode != 0 {
		err = &shellError{
			Process: ps,
		}
	}
	return &ps, err
}

// Exists returns whether the given program exists.
//
// This implementation of the method returns true if the SupportedPrograms
// field is empty, otherwise it checks the value of the key in that map.
// Use SupportedPrograms to mock the existence of a program in the PATH.
func (sh *Shell) Exists(program string) bool {
	if len(sh.SupportedPrograms) == 0 {
		return true
	}
	return sh.SupportedPrograms[program]
}

// CurrentDirectory returns the current directory.
func (sh *Shell) CurrentDirectory() string {
	return sh.Context.CurrentDirectory
}

type shellError struct {
	Process shell.CompletedProcess
}

func (err *shellError) Error() string {
	return strconv.Itoa(err.Process.ReturnCode)
}
