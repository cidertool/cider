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

package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

const licenseHeader = `/**
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
*/`

func main() {
	var cmd = &cobra.Command{
		Use:               "licensing",
		Short:             "Ensure every source file in the repo contains a license header",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE:              runLicensing,
	}

	cmd.SetArgs(os.Args[1:])

	if err := cmd.Execute(); err != nil {
		var code = 1

		var msg = "command failed"

		log.WithError(err).Error(msg)

		os.Exit(code)
	}
}

func runLicensing(cmd *cobra.Command, args []string) error {
	return filepath.Walk(".", checkFile)
}

func checkFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if filepath.Ext(path) != ".go" {
		return nil
	}

	f, err := os.ReadFile(path) // #nosec
	if err != nil {
		return err
	}

	source := string(f)

	if !strings.HasPrefix(source, licenseHeader) {
		source = licenseHeader + "\n\n" + source

		err := os.WriteFile(path, []byte(source), info.Mode())
		if err != nil {
			return err
		}

		log.Info(path)
	}

	return nil
}
