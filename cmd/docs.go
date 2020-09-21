package cmd

import (
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
		newDocsManCmd(),
		newDocsMdCmd(),
	)

	root.cmd = cmd
	return root
}
