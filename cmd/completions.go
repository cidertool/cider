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
