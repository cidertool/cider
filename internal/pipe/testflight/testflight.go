// Package testflight is a pipe that processes an app's release to Testflight
package testflight

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
	if !ctx.SkipUpdateMetadata {
		if err := updateBetaDetails(ctx, config, client, app, build); err != nil {
			return err
		}
	}
	if ctx.SkipSubmit {
		return pipe.ErrSkipSubmitEnabled
	}
	return client.SubmitBetaApp(ctx, build)
}

func updateBetaDetails(ctx *context.Context, config config.App, client client.Client, app *asc.App, build *asc.Build) error {
	if err := client.UpdateBetaAppLocalizations(ctx, app, config.Testflight.Localizations); err != nil {
		return err
	}
	if err := client.UpdateBetaBuildDetails(ctx, build, config.Testflight); err != nil {
		return err
	}
	if err := client.UpdateBetaBuildLocalizations(ctx, build, config.Testflight.Localizations); err != nil {
		return err
	}
	if err := client.UpdateBetaLicenseAgreement(ctx, app, config.Testflight); err != nil {
		return err
	}
	if err := client.AssignBetaGroups(ctx, build, config.Testflight.BetaGroups); err != nil {
		return err
	}
	if err := client.AssignBetaTesters(ctx, build, config.Testflight.BetaTesters); err != nil {
		return err
	}
	if config.Versions.ReviewDetails != nil {
		if err := client.UpdateBetaReviewDetails(ctx, app, *config.Testflight.ReviewDetails); err != nil {
			return err
		}
	}
	return nil
}
