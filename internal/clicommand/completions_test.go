package clicommand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletionsCmd(t *testing.T) {
	cmd := newCompletionsCmd()

	cmd.cmd.SetArgs([]string{})
	err := cmd.cmd.Execute()
	assert.Error(t, err)

	cmd.cmd.SetArgs([]string{"bash"})
	err = cmd.cmd.Execute()
	assert.NoError(t, err)

	cmd.cmd.SetArgs([]string{"zsh"})
	err = cmd.cmd.Execute()
	assert.NoError(t, err)

	cmd.cmd.SetArgs([]string{"fish"})
	err = cmd.cmd.Execute()
	assert.NoError(t, err)

	cmd.cmd.SetArgs([]string{"powershell"})
	err = cmd.cmd.Execute()
	assert.NoError(t, err)

	cmd.cmd.SetArgs([]string{"oil"})
	err = cmd.cmd.Execute()
	assert.Error(t, err)
}
