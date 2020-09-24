package cmd

import (
	"errors"
	"time"

	"github.com/apex/log"
	"github.com/cidertool/cider/internal/middleware"
	"github.com/cidertool/cider/internal/pipeline"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type releaseCmd struct {
	cmd  *cobra.Command
	opts releaseOpts
}

type releaseOpts struct {
	config              string
	appsToRelease       []string
	publishMode         context.PublishMode
	maxProcesses        int
	releaseAllApps      bool
	skipGit             bool
	skipUpdatePricing   bool
	skipUpdateMetadata  bool
	skipSubmit          bool
	timeout             time.Duration
	versionOverride     string
	buildOverride       string
	betaGroupsOverride  []string
	betaTestersOverride []string
	currentDirectory    string
}

func newReleaseCmd() *releaseCmd {
	var root = &releaseCmd{}
	var cmd = &cobra.Command{
		Use:   "release [path]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Release the selected apps in the current project",
		Long: `Release the selected apps in the current project.
		
You can provide a path to a project directory as an argument to be the root directory
of all relative path expansions in the program, such as the Git repository, preview sets,
and screenshot resources. An exception to this is if you set a custom configuration file
path with the ` + "`--config`" + ` flag.

Additionally, Cider requires a few environment variables to be set in order to operate.
They each correspond to an element of authorization described by the Apple Developer Documentation.

- ` + "`ASC_KEY_ID`" + `: The key's ID.
- ` + "`ASC_ISSUER_ID`" + `: Your team's issuer ID.
- ` + "`ASC_PRIVATE_KEY`" + ` or ` + "`ASC_PRIVATE_KEY_PATH`" + `: The .p8 private key issued by Apple.

These three values each have varying degrees of sensetivity and should be treated as secrets. Store
them securely in your environment so Cider can leverage them safely.

More info: https://developer.apple.com/documentation/appstoreconnectapi/creating_api_keys_for_app_store_connect_api`,

		Example: `cider release --mode=appstore --set-version="1.0"`,

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

	// Operative options

	cmd.Flags().StringVarP(
		&root.opts.config,
		"config",
		"f",
		"",
		"Load configuration from file",
	)
	cmd.Flags().StringArrayVarP(
		&root.opts.appsToRelease,
		"app",
		"a",
		[]string{},
		`Process the given app, providing the app key name used in your configuration file.

This flag can be provided repeatedly for each app you want to process. You can omit
this flag if your configuration file has only one app defined.`,
	)
	cmd.Flags().BoolVarP(
		&root.opts.releaseAllApps,
		"all-apps",
		"A",
		false,
		`Process all apps in the configuration file. Supercedes any usage of the `+"`--app`"+` flag.`,
	)
	cmd.Flags().Var(
		&root.opts.publishMode,
		"mode",
		`Mode used to declare the publishing target for submission.
		
The default is "testflight" for submitting to Testflight, and the other alternative
option is "appstore" for submitting to the App Store.`,
	)
	cmd.Flags().IntVarP(
		&root.opts.maxProcesses,
		"max-processes",
		"p",
		1,
		`Run certain metadata syncing and asset uploading logic in parallel with
the maximum allowable concurrency.`,
	)
	cmd.Flags().DurationVar(
		&root.opts.timeout,
		"timeout",
		30*time.Minute,
		`Timeout for the entire release process.
		
If the command takes longer than this amount of time to run, Cider will abort.`,
	)

	// Skip options

	cmd.Flags().BoolVar(
		&root.opts.skipGit,
		"skip-git",
		false,
		`Skips deriving version information from Git. Must only be used in conjunction with the `+"`--set-version`"+` flag.`,
	)
	cmd.Flags().BoolVar(
		&root.opts.skipUpdatePricing,
		"skip-update-pricing",
		false,
		"Skips updating app pricing",
	)
	cmd.Flags().BoolVar(
		&root.opts.skipUpdateMetadata,
		"skip-update-metadata",
		false,
		"Skips updating metadata (app info, localizations, assets, review details, etc.)",
	)
	cmd.Flags().BoolVar(
		&root.opts.skipSubmit,
		"skip-submit",
		false,
		"Skips submitting for review",
	)

	// Setting options

	cmd.Flags().StringVarP(
		&root.opts.versionOverride,
		"set-version",
		"V",
		"",
		`Version string override to use instead of parsing Git tags. Corresponds to the
CFBundleShortVersionString of your build.

Cider expects this string to follow the Major.Minor.Patch semantics outlined in Apple documentation
and Semantic Versioning (semver). If this flag is omitted, Git will be leveraged to determine the
latest tag. The tag will be used to calculate the version string under the same constraints.`,
	)
	cmd.Flags().StringVarP(
		&root.opts.buildOverride,
		"set-build",
		"B",
		"",
		`Build override to use instead of "latest". Corresponds to the CFBundleVersion
of your build.
		
The default behavior without this flag is to select the latest build. In both cases,
if the selected build has an invalid processing state, Cider will abort with an error
to ensure your release is handled safely.`,
	)
	cmd.Flags().StringArrayVar(
		&root.opts.betaGroupsOverride,
		"set-beta-group",
		[]string{},
		`Provide names of beta groups to release to instead of using
the configuration file.`,
	)
	cmd.Flags().StringArrayVar(
		&root.opts.betaTestersOverride,
		"set-beta-tester",
		[]string{},
		`Provide email addresses of beta testers to release to instead of
using the configuration file.`,
	)

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
	ctx.MaxProcesses = options.maxProcesses
	ctx.SkipGit = options.skipGit || forceAllSkips
	ctx.SkipUpdatePricing = options.skipUpdatePricing || forceAllSkips
	ctx.SkipUpdateMetadata = options.skipUpdateMetadata || forceAllSkips
	ctx.SkipSubmit = options.skipSubmit || forceAllSkips
	ctx.Version = options.versionOverride
	ctx.Build = options.buildOverride

	if len(options.betaGroupsOverride) > 0 || len(options.betaTestersOverride) > 0 {
		var betaGroups = make([]config.BetaGroup, len(options.betaGroupsOverride))
		var betaTesters = make([]config.BetaTester, len(options.betaTestersOverride))

		for i, groupName := range options.betaGroupsOverride {
			betaGroups[i] = config.BetaGroup{Name: groupName}
		}

		for i, email := range options.betaTestersOverride {
			betaTesters[i] = config.BetaTester{Email: email}
		}

		for appName, app := range ctx.Config {
			if len(options.betaGroupsOverride) > 0 {
				app.Testflight.BetaGroups = betaGroups
			}
			if len(betaTesters) > 0 {
				app.Testflight.BetaTesters = betaTesters
			}
			ctx.Config[appName] = app
		}
	}

	ctx.CurrentDirectory = options.currentDirectory

	return ctx
}
