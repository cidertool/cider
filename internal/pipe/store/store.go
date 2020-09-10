// Package store is a pipe that processes an app's release to the App Store
package store

import (
	"github.com/aaronsky/asc-go/asc"
	"github.com/apex/log"
	"github.com/cidertool/cider/internal/client"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
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
		return pipe.Skip("app store")
	}
	client := client.New(ctx)
	for _, name := range ctx.AppsToRelease {
		app := ctx.Config.Apps[name]
		log.WithField("app", name).Info("updating metadata")
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
	isInitial, err := client.ReleaseForAppIsInitial(ctx, app)
	if err != nil {
		return err
	}
	ctx.VersionIsInitialRelease = isInitial
	build, err := client.GetRelevantBuild(ctx, app)
	if err != nil {
		return err
	}
	version, err := client.CreateVersionIfNeeded(ctx, app, build, config.Versions)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"app":     *app.Attributes.BundleID,
		"build":   *build.Attributes.Version,
		"version": *version.Attributes.VersionString,
	}).Info("found resources")

	if ctx.SkipUpdateMetadata {
		log.Warn("skipping updating metdata")
	} else {
		log.Info("updating metadata")
		if err := updateVersionDetails(ctx, config, client, app, version); err != nil {
			return err
		}
	}

	if ctx.SkipSubmit {
		return pipe.ErrSkipSubmitEnabled
	}

	log.
		WithField("version", *version.Attributes.VersionString).
		Info("submitting to app store")
	return client.SubmitApp(ctx, version)
}

func updateVersionDetails(ctx *context.Context, config config.App, client client.Client, app *asc.App, version *asc.AppStoreVersion) error {
	appInfo, err := client.GetAppInfo(ctx, app)
	if err != nil {
		return err
	}
	log.Debug("updating app details")
	if err := client.UpdateApp(ctx, app, appInfo, config); err != nil {
		return err
	}
	log.Debugf("updating %d app localizations", len(config.Localizations))
	if err := client.UpdateAppLocalizations(ctx, app, appInfo, config.Localizations); err != nil {
		return err
	}
	log.Debugf("updating %d app store version localizations", len(config.Versions.Localizations))
	if err := client.UpdateVersionLocalizations(ctx, version, config.Versions.Localizations); err != nil {
		return err
	}
	if config.Versions.IDFADeclaration != nil {
		log.Debug("updating IDFA declaration")
		if err := client.UpdateIDFADeclaration(ctx, version, *config.Versions.IDFADeclaration); err != nil {
			return err
		}
	}
	if config.Versions.RoutingCoverage != nil {
		log.Debug("uploading routing coverage asset")
		if err := client.UploadRoutingCoverage(ctx, version, *config.Versions.RoutingCoverage); err != nil {
			return err
		}
	}
	if config.Versions.ReviewDetails != nil {
		log.Debug("updating review details")
		if err := client.UpdateReviewDetails(ctx, version, *config.Versions.ReviewDetails); err != nil {
			return err
		}
	}
	return nil
}
