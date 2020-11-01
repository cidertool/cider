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
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// ErrUnsupportedShell happens when a shell that Cobra does not support is
// passed for the completions command.
var ErrUnsupportedShell = errors.New("shell for completions is unsupported")

type completionsCmd struct {
	cmd *cobra.Command
}

func newCompletionsCmd() *completionsCmd {
	var root = &completionsCmd{}

	var cmd = &cobra.Command{
		Use:       "completions [bash|zsh|fish|powershell]",
		Short:     "Generate shell completions",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
			return ErrUnsupportedShell
		},
	}

	root.cmd = cmd

	return root
}
