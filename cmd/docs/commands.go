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
	"github.com/spf13/cobra"
)

const defaultDocsPath = "docs"

// CmdConfig returns the cobra.Command for the man subcommand.
func CmdConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Generate configuration file documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsConfigCmd,
	}
}

// CmdMan returns the cobra.Command for the man subcommand.
func CmdMan() *cobra.Command {
	return &cobra.Command{
		Use:   "man",
		Short: "Generate man documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsManCmd,
	}
}

// CmdMarkdown returns the cobra.Command for the man subcommand.
func CmdMarkdown() *cobra.Command {
	return &cobra.Command{
		Use:   "md",
		Short: "Generate Markdown documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsMdCmd,
	}
}
