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
	"errors"
	"os"
	"path/filepath"

	"github.com/cidertool/cider/pkg/config"
)

// ErrConfigNotFound happens if a config file could not be found at any of the default locations.
var ErrConfigNotFound = errors.New("config file not found at any default path")

func loadConfig(path string, wd string) (config.Project, error) {
	if path != "" {
		return config.Load(path)
	}

	for _, f := range [4]string{
		".cider.yml",
		".cider.yaml",
		"cider.yml",
		"cider.yaml",
	} {
		proj, err := config.Load(filepath.Join(wd, f))
		if err != nil && os.IsNotExist(err) {
			continue
		}

		return proj, err
	}

	return config.Project{}, ErrConfigNotFound
}
