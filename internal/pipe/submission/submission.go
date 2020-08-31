package submission

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
	version, err := client.CreateVersionIfNeeded(ctx, app, build, &appConfig.Versions)
	if err != nil {
		return err
	}
	err = client.UpdateAppLocalizations(ctx, app, appConfig.Localizations)
	if err != nil {
		return err
	}
	err = client.UpdateVersionLocalizations(ctx, version, appConfig.Versions.Localizations)
	if err != nil {
		return err
	}
	err = client.UpdateIDFADeclaration(ctx, version, appConfig.Versions.IDFADeclaration)
	if err != nil {
		return err
	}
	err = client.UploadRoutingCoverage(ctx, version, appConfig.Versions.RoutingCoverage)
	if err != nil {
		return err
	}
	err = client.UpdateReviewDetails(ctx, version, appConfig.Versions.ReviewDetails)
	if err != nil {
		return err
	}
	return client.SubmitApp(ctx, version)
}
