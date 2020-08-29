package cmd

import (
	"os"

	"github.com/aaronsky/applereleaser/internal/static"
	"github.com/spf13/cobra"
)

type initCmd struct {
	cmd    *cobra.Command
	config string
}

func newInitCmd() *initCmd {
	var root = &initCmd{}
	var cmd = &cobra.Command{
		Use:           "init",
		Short:         "Generates an .applereleaser.yml file",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.OpenFile(root.config, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.WriteString(static.ExampleConfig); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&root.config, "config", "f", ".applereleaser.yml", "Configuration file to load or create")

	root.cmd = cmd
	return root
}
