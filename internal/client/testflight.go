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
			_, _, err = c.client.TestFlight.UpdateBetaAppLocalization(ctx, loc.ID, betaAppLocalizationUpdateRequestAttributes(locConfig))
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
			_, _, err = c.client.TestFlight.CreateBetaAppLocalization(ctx.Context, betaAppLocalizationCreateRequestAttributes(locale, locConfig), appID)
			return err
		})
	}

	return g.Wait()
}

func betaAppLocalizationUpdateRequestAttributes(locConfig config.TestflightLocalization) *asc.BetaAppLocalizationUpdateRequestAttributes {
	attrs := asc.BetaAppLocalizationUpdateRequestAttributes{}

	if locConfig.Description != "" {
		attrs.Description = &locConfig.Description
	}

	if locConfig.FeedbackEmail != "" {
		attrs.FeedbackEmail = &locConfig.FeedbackEmail
	}

	if locConfig.MarketingURL != "" {
		attrs.MarketingURL = &locConfig.MarketingURL
	}

	if locConfig.PrivacyPolicyURL != "" {
		attrs.PrivacyPolicyURL = &locConfig.PrivacyPolicyURL
	}

	if locConfig.TVOSPrivacyPolicy != "" {
		attrs.TVOSPrivacyPolicy = &locConfig.TVOSPrivacyPolicy
	}

	return &attrs
}

func betaAppLocalizationCreateRequestAttributes(locale string, locConfig config.TestflightLocalization) asc.BetaAppLocalizationCreateRequestAttributes {
	attrs := asc.BetaAppLocalizationCreateRequestAttributes{
		Locale: locale,
	}

	if locConfig.Description != "" {
		attrs.Description = &locConfig.Description
	}

	if locConfig.FeedbackEmail != "" {
		attrs.FeedbackEmail = &locConfig.FeedbackEmail
	}

	if locConfig.MarketingURL != "" {
		attrs.MarketingURL = &locConfig.MarketingURL
	}

	if locConfig.PrivacyPolicyURL != "" {
		attrs.PrivacyPolicyURL = &locConfig.PrivacyPolicyURL
	}

	if locConfig.TVOSPrivacyPolicy != "" {
		attrs.TVOSPrivacyPolicy = &locConfig.TVOSPrivacyPolicy
	}

	return attrs
}

func (c *ascClient) UpdateBetaBuildDetails(ctx *context.Context, buildID string, config config.Testflight) error {
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
		} else if locConfig.WhatsNew == "" {
			log.WithField("locale", locale).Warn("skipping updating beta build localization due to empty What's New text")
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

		if locConfig.WhatsNew == "" {
			log.WithField("locale", locale).Warn("skipping updating beta build localization due to empty What's New text")
			continue
		}

		g.Go(func() error {
			_, _, err := c.client.TestFlight.CreateBetaBuildLocalization(ctx.Context, locale, &locConfig.WhatsNew, buildID)
			return err
		})
	}

	return g.Wait()
}

func (c *ascClient) UpdateBetaLicenseAgreement(ctx *context.Context, appID string, config config.Testflight) error {
	if config.LicenseAgreement == "" {
		return nil
	}

	resp, _, err := c.client.TestFlight.GetBetaLicenseAgreementForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	_, _, err = c.client.TestFlight.UpdateBetaLicenseAgreement(ctx, resp.Data.ID, &config.LicenseAgreement)

	return err
}

