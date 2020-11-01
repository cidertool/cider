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

package cmd

import (
	"github.com/cidertool/cider/pkg/cmd/docs"
	"github.com/spf13/cobra"
)

type docsCmd struct {
	cmd *cobra.Command
}

func newDocsCmd() *docsCmd {
	var root = &docsCmd{}

	var cmd = &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation for Cider",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, sub := range cmd.Commands() {
				if err := sub.RunE(sub, args); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.AddCommand(
		docs.CmdConfig(),
		docs.CmdMan(),
		docs.CmdMarkdown(),
	)

	root.cmd = cmd

	return root
}
