package store

import (
	"github.com/aaronsky/applereleaser/internal/client"
	"github.com/aaronsky/applereleaser/internal/pipe"
	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
	"github.com/apex/log"
)

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "committing to app store"
}

// Publish to App Store Review.
func (p Pipe) Publish(ctx *context.Context) error {
	if ctx.PublishMode != context.PublishModeAppStore {
		return pipe.Skip("testflight")
	}
	client := client.New(ctx)
	for _, name := range ctx.AppsToRelease {
		app := ctx.Config.Apps[name]
		log.WithField("testflight", name).Info("updating metadata")
		err := doRelease(ctx, app, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func doRelease(ctx *context.Context, config config.App, client client.Client) error {
	app, err := client.GetAppForBundleID(ctx, config.BundleID)
	if err != nil {
		return err
	}
	build, err := client.GetRelevantBuild(ctx, app)
	if err != nil {
		return err
	}
	version, err := client.CreateVersionIfNeeded(ctx, app, build, config.Versions)
	if err != nil {
		return err
	}
	if !ctx.SkipUpdateMetadata {
		if err := updateVersionDetails(ctx, config, client, app, version); err != nil {
			return err
		}
	}
	if ctx.SkipSubmit {
		return pipe.ErrSkipSubmitEnabled
	}
	return client.SubmitApp(ctx, version)
}

func updateVersionDetails(ctx *context.Context, config config.App, client client.Client, app *asc.App, version *asc.AppStoreVersion) error {
	if err := client.UpdateAppLocalizations(ctx, app, config.Localizations); err != nil {
		return err
	}
	if err := client.UpdateVersionLocalizations(ctx, version, config.Versions.Localizations); err != nil {
		return err
	}
	if config.Versions.IDFADeclaration != nil {
		if err := client.UpdateIDFADeclaration(ctx, version, *config.Versions.IDFADeclaration); err != nil {
			return err
		}
	}
	if config.Versions.RoutingCoverage != nil {
		if err := client.UploadRoutingCoverage(ctx, version, *config.Versions.RoutingCoverage); err != nil {
			return err
		}
	}
	if config.Versions.ReviewDetails != nil {
		if err := client.UpdateReviewDetails(ctx, version, *config.Versions.ReviewDetails); err != nil {
			return err
		}
	}
	return nil
}
