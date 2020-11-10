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

// Package clicommand declares the command line interface for Cider.
package clicommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Execute is the primary function to initiate the command line interface for Cider.
func Execute(version string, exit func(int), args []string) {
	if os.Getenv("CI") != "" {
		color.NoColor = false
	}

	log.SetHandler(cli.Default)

	fmt.Println()
	defer fmt.Println()
	NewRoot(version, exit).Execute(args)
}

// Root defines a rough structure for a root command type.
type Root struct {
	Cmd   *cobra.Command
	debug bool
	exit  func(int)
}

// NewRoot creates a new instance of the root command for the cider executable.
func NewRoot(version string, exit func(int)) *Root {
	var root = &Root{
		exit: exit,
	}

	var cmd = &cobra.Command{
		Use:               "cider",
		Short:             "Submit your builds to the Apple App Store in seconds",
		Version:           version,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		PersistentPreRun:  root.setDebug,
	}

	cmd.PersistentFlags().BoolVar(&root.debug, "debug", false, "Enable debug mode")

	cmd.AddCommand(
		newInitCmd().cmd,
		newCheckCmd().cmd,
		newReleaseCmd().cmd,
		newCompletionsCmd().cmd,
	)

	root.Cmd = cmd

	return root
}

// Execute executes the root command.
func (cmd *Root) Execute(args []string) {
	cmd.Cmd.SetArgs(args)

	if err := cmd.Cmd.Execute(); err != nil {
		var code = 1

		var msg = "command failed"

		var eerr *exitError

		if ok := errors.As(err, &eerr); ok {
			code = eerr.code

			if eerr.details != "" {
				msg = eerr.details
			}
		}

		log.WithError(err).Error(msg)

		cmd.exit(code)
	}
}

func (cmd *Root) setDebug(c *cobra.Command, args []string) {
	if cmd.debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logs enabled")
	}
}
