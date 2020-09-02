package cmd

import (
	"time"

	"github.com/aaronsky/applereleaser/internal/middleware"
	"github.com/aaronsky/applereleaser/internal/pipeline"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type releaseCmd struct {
	cmd  *cobra.Command
	opts releaseOpts
}

type releaseOpts struct {
	config             string
	appsToRelease      []string
	releaseAllApps     bool
	skipUpdateMetadata bool
	skipSubmit         bool
	timeout            time.Duration
	currentDirectory   string
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
			start := time.Now()

			log.Infof(color.New(color.Bold).Sprint("releasing..."))

			_, err := releaseProject(root.opts)
			if err != nil {
				return wrapError(err, color.New(color.Bold).Sprintf("release failed after %0.2fs", time.Since(start).Seconds()))
			}

			log.Infof(color.New(color.Bold).Sprintf("release succeeded after %0.2fs", time.Since(start).Seconds()))
			return nil
		},
	}

	cmd.Flags().StringVarP(&root.opts.config, "config", "f", "", "Load configuration from file")
	cmd.Flags().StringArrayVarP(&root.opts.appsToRelease, "app", "a", make([]string, 0), "App to release, using key name in configuration")
	cmd.Flags().BoolVarP(&root.opts.releaseAllApps, "all-apps", "A", false, "Release all apps")
	cmd.Flags().BoolVar(&root.opts.skipUpdateMetadata, "skip-update-metadata", false, "Skips updating metadata")
	cmd.Flags().BoolVar(&root.opts.skipSubmit, "skip-submit", false, "Skips submitting for review")
	cmd.Flags().DurationVar(&root.opts.timeout, "timeout", 30*time.Minute, "Timeout to the entire release process")

	root.cmd = cmd
	return root
}

func releaseProject(options releaseOpts) (*context.Context, error) {
	cfg, err := loadConfig(options.config, options.currentDirectory)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.NewWithTimeout(cfg, options.timeout)
	defer cancel()
	setupReleaseContext(ctx, options)
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

func setupReleaseContext(ctx *context.Context, options releaseOpts) *context.Context {
	ctx.AppsToRelease = ctx.Config.AppsMatching(options.appsToRelease, options.releaseAllApps)
	ctx.SkipUpdateMetadata = options.skipUpdateMetadata
	ctx.SkipSubmit = options.skipSubmit
	ctx.CurrentDirectory = options.currentDirectory
	return ctx
}
