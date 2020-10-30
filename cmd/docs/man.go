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

package docs

import (
	"path/filepath"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func runDocsManCmd(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) == 0 {
		path = defaultDocsPath
	} else {
		path = args[0]
	}

	path = filepath.Join(path, "man")

	log.WithField("path", path).Info("generating man documentation")

	err := doc.GenManTreeFromOpts(cmd.Root(), doc.GenManTreeOptions{
		Path: path,
	})
	if err != nil {
		log.Error("generation failed")
	} else {
		log.Info("generation completed successfully")
	}

	return err
}
