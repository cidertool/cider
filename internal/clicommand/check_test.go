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

package clicommand

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCheckCmd(t *testing.T) {
	var cmd = newCheckCmd()

	var path = filepath.Join(t.TempDir(), "foo.yaml")

	var proj config.Project

	s, err := proj.String()
	assert.NoError(t, err)
	err = os.WriteFile(path, []byte(s), 0600)
	assert.NoError(t, err)

	err = cmd.cmd.Execute()
	assert.Error(t, err)

	cmd.config = path

	err = cmd.cmd.Execute()
	assert.NoError(t, err)
}
