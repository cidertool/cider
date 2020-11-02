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

// Package shelltest provides utilities for mocking the login shell.
package shelltest

import (
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/cidertool/cider/internal/shell"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

// ErrCommandOverflow happens when the command that is trying to be run would result in a
// buffer overflow in the Shell mock.
type ErrCommandOverflow struct {
	Index   int
	Len     int
	Command string
}

func (e ErrCommandOverflow) Error() string {
	return fmt.Sprintf("command out of bounds: i=%d,len=%d (command: `%s`)", e.Index, e.Len, e.Command)
}

// Shell is a type that conforms to shell.Shell.
type Shell struct {
	T                   *testing.T
	Context             *context.Context
	SupportedPrograms   map[string]bool
	Commands            []Command
	index               int
	expectOverflowError bool
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
		err := ErrCommandOverflow{
			Index:   sh.index,
			Len:     len(sh.Commands),
			Command: cmd.String(),
		}

		if !sh.expectOverflowError {
			assert.FailNow(sh.T, err.Error())
		}

		return nil, err
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
