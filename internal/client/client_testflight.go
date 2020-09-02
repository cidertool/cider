package client

import (
	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

func (c *ascClient) UpdateBetaAppLocalizations(ctx *context.Context, app *asc.App, config config.TestflightLocalizations) error {
	locListResp, _, err := c.client.TestFlight.ListBetaAppLocalizationsForApp(ctx, app.ID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for _, loc := range locListResp.Data {
		locale := *loc.Attributes.Locale
		found[locale] = true
		locConfig := config[locale]
		_, _, err := c.client.TestFlight.UpdateBetaAppLocalization(ctx, loc.ID, &asc.BetaAppLocalizationUpdateRequestAttributes{
			Description:       &locConfig.Description,
			FeedbackEmail:     &locConfig.FeedbackEmail,
			MarketingURL:      &locConfig.MarketingURL,
			PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
			TVOSPrivacyPolicy: &locConfig.TVOSPrivacyPolicy,
		})
		if err != nil {
			return err
		}
	}

	for locale, locConfig := range config {
		if found[locale] {
			continue
		}
		_, _, err := c.client.TestFlight.CreateBetaAppLocalization(ctx.Context, asc.BetaAppLocalizationCreateRequestAttributes{
			Description:       &locConfig.Description,
			FeedbackEmail:     &locConfig.FeedbackEmail,
			Locale:            locale,
			MarketingURL:      &locConfig.MarketingURL,
			PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
			TVOSPrivacyPolicy: &locConfig.TVOSPrivacyPolicy,
		}, app.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config config.TestflightForApp) error {
	_, _, err := c.client.TestFlight.UpdateBuildBetaDetail(ctx, build.ID, &config.EnableAutoNotify)
	return err
}

func (c *ascClient) UpdateBetaBuildLocalizations(ctx *context.Context, build *asc.Build, config config.TestflightLocalizations) error {
	locListResp, _, err := c.client.TestFlight.ListBetaBuildLocalizationsForBuild(ctx, build.ID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for _, loc := range locListResp.Data {
		locale := *loc.Attributes.Locale
		found[locale] = true
		locConfig := config[locale]
		_, _, err := c.client.TestFlight.UpdateBetaBuildLocalization(ctx, loc.ID, &locConfig.WhatsNew)
		if err != nil {
			return err
		}
	}

	for locale, locConfig := range config {
		if found[locale] {
			continue
		}
		_, _, err := c.client.TestFlight.CreateBetaBuildLocalization(ctx.Context, locale, &locConfig.WhatsNew, build.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UpdateBetaLicenseAgreement(ctx *context.Context, app *asc.App, config config.TestflightForApp) error {
	resp, _, err := c.client.TestFlight.GetBetaLicenseAgreementForApp(ctx, app.ID, nil)
	if err != nil {
		return err
	}
	_, _, err = c.client.TestFlight.UpdateBetaLicenseAgreement(ctx, resp.Data.ID, &config.LicenseAgreement)
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

func (c *ascClient) UpdateBetaReviewDetails(ctx *context.Context, app *asc.App, config config.ReviewDetails) error {
	detailsResp, _, err := c.client.TestFlight.GetBetaAppReviewDetailsForApp(ctx, app.ID, nil)
	if err != nil {
		return err
	}
	_, _, err = c.client.TestFlight.UpdateBetaAppReviewDetail(ctx, detailsResp.Data.ID, &asc.BetaAppReviewDetailUpdateRequestAttributes{
		ContactEmail:        &config.Contact.Email,
		ContactFirstName:    &config.Contact.FirstName,
		ContactLastName:     &config.Contact.LastName,
		ContactPhone:        &config.Contact.Phone,
		DemoAccountName:     &config.DemoAccount.Name,
		DemoAccountPassword: &config.DemoAccount.Password,
		DemoAccountRequired: &config.DemoAccount.Required,
		Notes:               &config.Notes,
	})
	return err
}

func (c *ascClient) SubmitBetaApp(ctx *context.Context, build *asc.Build) error {
	_, _, err := c.client.TestFlight.CreateBetaAppReviewSubmission(ctx, build.ID)
	return err
}
