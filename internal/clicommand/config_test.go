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
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Happy_CustomPath(t *testing.T) {
	var path = filepath.Join(t.TempDir(), "foo.yaml")

	var proj config.Project

	s, err := proj.String()
	assert.NoError(t, err)
	err = ioutil.WriteFile(path, []byte(s), 0600)
	assert.NoError(t, err)
	cfg, err := loadConfig(path, "")
	assert.NoError(t, err)
	assert.Empty(t, cfg)
}

func TestConfig_Happy_DefaultPath(t *testing.T) {
	var folder = t.TempDir()

	var path = filepath.Join(folder, "cider.yaml")

	var proj config.Project

	s, err := proj.String()
	assert.NoError(t, err)
	err = ioutil.WriteFile(path, []byte(s), 0600)
	assert.NoError(t, err)
	cfg, err := loadConfig("", folder)
	assert.NoError(t, err)
	assert.Empty(t, cfg)
}

func TestConfig_Err_DoesntExist(t *testing.T) {
	cfg, err := loadConfig("", t.TempDir())
	assert.Error(t, err)
	assert.Empty(t, cfg)
}
