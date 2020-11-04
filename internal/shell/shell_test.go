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
	assert.NotNil(t, ps)
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

func TestCurrentDirectory(t *testing.T) {
	ctx := context.New(config.Project{})
	ctx.CurrentDirectory = "TEST"
	sh := New(ctx)
	assert.Equal(t, ctx.CurrentDirectory, sh.CurrentDirectory())
}
