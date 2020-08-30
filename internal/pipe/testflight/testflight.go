package testflight

import (
	"github.com/aaronsky/applereleaser/internal/client"
	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/apex/log"
)

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "choosing processed build"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	client := client.New(ctx)
	for name, app := range ctx.Config.Apps {
		log.WithField("testflight", name).Info("updating metadata")
		err := doRelease(ctx, app, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func doRelease(ctx *context.Context, appConfig config.App, client client.Client) error {
	app, err := client.GetAppForBundleID(ctx, appConfig.BundleID)
	if err != nil {
		return err
	}
	build, err := client.GetRelevantBuild(ctx, app)
	if err != nil {
		return err
	}
	err = client.UpdateBetaAppLocalizations(ctx, app, &appConfig.Testflight.Localizations)
	if err != nil {
		return err
	}
	err = client.UpdateBetaBuildDetails(ctx, build, &appConfig.Testflight)
	if err != nil {
		return err
	}
	err = client.UpdateBetaBuildLocalizations(ctx, build, &appConfig.Testflight.Localizations)
	if err != nil {
		return err
	}
	// 5. create or update beta build localizations and assign to build
	// 6. create or update beta license agreements
	// 7. assign beta groups and beta testers to build
	// 8. create beta app submission for build
	return client.SubmitBetaApp(ctx, build)
}
