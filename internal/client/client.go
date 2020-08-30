package client

import (
	"fmt"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

type Client interface {
	GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error)
	GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error)
	UpdateBetaAppLocalizations(ctx *context.Context, app *asc.App, config *config.TestflightLocalizations) error
	UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config *config.TestflightForApp) error
	UpdateBetaBuildLocalizations(ctx *context.Context, build *asc.Build, config *config.TestflightLocalizations) error
	SubmitBetaApp(ctx *context.Context, build *asc.Build) error
}

func New(ctx *context.Context) Client {
	client := asc.NewClient(ctx.Credentials.Client())
	return &ascClient{client: client}
}

type ascClient struct {
	client *asc.Client
}

func (c *ascClient) GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error) {
	resp, _, err := c.client.Apps.ListApps(ctx, &asc.ListAppsQuery{
		FilterBundleID: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("app not found matching %s: %w", id, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("app not found matching %s", id)
	}
	return &resp.Data[0], nil
}

func (c *ascClient) GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error) {
	resp, _, err := c.client.Builds.ListBuilds(ctx, &asc.ListBuildsQuery{
		FilterApp:     []string{app.ID},
		FilterVersion: []string{ctx.Version},
	})
	if err != nil {
		return nil, fmt.Errorf("build not found matching app %s and version %s: %w", *app.Attributes.BundleID, ctx.Version, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("build not found matching app %s and version %s", *app.Attributes.BundleID, ctx.Version)
	}
	return &resp.Data[0], nil
}

func (c *ascClient) UpdateBetaAppLocalizations(ctx *context.Context, app *asc.App, config *config.TestflightLocalizations) error {
	locListResp, _, err := c.client.TestFlight.ListBetaAppLocalizationsForApp(ctx, app.ID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for _, loc := range locListResp.Data {
		locale := *loc.Attributes.Locale
		found[locale] = true
		locConfig := (*config)[locale]
		_, _, err := c.client.TestFlight.UpdateBetaAppLocalization(ctx, loc.ID, &asc.BetaAppLocalizationUpdateRequestAttributes{
			Description:       asc.String(locConfig.Description),
			FeedbackEmail:     asc.String(locConfig.FeedbackEmail),
			MarketingURL:      asc.String(locConfig.MarketingURL),
			PrivacyPolicyURL:  asc.String(locConfig.PrivacyPolicyURL),
			TVOSPrivacyPolicy: asc.String(locConfig.TVOSPrivacyPolicy),
		})
		if err != nil {
			return err
		}
	}

	for locale, locConfig := range *config {
		if found[locale] {
			continue
		}
		_, _, err := c.client.TestFlight.CreateBetaAppLocalization(ctx.Context, asc.BetaAppLocalizationCreateRequestAttributes{
			Description:       asc.String(locConfig.Description),
			FeedbackEmail:     asc.String(locConfig.FeedbackEmail),
			Locale:            locale,
			MarketingURL:      asc.String(locConfig.MarketingURL),
			PrivacyPolicyURL:  asc.String(locConfig.PrivacyPolicyURL),
			TVOSPrivacyPolicy: asc.String(locConfig.TVOSPrivacyPolicy),
		}, app.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config *config.TestflightForApp) error {
	_, _, err := c.client.TestFlight.UpdateBuildBetaDetail(ctx, build.ID, asc.Bool(config.EnableAutoNotify))
	return err
}

func (c *ascClient) UpdateBetaBuildLocalizations(ctx *context.Context, build *asc.Build, config *config.TestflightLocalizations) error {
	locListResp, _, err := c.client.TestFlight.ListBetaBuildLocalizationsForBuild(ctx, build.ID, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ascClient) SubmitBetaApp(ctx *context.Context, build *asc.Build) error {
	_, _, err := c.client.TestFlight.CreateBetaAppReviewSubmission(ctx, build.ID)
	return err
}
