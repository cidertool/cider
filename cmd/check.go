package cmd

import (
	"fmt"

	"github.com/aaronsky/applereleaser/internal/pipe/defaults"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/apex/log"
	"github.com/fatih/color"
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
			var ctx = context.New(cfg)

			if err := context.NewInterrupt().Run(ctx, func() error {
				log.Info(color.New(color.Bold).Sprint("checking config:"))
				return defaults.Pipe{}.Run(ctx)
			}); err != nil {
				log.WithError(err).Error(color.New(color.Bold).Sprintf("config is invalid"))
				return fmt.Errorf("invalid config: %w", err)
			}

			log.Info(color.New(color.Bold).Sprintf("config is valid"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&root.config, "config", "f", "", "Configuration file to check")

	root.cmd = cmd
	return root
}
