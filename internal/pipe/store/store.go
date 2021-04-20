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

// Package store is a pipe that processes an app's release to the App Store
package store

import (
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/client"
	"github.com/cidertool/cider/internal/log"
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
			return pipe.ErrMissingApp{Name: name}
		}

		ctx.Log.WithField("app", name).Info("updating metadata")

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

	ctx.Log.WithFields(log.Fields{
		"app":     *app.Attributes.BundleID,
		"build":   *build.Attributes.Version,
		"version": *version.Attributes.VersionString,
	}).Info("found resources")

	if ctx.SkipUpdateMetadata {
		ctx.Log.Warn("skipping updating metdata")
	} else {
		ctx.Log.Info("updating metadata")
		if err := p.updateVersionDetails(ctx, config, app, version); err != nil {
			return err
		}
	}

	if ctx.SkipSubmit {
		return pipe.ErrSkipSubmitEnabled
	}

	if config.Versions.PhasedReleaseEnabled && !ctx.VersionIsInitialRelease {
		ctx.Log.Info("preparing phased release details")

		if err := p.Client.EnablePhasedRelease(ctx, version.ID); err != nil {
			return err
		}
	}

	ctx.Log.
		WithField("version", *version.Attributes.VersionString).
		Info("submitting to app store")

	return p.Client.SubmitApp(ctx, version.ID)
}

func (p *Pipe) updateVersionDetails(ctx *context.Context, config config.App, app *asc.App, version *asc.AppStoreVersion) error {
	appInfo, err := p.Client.GetAppInfo(ctx, app.ID)
	if err != nil {
		return err
	}

	ctx.Log.Info("updating app details")

	if err := p.Client.UpdateApp(ctx, app.ID, appInfo.ID, version.ID, config); err != nil {
		return err
	}

	ctx.Log.Infof("updating %d app localizations", len(config.Localizations))

	if err := p.Client.UpdateAppLocalizations(ctx, app.ID, config.Localizations); err != nil {
		return err
	}

	ctx.Log.Infof("updating %d app store version localizations", len(config.Versions.Localizations))

	if err := p.Client.UpdateVersionLocalizations(ctx, version.ID, config.Versions.Localizations); err != nil {
		return err
	}

	if config.Versions.IDFADeclaration != nil {
		ctx.Log.Info("updating IDFA declaration")

		if err := p.Client.UpdateIDFADeclaration(ctx, version.ID, *config.Versions.IDFADeclaration); err != nil {
			return err
		}
	}

	if config.Versions.RoutingCoverage != nil {
		ctx.Log.Info("uploading routing coverage asset")

		if err := p.Client.UploadRoutingCoverage(ctx, version.ID, *config.Versions.RoutingCoverage); err != nil {
			return err
		}
	}

	if config.Versions.ReviewDetails != nil {
		ctx.Log.Info("updating review details")

		if err := p.Client.UpdateReviewDetails(ctx, version.ID, *config.Versions.ReviewDetails); err != nil {
			return err
		}
	}

	return nil
}
