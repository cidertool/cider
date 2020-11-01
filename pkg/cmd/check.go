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
	"fmt"

	"github.com/apex/log"
	"github.com/cidertool/cider/internal/pipe/defaults"
	"github.com/cidertool/cider/pkg/context"
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
		Long:          `Use to validate your configuration file.`,
		Example:       "cider check",
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
