// Package testflight is a pipe that processes an app's release to Testflight
package testflight

import (
	"fmt"

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
	return "committing to testflight"
}

// Publish to Testflight.
func (p *Pipe) Publish(ctx *context.Context) error {
	if p.Client == nil {
		p.Client = client.New(ctx)
	}

	for _, name := range ctx.AppsToRelease {
		app, ok := ctx.Config[name]
		if !ok {
			return pipe.ErrMissingApp(name)
		}

		log.WithField("name", name).Info("preparing")

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

	build, err := p.Client.GetBuild(ctx, app)
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
		if err := p.updateBetaDetails(ctx, config, app, build); err != nil {
			return err
		}
	}

	if !ctx.SkipUpdateMetadata || ctx.OverrideBetaGroups {
		if err := p.updateBetaGroups(ctx, config, app, build); err != nil {
			return err
		}
	}

	if !ctx.SkipUpdateMetadata || ctx.OverrideBetaTesters {
		if err := p.updateBetaTesters(ctx, config, app, build); err != nil {
			return err
		}
	}

	if ctx.SkipSubmit {
		return pipe.ErrSkipSubmitEnabled
	}

	log.
		WithField("build", buildVersionLog).
		Info("submitting to testflight")

	return p.Client.SubmitBetaApp(ctx, build.ID)
}

func (p *Pipe) updateBetaDetails(ctx *context.Context, config config.App, app *asc.App, build *asc.Build) error {
	log.Infof("updating %d beta app localizations", len(config.Testflight.Localizations))

	if err := p.Client.UpdateBetaAppLocalizations(ctx, app.ID, config.Testflight.Localizations); err != nil {
		return err
	}

	log.Info("updating beta build details")

	if err := p.Client.UpdateBetaBuildDetails(ctx, build.ID, config.Testflight); err != nil {
		return err
	}

	log.Infof("updating %d beta build localizations", len(config.Testflight.Localizations))

	if err := p.Client.UpdateBetaBuildLocalizations(ctx, build.ID, config.Testflight.Localizations); err != nil {
		return err
	}

	log.Info("updating beta license agreement")

	if err := p.Client.UpdateBetaLicenseAgreement(ctx, app.ID, config.Testflight); err != nil {
		return err
	}

	if config.Testflight.ReviewDetails != nil {
		log.Info("updating beta review details")

		if err := p.Client.UpdateBetaReviewDetails(ctx, app.ID, *config.Testflight.ReviewDetails); err != nil {
			return err
		}
	}

	return nil
}

func (p *Pipe) updateBetaGroups(ctx *context.Context, config config.App, app *asc.App, build *asc.Build) error {
	log.Info("updating build beta groups")
	return p.Client.AssignBetaGroups(ctx, app.ID, build.ID, config.Testflight.BetaGroups)
}

func (p *Pipe) updateBetaTesters(ctx *context.Context, config config.App, app *asc.App, build *asc.Build) error {
	log.Info("updating build beta testers")
	return p.Client.AssignBetaTesters(ctx, app.ID, build.ID, config.Testflight.BetaTesters)
}
