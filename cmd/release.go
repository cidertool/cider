package cmd

import (
	"time"

	"github.com/aaronsky/applereleaser/internal/middleware"
	"github.com/aaronsky/applereleaser/internal/pipeline"
	"github.com/aaronsky/applereleaser/pkg/context"

	"github.com/spf13/cobra"
)

type releaseCmd struct {
	cmd  *cobra.Command
	opts releaseOpts
}

type releaseOpts struct {
	config           string
	skipPublish      bool
	timeout          time.Duration
	currentDirectory string
}

func newReleaseCmd() *releaseCmd {
	var root = &releaseCmd{}
	var cmd = &cobra.Command{
		Use:           "release [path]",
		Args:          cobra.MaximumNArgs(1),
		Short:         "Releases all the apps in the current project",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				root.opts.currentDirectory = args[0]
			}
			_, err := releaseProject(root.opts)
			return err
		},
	}

	cmd.Flags().StringVarP(&root.opts.config, "config", "f", "", "Load configuration from file")
	cmd.Flags().BoolVar(&root.opts.skipPublish, "skip-publish", false, "Skips publishing artifacts")
	cmd.Flags().DurationVar(&root.opts.timeout, "timeout", 30*time.Minute, "Timeout to the entire release process")

	root.cmd = cmd
	return root
}

func releaseProject(options releaseOpts) (*context.Context, error) {
	cfg, err := loadConfig(options.config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.NewWithTimeout(cfg, options.timeout)
	defer cancel()
	ctx.SkipPublish = options.skipPublish
	ctx.CurrentDirectory = options.currentDirectory
	return ctx, context.NewInterrupt().Run(ctx, func() error {
		for _, pipe := range pipeline.Pipeline {
			if err := middleware.Logging(
				pipe.String(),
				middleware.ErrHandler(pipe.Run),
				middleware.DefaultInitialPadding,
			)(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}
