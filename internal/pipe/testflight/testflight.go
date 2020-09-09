// Package testflight is a pipe that processes an app's release to Testflight
package testflight

import (
	"fmt"

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
	return "committing to testflight"
}

// Publish to Testflight.
func (p Pipe) Publish(ctx *context.Context) error {
	if ctx.PublishMode != context.PublishModeTestflight {
		return pipe.Skip("testflight")
	}
	client := client.New(ctx)
	for _, name := range ctx.AppsToRelease {
		app := ctx.Config.Apps[name]
		log.WithField("name", name).Info("preparing")
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

	buildVersionLog := fmt.Sprintf("%s (%s)", ctx.Version, *build.Attributes.Version)
	log.WithFields(log.Fields{
		"app":   *app.Attributes.BundleID,
		"build": buildVersionLog,
	}).Info("found resources")

	if ctx.SkipUpdateMetadata {
		log.Warn("skipping updating metdata")
	} else {
		log.Info("updating metadata")
		if err := updateBetaDetails(ctx, config, client, app, build); err != nil {
			return err
		}
	}

	if ctx.SkipSubmit {
		return pipe.ErrSkipSubmitEnabled
	}

	log.
		WithField("build", buildVersionLog).
		Info("submitting to testflight")

	return client.SubmitBetaApp(ctx, build)
}

func updateBetaDetails(ctx *context.Context, config config.App, client client.Client, app *asc.App, build *asc.Build) error {
	log.Debugf("updating %d beta app localizations", len(config.Testflight.Localizations))
	if err := client.UpdateBetaAppLocalizations(ctx, app, config.Testflight.Localizations); err != nil {
		return err
	}
	log.Debug("updating beta build details")
	if err := client.UpdateBetaBuildDetails(ctx, build, config.Testflight); err != nil {
		return err
	}
	log.Debugf("updating %d beta build localizations", len(config.Testflight.Localizations))
	if err := client.UpdateBetaBuildLocalizations(ctx, build, config.Testflight.Localizations); err != nil {
		return err
	}
	log.Debug("updating beta license agreement")
	if err := client.UpdateBetaLicenseAgreement(ctx, app, config.Testflight); err != nil {
		return err
	}
	log.Debug("updating build beta groups")
	if err := client.AssignBetaGroups(ctx, build, config.Testflight.BetaGroups); err != nil {
		return err
	}
	log.Debug("updating build beta testers")
	if err := client.AssignBetaTesters(ctx, build, config.Testflight.BetaTesters); err != nil {
		return err
	}
	if config.Testflight.ReviewDetails != nil {
		log.Debug("updating beta review details")
		if err := client.UpdateBetaReviewDetails(ctx, app, *config.Testflight.ReviewDetails); err != nil {
			return err
		}
	}
	return nil
}
