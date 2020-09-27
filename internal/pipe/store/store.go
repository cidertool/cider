// Package store is a pipe that processes an app's release to the App Store
package store

import (
	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/client"
	"github.com/cidertool/cider/internal/pipe"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

// Pipe is a global hook pipe.
type Pipe struct {
	Client client.Client
}

// String is the name of this pipe.
func (Pipe) String() string {
	return "committing to app store"
}

// Publish to App Store Review.
func (p *Pipe) Publish(ctx *context.Context) error {
	if p.Client == nil {
		p.Client = client.New(ctx)
	}
	for _, name := range ctx.AppsToRelease {
		app, ok := ctx.Config[name]
		if !ok {
			return pipe.ErrMissingApp(name)
		}
		log.WithField("app", name).Info("updating metadata")
		err := p.doRelease(ctx, app)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipe) doRelease(ctx *context.Context, config config.App) error {
	app, err := p.Client.GetAppForBundleID(ctx, config.BundleID)
	if err != nil {
		return err
	}
	isInitial, err := p.Client.ReleaseForAppIsInitial(ctx, app.ID)
	if err != nil {
		return err
	}
	ctx.VersionIsInitialRelease = isInitial
	build, err := p.Client.GetBuild(ctx, app)
	if err != nil {
		return err
	}
	version, err := p.Client.CreateVersionIfNeeded(ctx, app.ID, build.ID, config.Versions)
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
		if err := p.updateVersionDetails(ctx, config, app, version); err != nil {
			return err
		}
	}

	if ctx.SkipSubmit {
		return pipe.ErrSkipSubmitEnabled
	}

	if config.Versions.PhasedReleaseEnabled && !ctx.VersionIsInitialRelease {
		log.Info("preparing phased release details")
		if err := p.Client.EnablePhasedRelease(ctx, version.ID); err != nil {
			return err
		}
	}
	log.
		WithField("version", *version.Attributes.VersionString).
		Info("submitting to app store")
	return p.Client.SubmitApp(ctx, version.ID)
}

func (p *Pipe) updateVersionDetails(ctx *context.Context, config config.App, app *asc.App, version *asc.AppStoreVersion) error {
	appInfo, err := p.Client.GetAppInfo(ctx, app.ID)
	if err != nil {
		return err
	}
	log.Info("updating app details")
	if err := p.Client.UpdateApp(ctx, app.ID, appInfo.ID, version.ID, config); err != nil {
		return err
	}
	log.Infof("updating %d app localizations", len(config.Localizations))
	if err := p.Client.UpdateAppLocalizations(ctx, app.ID, config.Localizations); err != nil {
		return err
	}
	log.Infof("updating %d app store version localizations", len(config.Versions.Localizations))
	if err := p.Client.UpdateVersionLocalizations(ctx, version.ID, config.Versions.Localizations); err != nil {
		return err
	}
	if config.Versions.IDFADeclaration != nil {
		log.Info("updating IDFA declaration")
		if err := p.Client.UpdateIDFADeclaration(ctx, version.ID, *config.Versions.IDFADeclaration); err != nil {
			return err
		}
	}
	if config.Versions.RoutingCoverage != nil {
		log.Info("uploading routing coverage asset")
		if err := p.Client.UploadRoutingCoverage(ctx, version.ID, *config.Versions.RoutingCoverage); err != nil {
			return err
		}
	}
	if config.Versions.ReviewDetails != nil {
		log.Info("updating review details")
		if err := p.Client.UpdateReviewDetails(ctx, version.ID, *config.Versions.ReviewDetails); err != nil {
			return err
		}
	}
	return nil
}
