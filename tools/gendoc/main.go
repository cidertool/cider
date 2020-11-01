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

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

const defaultDocsPath = "docs"

type rootCmd struct {
	cmd  *cobra.Command
	exit func(int)
}

func main() {
	var root = &rootCmd{
		exit: os.Exit,
	}

	var cmd = &cobra.Command{
		Use:               "gendoc",
		Short:             "Generate documentation for Cider",
		Args:              cobra.MaximumNArgs(1),
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, sub := range cmd.Commands() {
				if sub.Name() == "help" {
					continue
				}
				if err := sub.RunE(sub, args); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.AddCommand(
		cmdConfig(),
		cmdMan(),
		cmdMarkdown(),
	)

	root.cmd = cmd

	root.Execute(os.Args[1:])
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		var code = 1

		var msg = "command failed"

		log.WithError(err).Error(msg)

		cmd.exit(code)
	}
}

// CmdConfig returns the cobra.Command for the man subcommand.
func cmdConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Generate configuration file documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsConfigCmd,
	}
}

// CmdMan returns the cobra.Command for the man subcommand.
func cmdMan() *cobra.Command {
	return &cobra.Command{
		Use:   "man",
		Short: "Generate man documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsManCmd,
	}
}

// CmdMarkdown returns the cobra.Command for the man subcommand.
func cmdMarkdown() *cobra.Command {
	return &cobra.Command{
		Use:   "md",
		Short: "Generate Markdown documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsMdCmd,
	}
}
