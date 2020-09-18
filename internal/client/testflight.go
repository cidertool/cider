package client

import (
	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/parallel"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

func (c *ascClient) UpdateBetaAppLocalizations(ctx *context.Context, appID string, config config.TestflightLocalizations) error {
	var g = parallel.New(ctx.MaxProcesses)
	locListResp, _, err := c.client.TestFlight.ListBetaAppLocalizationsForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for i := range locListResp.Data {
		loc := locListResp.Data[i]
		locale := *loc.Attributes.Locale
		log.WithField("locale", locale).Debug("found beta app locale")
		locConfig, ok := config[locale]
		if !ok {
			log.WithField("locale", locale).Debug("not in configuration. skipping...")
			continue
		}
		found[locale] = true

		g.Go(func() error {
			_, _, err = c.client.TestFlight.UpdateBetaAppLocalization(ctx, loc.ID, &asc.BetaAppLocalizationUpdateRequestAttributes{
				Description:       &locConfig.Description,
				FeedbackEmail:     &locConfig.FeedbackEmail,
				MarketingURL:      &locConfig.MarketingURL,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
				TVOSPrivacyPolicy: &locConfig.TVOSPrivacyPolicy,
			})
			return err
		})
	}

	for locale := range config {
		locale := locale
		if found[locale] {
			continue
		}
		locConfig := config[locale]

		g.Go(func() error {
			_, _, err = c.client.TestFlight.CreateBetaAppLocalization(ctx.Context, asc.BetaAppLocalizationCreateRequestAttributes{
				Description:       &locConfig.Description,
				FeedbackEmail:     &locConfig.FeedbackEmail,
				Locale:            locale,
				MarketingURL:      &locConfig.MarketingURL,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
				TVOSPrivacyPolicy: &locConfig.TVOSPrivacyPolicy,
			}, appID)
			return err
		})
	}

	return g.Wait()
}

func (c *ascClient) UpdateBetaBuildDetails(ctx *context.Context, buildID string, config config.TestflightForApp) error {
	_, _, err := c.client.TestFlight.UpdateBuildBetaDetail(ctx, buildID, &config.EnableAutoNotify)
	return err
}

func (c *ascClient) UpdateBetaBuildLocalizations(ctx *context.Context, buildID string, config config.TestflightLocalizations) error {
	var g = parallel.New(ctx.MaxProcesses)
	locListResp, _, err := c.client.TestFlight.ListBetaBuildLocalizationsForBuild(ctx, buildID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for i := range locListResp.Data {
		loc := locListResp.Data[i]
		locale := *loc.Attributes.Locale
		log.WithField("locale", locale).Debug("found beta build locale")
		locConfig, ok := config[locale]
		if !ok {
			log.WithField("locale", locale).Debug("not in configuration. skipping...")
			continue
		}
		found[locale] = true

		g.Go(func() error {
			_, _, err := c.client.TestFlight.UpdateBetaBuildLocalization(ctx, loc.ID, &locConfig.WhatsNew)
			return err
		})
	}

	for locale := range config {
		locale := locale
		if found[locale] {
			continue
		}
		locConfig := config[locale]

		g.Go(func() error {
			_, _, err := c.client.TestFlight.CreateBetaBuildLocalization(ctx.Context, locale, &locConfig.WhatsNew, buildID)
			return err
		})
	}

	return g.Wait()
}

func (c *ascClient) UpdateBetaLicenseAgreement(ctx *context.Context, appID string, config config.TestflightForApp) error {
	resp, _, err := c.client.TestFlight.GetBetaLicenseAgreementForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	_, _, err = c.client.TestFlight.UpdateBetaLicenseAgreement(ctx, resp.Data.ID, &config.LicenseAgreement)
	return err
}

func (c *ascClient) AssignBetaGroups(ctx *context.Context, appID string, buildID string, groups []string) error {
	if len(groups) == 0 {
		log.Debug("no groups provided as input to add")
		return nil
	}
	groupsResp, _, err := c.client.TestFlight.ListBetaGroups(ctx, &asc.ListBetaGroupsQuery{
		FilterName: groups,
		FilterApp:  []string{appID},
	})
	if err != nil {
		return err
	}
	if len(groupsResp.Data) == 0 {
		log.WithField("groups", groups).Warn("no matching groups found")
	}
	for _, group := range groupsResp.Data {
		_, err := c.client.TestFlight.AddBuildsToBetaGroup(ctx, group.ID, []string{buildID})
		if err != nil {
			return err
		}
	}
	log.Infof("assigned %d beta group(s)", len(groupsResp.Data))
	return nil
}

func (c *ascClient) AssignBetaTesters(ctx *context.Context, appID string, buildID string, testers []config.BetaTester) error {
	if len(testers) == 0 {
		log.Debug("no testers provided as input to add")
		return nil
	}
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
		FilterApps:      []string{appID},
	})
	if err != nil {
		return err
	}
	if len(testersResp.Data) == 0 {
		log.WithFields(log.Fields{
			"emails":     emailFilters,
			"firstNames": firstNameFilters,
			"lastNames":  lastNameFilters,
		}).Warn("no matching testers found")
	}
	for _, tester := range testersResp.Data {
		_, err := c.client.TestFlight.AssignSingleBetaTesterToBuilds(ctx, tester.ID, []string{buildID})
		if err != nil {
			return err
		}
	}
	log.Infof("assigned %d beta tester(s)", len(testersResp.Data))
	return nil
}

func (c *ascClient) UpdateBetaReviewDetails(ctx *context.Context, appID string, config config.ReviewDetails) error {
	detailsResp, _, err := c.client.TestFlight.GetBetaAppReviewDetailsForApp(ctx, appID, nil)
	if err != nil {
		return err
	}
	attributes := asc.BetaAppReviewDetailUpdateRequestAttributes{}
	if config.Contact != nil {
		attributes.ContactEmail = &config.Contact.Email
		attributes.ContactFirstName = &config.Contact.FirstName
		attributes.ContactLastName = &config.Contact.LastName
		attributes.ContactPhone = &config.Contact.Phone
	}
	if config.DemoAccount != nil {
		attributes.DemoAccountName = &config.DemoAccount.Name
		attributes.DemoAccountPassword = &config.DemoAccount.Password
		attributes.DemoAccountRequired = &config.DemoAccount.Required
	}
	attributes.Notes = &config.Notes
	_, _, err = c.client.TestFlight.UpdateBetaAppReviewDetail(ctx, detailsResp.Data.ID, &attributes)
	return err
}

func (c *ascClient) SubmitBetaApp(ctx *context.Context, buildID string) error {
	_, _, err := c.client.TestFlight.CreateBetaAppReviewSubmission(ctx, buildID)
	return err
}