func (c *ascClient) AssignBetaGroups(ctx *context.Context, appID string, buildID string, groups []config.BetaGroup) error {
	var g = parallel.New(ctx.MaxProcesses)

	if len(groups) == 0 {
		log.Debug("no groups in configuration")
		return nil
	}

	existingGroupsResp, _, err := c.client.TestFlight.ListBetaGroups(ctx, &asc.ListBetaGroupsQuery{
		FilterApp: []string{appID},
	})
	if err != nil {
		return err
	}

	// Map of group names -> config.BetaGroup
	var groupConfigs = make(map[string]config.BetaGroup, len(groups))

	for i := range groups {
		group := groups[i]
		groupConfigs[group.Name] = group
	}

	// Map of group names -> whether or not they exist in the configuration
	var found = make(map[string]bool)

	for i := range existingGroupsResp.Data {
		group := existingGroupsResp.Data[i]
		if group.Attributes == nil || group.Attributes.Name == nil {
			continue
		}

		name := *group.Attributes.Name
		found[name] = true
		configGroup, ok := groupConfigs[name]

		if !ok {
			log.WithField("group", name).Debug("not in configuration. skipping...")
			continue
		}

		g.Go(func() error {
			log.WithField("group", name).Debug("update beta group")
			return c.updateBetaGroup(ctx, g, appID, group.ID, buildID, configGroup)
		})
	}

	for i := range groups {
		group := groups[i]
		if group.Name == "" {
			log.Warn("skipping a beta group with a missing name")
			continue
		} else if found[group.Name] {
			continue
		}

		g.Go(func() error {
			log.WithField("group", group.Name).Debug("create beta group")
			return c.createBetaGroup(ctx, g, appID, buildID, group)
		})
	}

	return g.Wait()
}

func (c *ascClient) updateBetaGroup(ctx *context.Context, g parallel.Group, appID string, groupID string, buildID string, group config.BetaGroup) error {
	g.Go(func() error {
		_, _, err := c.client.TestFlight.UpdateBetaGroup(ctx, groupID, &asc.BetaGroupUpdateRequestAttributes{
			FeedbackEnabled:        &group.FeedbackEnabled,
			Name:                   &group.Name,
			PublicLinkEnabled:      &group.EnablePublicLink,
			PublicLinkLimit:        &group.PublicLinkLimit,
			PublicLinkLimitEnabled: &group.EnablePublicLinkLimit,
		})
		return err
	})
	g.Go(func() error {
		_, err := c.client.TestFlight.AddBuildsToBetaGroup(ctx, groupID, []string{buildID})
		return err
	})
	g.Go(func() error {
		return c.updateBetaTestersForGroup(ctx, g, appID, groupID, group.Testers)
	})

	return nil
}

func (c *ascClient) createBetaGroup(ctx *context.Context, g parallel.Group, appID string, buildID string, group config.BetaGroup) error {
	newGroupResp, _, err := c.client.TestFlight.CreateBetaGroup(ctx, asc.BetaGroupCreateRequestAttributes{
		FeedbackEnabled:        &group.FeedbackEnabled,
		Name:                   group.Name,
		PublicLinkEnabled:      &group.EnablePublicLink,
		PublicLinkLimit:        &group.PublicLinkLimit,
		PublicLinkLimitEnabled: &group.EnablePublicLinkLimit,
	}, appID, nil, []string{buildID})
	if err != nil {
		return err
	}

	g.Go(func() error {
		return c.updateBetaTestersForGroup(ctx, g, appID, newGroupResp.Data.ID, group.Testers)
	})

	return nil
}

func (c *ascClient) updateBetaTestersForGroup(ctx *context.Context, g parallel.Group, appID string, groupID string, testers []config.BetaTester) error {
	if len(testers) == 0 {
		return nil
	}

	existingTesters, err := c.listBetaTesters(ctx, appID, testers)
	if err != nil {
		return err
	}

	betaTesterIDs, found := filterTestersNotInBetaGroup(existingTesters, groupID)

	g.Go(func() error {
		_, err = c.client.TestFlight.AddBetaTestersToBetaGroup(ctx, groupID, betaTesterIDs)
		return err
	})

	for i := range testers {
		tester := testers[i]
		if tester.Email == "" {
			log.Warnf("skipping a beta tester in beta group %s with a missing email", groupID)
			continue
		} else if found[tester.Email] {
			continue
		}

		g.Go(func() error {
			return c.createBetaTester(ctx, tester, []string{groupID}, nil)
		})
	}

	return err
}

