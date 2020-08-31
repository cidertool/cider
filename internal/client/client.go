package client

import (
	"fmt"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

// Client is an abstraction of an App Store Connect API client's functionality
type Client interface {
	GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error)
	GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error)
	UpdateBetaAppLocalizations(ctx *context.Context, app *asc.App, config *config.TestflightLocalizations) error
	UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config *config.TestflightForApp) error
	UpdateBetaBuildLocalizations(ctx *context.Context, build *asc.Build, config *config.TestflightLocalizations) error
	UpdateBetaLicenseAgreement(ctx *context.Context, app *asc.App, config *config.TestflightForApp) error
	AssignBetaGroups(ctx *context.Context, build *asc.Build, groups []string) error
	AssignBetaTesters(ctx *context.Context, build *asc.Build, testers []config.BetaTester) error
	SubmitBetaApp(ctx *context.Context, build *asc.Build) error
}

// New returns a new Client
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

	found := make(map[string]bool)
	for _, loc := range locListResp.Data {
		locale := *loc.Attributes.Locale
		found[locale] = true
		locConfig := (*config)[locale]
		_, _, err := c.client.TestFlight.UpdateBetaBuildLocalization(ctx, loc.ID, asc.String(locConfig.WhatsNew))
		if err != nil {
			return err
		}
	}

	for locale, locConfig := range *config {
		if found[locale] {
			continue
		}
		_, _, err := c.client.TestFlight.CreateBetaBuildLocalization(ctx.Context, locale, asc.String(locConfig.WhatsNew), build.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UpdateBetaLicenseAgreement(ctx *context.Context, app *asc.App, config *config.TestflightForApp) error {
	resp, _, err := c.client.TestFlight.GetBetaLicenseAgreementForApp(ctx, app.ID, nil)
	if err != nil {
		return err
	}
	_, _, err = c.client.TestFlight.UpdateBetaLicenseAgreement(ctx, resp.Data.ID, asc.String(config.LicenseAgreement))
	return err
}

func (c *ascClient) AssignBetaGroups(ctx *context.Context, build *asc.Build, groups []string) error {
	groupsResp, _, err := c.client.TestFlight.ListBetaGroups(ctx, &asc.ListBetaGroupsQuery{
		FilterName: groups,
	})
	if err != nil {
		return err
	}
	for _, group := range groupsResp.Data {
		_, err := c.client.TestFlight.AddBuildsToBetaGroup(ctx, group.ID, []string{build.ID})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) AssignBetaTesters(ctx *context.Context, build *asc.Build, testers []config.BetaTester) error {
	emailFilters := make([]string, 0)
	firstNameFilters := make([]string, 0)
	lastNameFilters := make([]string, 0)
	for _, tester := range testers {
		if tester.Email != "" {
			emailFilters = append(emailFilters, tester.Email)
		}
		if tester.FirstName != "" {
			firstNameFilters = append(firstNameFilters, tester.FirstName)
		}
		if tester.LastName != "" {
			lastNameFilters = append(lastNameFilters, tester.LastName)
		}
	}

	testersResp, _, err := c.client.TestFlight.ListBetaTesters(ctx, &asc.ListBetaTestersQuery{
		FilterEmail:     emailFilters,
		FilterFirstName: firstNameFilters,
		FilterLastName:  lastNameFilters,
	})
	if err != nil {
		return err
	}
	for _, tester := range testersResp.Data {
		_, err := c.client.TestFlight.AssignSingleBetaTesterToBuilds(ctx, tester.ID, []string{build.ID})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) SubmitBetaApp(ctx *context.Context, build *asc.Build) error {
	_, _, err := c.client.TestFlight.CreateBetaAppReviewSubmission(ctx, build.ID)
	return err
}
