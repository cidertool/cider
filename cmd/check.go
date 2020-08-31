package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type checkCmd struct {
	cmd    *cobra.Command
	config string
}

func newCheckCmd() *checkCmd {
	var root = &checkCmd{}
	var cmd = &cobra.Command{
		Use:           "check",
		Short:         "Checks if the configuration is valid",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(root.config, "")
			if err != nil {
				return err
			}
			fmt.Println(cfg.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&root.config, "config", "f", "", "Configuration file to check")

	root.cmd = cmd
	return root
}