func filterTestersNotInBetaGroup(testers []asc.BetaTester, groupID string) (betaTesterIDs []string, found map[string]bool) {
	betaTesterIDs = make([]string, 0)
	found = make(map[string]bool)

	for i := range testers {
		tester := testers[i]
		if tester.Attributes == nil || tester.Attributes.Email == nil {
			continue
		}

		email := string(*tester.Attributes.Email)
		found[email] = true
		inGroup := false

		if tester.Relationships != nil &&
			tester.Relationships.BetaGroups != nil &&
			len(tester.Relationships.BetaGroups.Data) > 0 {
			for _, rel := range tester.Relationships.BetaGroups.Data {
				if rel.ID == groupID {
					inGroup = true
					break
				}
			}
		}

		if !inGroup {
			betaTesterIDs = append(betaTesterIDs, tester.ID)
		}
	}

	return betaTesterIDs, found
}

func (c *ascClient) AssignBetaTesters(ctx *context.Context, appID string, buildID string, testers []config.BetaTester) error {
	var g = parallel.New(ctx.MaxProcesses)

	if len(testers) == 0 {
		return nil
	}

	existingTesters, err := c.listBetaTesters(ctx, appID, testers)
	if err != nil {
		return err
	}
	// Map of tester emails -> whether or not they exist in the configuration
	var found = make(map[string]bool)

	for i := range existingTesters {
		tester := existingTesters[i]
		if tester.Attributes == nil || tester.Attributes.Email == nil {
			continue
		}

		email := string(*tester.Attributes.Email)
		found[email] = true

		g.Go(func() error {
			log.
				WithFields(log.Fields{
					"email": email,
					"build": buildID,
				}).
				Debug("assign individual beta tester")
			_, err := c.client.TestFlight.AssignSingleBetaTesterToBuilds(ctx, tester.ID, []string{buildID})
			return err
		})
	}

	for i := range testers {
		tester := testers[i]
		if tester.Email == "" {
			log.Warn("beta tester email missing")
			continue
		} else if found[tester.Email] {
			continue
		}

		g.Go(func() error {
			log.WithField("email", tester.Email).Debug("create individual beta tester")
			return c.createBetaTester(ctx, tester, nil, []string{buildID})
		})
	}

	return g.Wait()
}

func (c *ascClient) listBetaTesters(ctx *context.Context, appID string, config []config.BetaTester) ([]asc.BetaTester, error) {
	emailFilters := make([]string, 0)
	firstNameFilters := make([]string, 0)
	lastNameFilters := make([]string, 0)

	for _, tester := range config {
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
		return nil, err
	}

	return testersResp.Data, nil
}

func (c *ascClient) createBetaTester(ctx *context.Context, tester config.BetaTester, betaGroupIDs []string, buildIDs []string) error {
	if betaGroupIDs != nil && buildIDs != nil {
		log.WithField("tester", tester).Warn("authors note: you can't define betaGroupIDs and buildIDs at the same time")
	}

	_, _, err := c.client.TestFlight.CreateBetaTester(ctx, asc.BetaTesterCreateRequestAttributes{
		Email:     asc.Email(tester.Email),
		FirstName: &tester.FirstName,
		LastName:  &tester.LastName,
	}, betaGroupIDs, buildIDs)

	return err
}

func (c *ascClient) UpdateBetaReviewDetails(ctx *context.Context, appID string, config config.ReviewDetails) error {
	detailsResp, _, err := c.client.TestFlight.GetBetaAppReviewDetailsForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	attrs := asc.BetaAppReviewDetailUpdateRequestAttributes{}

	if config.Contact != nil {
		attrs.ContactEmail = &config.Contact.Email
		attrs.ContactFirstName = &config.Contact.FirstName
		attrs.ContactLastName = &config.Contact.LastName
		attrs.ContactPhone = &config.Contact.Phone
	}

	if config.DemoAccount != nil {
		attrs.DemoAccountName = &config.DemoAccount.Name
		attrs.DemoAccountPassword = &config.DemoAccount.Password
		attrs.DemoAccountRequired = &config.DemoAccount.Required
	}

	attrs.Notes = &config.Notes

	if len(config.Attachments) > 0 {
		log.Warn("attachments are not supported for beta review details and will be ignored")
	}

	_, _, err = c.client.TestFlight.UpdateBetaAppReviewDetail(ctx, detailsResp.Data.ID, &attrs)

	return err
}

func (c *ascClient) SubmitBetaApp(ctx *context.Context, buildID string) error {
	_, _, err := c.client.TestFlight.CreateBetaAppReviewSubmission(ctx, buildID)
	return err
}
