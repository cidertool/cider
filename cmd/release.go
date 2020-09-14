package cmd

import (
	"errors"
	"time"

	"github.com/apex/log"
	"github.com/cidertool/cider/internal/middleware"
	"github.com/cidertool/cider/internal/pipeline"
	"github.com/cidertool/cider/pkg/context"
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
	publishMode        context.PublishMode
	releaseAllApps     bool
	skipGit            bool
	skipUpdatePricing  bool
	skipUpdateMetadata bool
	skipSubmit         bool
	timeout            time.Duration
	versionOverride    string
	buildOverride      string
	currentDirectory   string
}

func newReleaseCmd() *releaseCmd {
	var root = &releaseCmd{}
	var cmd = &cobra.Command{
		Use:           "release [path]",
		Args:          cobra.MaximumNArgs(1),
		Short:         "Release the selected apps in the current project",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				root.opts.currentDirectory = args[0]
			}
			if root.opts.skipGit && root.opts.versionOverride == "" {
				return errors.New("if --skip-git is set, --set-version must also be set")
			}

			start := time.Now()

			log.Info(color.New(color.Bold).Sprint("releasing..."))

			_, err := releaseProject(root.opts)
			if err != nil {
				return wrapError(err, color.New(color.Bold).Sprintf("release failed after %0.2fs", time.Since(start).Seconds()))
			}

			log.Info(color.New(color.Bold).Sprintf("release succeeded after %0.2fs", time.Since(start).Seconds()))
			return nil
		},
	}

	cmd.Flags().StringVarP(&root.opts.config, "config", "f", "", "Load configuration from file")
	cmd.Flags().StringArrayVarP(&root.opts.appsToRelease, "app", "a", make([]string, 0), "App to release, using key name in configuration")
	cmd.Flags().Var(&root.opts.publishMode, "mode", `Publish mode (default: "testflight")`)
	cmd.Flags().BoolVarP(&root.opts.releaseAllApps, "all-apps", "A", false, "Release all apps")
	cmd.Flags().BoolVar(&root.opts.skipGit, "skip-git", false, "Skips deriving version information from Git. Must only be used in conjunction with --set-version")
	cmd.Flags().BoolVar(&root.opts.skipUpdatePricing, "skip-update-pricing", false, "Skips updating pricing")
	cmd.Flags().BoolVar(&root.opts.skipUpdateMetadata, "skip-update-metadata", false, "Skips updating metadata")
	cmd.Flags().BoolVar(&root.opts.skipSubmit, "skip-submit", false, "Skips submitting for review")
	cmd.Flags().StringVarP(&root.opts.versionOverride, "set-version", "V", "", "Version override to use instead of Git tags")
	cmd.Flags().StringVarP(&root.opts.buildOverride, "set-build", "B", "", `Build override to use instead of "latest".`)
	cmd.Flags().DurationVar(&root.opts.timeout, "timeout", 30*time.Minute, "Timeout to the entire release process")

	root.cmd = cmd
	return root
}

func releaseProject(options releaseOpts) (*context.Context, error) {
	var forceAllSkips bool
	cfg, err := loadConfig(options.config, options.currentDirectory)
	if err != nil {
		if err == ErrConfigNotFound {
			log.Warn(err.Error())
			log.Warn("using defaults and enabling all skips to avoid dangerous consequences...")
			forceAllSkips = true
		} else {
			return nil, err
		}
	}
	ctx, cancel := context.NewWithTimeout(cfg, options.timeout)
	defer cancel()
	setupReleaseContext(ctx, options, forceAllSkips)
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

func setupReleaseContext(ctx *context.Context, options releaseOpts, forceAllSkips bool) *context.Context {
	ctx.AppsToRelease = ctx.Config.AppsMatching(options.appsToRelease, options.releaseAllApps)
	if options.publishMode == "" {
		ctx.PublishMode = context.PublishModeTestflight
	} else {
		ctx.PublishMode = options.publishMode
	}
	ctx.SkipGit = options.skipGit || forceAllSkips
	ctx.SkipUpdatePricing = options.skipUpdatePricing || forceAllSkips
	ctx.SkipUpdateMetadata = options.skipUpdateMetadata || forceAllSkips
	ctx.SkipSubmit = options.skipSubmit || forceAllSkips
	ctx.Version = options.versionOverride
	ctx.Build = options.buildOverride
	ctx.CurrentDirectory = options.currentDirectory
	return ctx
}
