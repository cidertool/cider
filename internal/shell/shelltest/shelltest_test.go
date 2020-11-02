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

package shelltest

import (
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	ctx := context.New(config.Project{})
	ctx.CurrentDirectory = "TEST"
	sh := Shell{
		Context: ctx,
		Commands: []Command{
			{
				ReturnCode: 0,
				Stdout:     "TEST",
				Stderr:     "TEST",
			},
			{
				ReturnCode: 128,
				Stdout:     "TEST",
				Stderr:     "TEST",
			},
		},
	}

	dir := sh.CurrentDirectory()
	assert.Equal(t, dir, ctx.CurrentDirectory)

	exists := sh.Exists("echo")
	assert.True(t, exists)

	sh.SupportedPrograms = map[string]bool{
		"echo": false,
	}
	exists = sh.Exists("echo")
	assert.False(t, exists)

	cmd := sh.NewCommand("echo", "true")
	assert.NotNil(t, cmd)

	proc, err := sh.Exec(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, proc)

	proc, err = sh.Exec(cmd)
	assert.EqualError(t, err, "128")
	assert.NotNil(t, proc)

	sh.expectOverflowError = true
	expectedErr := ErrCommandOverflow{
		Index:   2,
		Len:     2,
		Command: cmd.String(),
	}
	_, err = sh.Exec(cmd)
	assert.EqualError(t, err, expectedErr.Error())
}
