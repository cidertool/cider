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

// Package cmd declares the command line interface for Cider
package cmd

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
	newRootCmd(version, exit).Execute(args)
}

type rootCmd struct {
	cmd   *cobra.Command
	debug bool
	exit  func(int)
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	var root = &rootCmd{
		exit: exit,
	}

	var cmd = &cobra.Command{
		Use:   "cider",
		Short: "Submit your builds to the Apple App Store in seconds",
		Long: `Cider  Copyright (C) 2020  Aaron Sky
This program comes with ABSOLUTELY NO WARRANTY; for details type ` + "`help'" + `.
This is free software, and you are welcome to redistribute it
under certain conditions; type ` + "`help'" + ` for details.`,
		Version:           version,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if root.debug {
				log.SetLevel(log.DebugLevel)
				log.Debug("debug logs enabled")
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&root.debug, "debug", false, "Enable debug mode")

	cmd.AddCommand(
		newInitCmd().cmd,
		newCheckCmd().cmd,
		newReleaseCmd().cmd,
		newCompletionsCmd().cmd,
	)

	if version == "dev" {
		cmd.AddCommand(newDocsCmd().cmd)
	}

	root.cmd = cmd

	return root
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
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
