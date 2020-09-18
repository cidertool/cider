package client

import (
	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/parallel"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

func (c *ascClient) UpdateApp(ctx *context.Context, appID string, appInfoID string, config config.App) error {
	updateParams := asc.AppUpdateRequestAttributes{
		ContentRightsDeclaration: contentRightsDeclaration(config.UsesThirdPartyContent),
	}

	var availableTerritoryIDs []string
	var prices []asc.NewAppPriceRelationship

	if !ctx.SkipUpdatePricing && config.Availability != nil {
		var err error
		availableTerritoryIDs, err = c.AvailableTerritoryIDsInConfig(ctx, config.Availability.Territories)
		if err != nil {
			return err
		}
		prices = priceSchedules(config.Availability.Pricing)
		updateParams.AvailableInNewTerritories = config.Availability.AvailableInNewTerritories
	}
	if config.PrimaryLocale != "" {
		updateParams.PrimaryLocale = &config.PrimaryLocale
	}

	if _, _, err := c.client.Apps.UpdateApp(ctx, appID, &updateParams, availableTerritoryIDs, prices); err != nil {
		return err
	}

	// TODO: Man, category IDs are kind of wild, aren't they? Will need to fix this at some point.

	// if _, _, err := c.client.Apps.UpdateAppInfo(ctx, appInfoID, &asc.AppInfoUpdateRequestRelationships{
	// 	PrimaryCategoryID:         nil,
	// 	PrimarySubcategoryOneID:   nil,
	// 	PrimarySubcategoryTwoID:   nil,
	// 	SecondaryCategoryID:       nil,
	// 	SecondarySubcategoryOneID: nil,
	// 	SecondarySubcategoryTwoID: nil,
	// }); err != nil {
	// 	return err
	// }
	return nil
}

func (c *ascClient) AvailableTerritoryIDsInConfig(ctx *context.Context, config []string) (availableTerritoryIDs []string, err error) {
	availableTerritoryIDs = make([]string, 0)
	if len(config) == 0 {
		return availableTerritoryIDs, nil
	}
	territoriesResp, _, err := c.client.Pricing.ListTerritories(ctx, &asc.ListTerritoriesQuery{Limit: 200})
	if err != nil {
		return nil, err
	}
	found := make(map[string]bool)
	for _, territory := range territoriesResp.Data {
		found[territory.ID] = true
	}
	for _, id := range config {
		if !found[id] {
			continue
		}
		availableTerritoryIDs = append(availableTerritoryIDs, id)
	}
	return availableTerritoryIDs, nil
}

func contentRightsDeclaration(flag *bool) *string {
	if flag != nil {
		if *flag {
			return asc.String("USES_THIRD_PARTY_CONTENT")
		}
		return asc.String("DOES_NOT_USE_THIRD_PARTY_CONTENT")
	}
	return nil
}

func priceSchedules(schedules []config.PriceSchedule) (priceSchedules []asc.NewAppPriceRelationship) {
	priceSchedules = make([]asc.NewAppPriceRelationship, len(schedules))
	for i, price := range schedules {
		var startDate *asc.Date
		if price.StartDate != nil {
			startDate = &asc.Date{Time: *price.StartDate}
		}
		tier := price.Tier
		priceSchedules[i] = asc.NewAppPriceRelationship{
			StartDate:   startDate,
			PriceTierID: &tier,
		}
	}
	return priceSchedules
}

func (c *ascClient) UpdateAppLocalizations(ctx *context.Context, appID string, config config.AppLocalizations) error {
	var g = parallel.New(ctx.MaxProcesses)
	appInfosResp, _, err := c.client.Apps.ListAppInfosForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	for i := range appInfosResp.Data {
		appInfo := appInfosResp.Data[i]
		if *appInfo.Attributes.AppStoreState != asc.AppStoreVersionStatePrepareForSubmission {
			continue
		}
		appLocResp, _, err := c.client.Apps.ListAppInfoLocalizationsForAppInfo(ctx, appInfo.ID, nil)
		if err != nil {
			return err
		}

		found := make(map[string]bool)
		for i := range appLocResp.Data {
			loc := appLocResp.Data[i]
			locale := *loc.Attributes.Locale
			log.WithField("locale", locale).Debug("found app locale")
			locConfig, ok := config[locale]
			if !ok {
				log.WithField("locale", locale).Debug("not in configuration. skipping...")
				continue
			}
			found[locale] = true

			g.Go(func() error {
				_, _, err := c.client.Apps.UpdateAppInfoLocalization(ctx, loc.ID, &asc.AppInfoLocalizationUpdateRequestAttributes{
					Name:              &locConfig.Name,
					Subtitle:          &locConfig.Subtitle,
					PrivacyPolicyText: &locConfig.PrivacyPolicyText,
					PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
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
				_, _, err := c.client.Apps.CreateAppInfoLocalization(ctx.Context, asc.AppInfoLocalizationCreateRequestAttributes{
					Locale:            locale,
					Name:              &locConfig.Name,
					Subtitle:          &locConfig.Subtitle,
					PrivacyPolicyText: &locConfig.PrivacyPolicyText,
					PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
				}, appInfo.ID)
				return err
			})
		}
	}

	return g.Wait()
}

func (c *ascClient) CreateVersionIfNeeded(ctx *context.Context, appID string, buildID string, config config.Version) (*asc.AppStoreVersion, error) {
	platform, err := config.Platform.APIValue()
	if err != nil {
		return nil, err
	}
	var releaseTypeP *string
	releaseType, _ := config.ReleaseType.APIValue()
	if releaseType != "" {
		releaseTypeP = &releaseType
	}
	var earliestReleaseDate *asc.DateTime
	if config.EarliestReleaseDate != nil {
		earliestReleaseDate = &asc.DateTime{Time: *config.EarliestReleaseDate}
	}
	var versionResp *asc.AppStoreVersionResponse
	versionsResp, _, err := c.client.Apps.ListAppStoreVersionsForApp(ctx, appID, &asc.ListAppStoreVersionsQuery{
		FilterVersionString: []string{ctx.Version},
		FilterPlatform:      []string{string(platform)},
	})
	if err != nil || len(versionsResp.Data) == 0 {
		versionResp, _, err = c.client.Apps.CreateAppStoreVersion(ctx, asc.AppStoreVersionCreateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: earliestReleaseDate,
			Platform:            platform,
			ReleaseType:         releaseTypeP,
			UsesIDFA:            asc.Bool(config.IDFADeclaration != nil),
			VersionString:       ctx.Version,
		}, appID, &buildID)
	} else {
		latestVersion := versionsResp.Data[0]
		versionResp, _, err = c.client.Apps.UpdateAppStoreVersion(ctx, latestVersion.ID, &asc.AppStoreVersionUpdateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: earliestReleaseDate,
			ReleaseType:         releaseTypeP,
			UsesIDFA:            asc.Bool(config.IDFADeclaration != nil),
			VersionString:       &ctx.Version,
		}, &buildID)
	}
	return &versionResp.Data, err
}

func (c *ascClient) UpdateVersionLocalizations(ctx *context.Context, versionID string, config config.VersionLocalizations) error {
	var g = parallel.New(ctx.MaxProcesses)
	locListResp, _, err := c.client.Apps.ListLocalizationsForAppStoreVersion(ctx, versionID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for i := range locListResp.Data {
		loc := locListResp.Data[i]
		locale := *loc.Attributes.Locale
		log.WithField("locale", locale).Debug("found version locale")
		locConfig, ok := config[locale]
		if !ok {
			log.WithField("locale", locale).Debug("not in configuration. skipping...")
			continue
		}
		found[locale] = true

		g.Go(func() error {
			attrs := asc.AppStoreVersionLocalizationUpdateRequestAttributes{
				Description:     &locConfig.Description,
				Keywords:        &locConfig.Keywords,
				MarketingURL:    &locConfig.MarketingURL,
				PromotionalText: &locConfig.PromotionalText,
				SupportURL:      &locConfig.SupportURL,
			}
			// If WhatsNew is set on an app that has never released before, the API will respond with a 409 Conflict when attempting to set the value.
			if !ctx.VersionIsInitialRelease {
				attrs.WhatsNew = &locConfig.WhatsNewText
			}
			log.WithField("locale", locale).Debug("update version locale")
			updatedLocResp, _, err := c.client.Apps.UpdateAppStoreVersionLocalization(ctx, loc.ID, &attrs)
			if err != nil {
				return err
			}
			return c.UpdatePreviewsAndScreenshotsIfNeeded(ctx, g, &updatedLocResp.Data, locConfig)
		})
	}

	for locale := range config {
		locale := locale
		if found[locale] {
			continue
		}
		locConfig := config[locale]

		g.Go(func() error {
			attrs := asc.AppStoreVersionLocalizationCreateRequestAttributes{
				Description:     &locConfig.Description,
				Keywords:        &locConfig.Keywords,
				MarketingURL:    &locConfig.MarketingURL,
				PromotionalText: &locConfig.PromotionalText,
				SupportURL:      &locConfig.SupportURL,
			}
			// If WhatsNew is set on an app that has never released before, the API will respond with a 409 Conflict when attempting to set the value.
			if !ctx.VersionIsInitialRelease {
				attrs.WhatsNew = &locConfig.WhatsNewText
			}
			log.WithField("locale", locale).Debug("create version locale")
			locResp, _, err := c.client.Apps.CreateAppStoreVersionLocalization(ctx.Context, attrs, versionID)
			if err != nil {
				return err
			}
			return c.UpdatePreviewsAndScreenshotsIfNeeded(ctx, g, &locResp.Data, locConfig)
		})
	}

	return g.Wait()
}

func (c *ascClient) UpdateIDFADeclaration(ctx *context.Context, versionID string, config config.IDFADeclaration) error {
	existingDeclResp, _, err := c.client.Submission.GetIDFADeclarationForAppStoreVersion(ctx, versionID, nil)
	if err != nil || existingDeclResp.Data.ID == "" {
		_, _, err = c.client.Submission.CreateIDFADeclaration(ctx, asc.IDFADeclarationCreateRequestAttributes{
			AttributesActionWithPreviousAd:        config.AttributesActionWithPreviousAd,
			AttributesAppInstallationToPreviousAd: config.AttributesAppInstallationToPreviousAd,
			HonorsLimitedAdTracking:               config.HonorsLimitedAdTracking,
			ServesAds:                             config.ServesAds,
		}, versionID)
	} else {
		_, _, err = c.client.Submission.UpdateIDFADeclaration(ctx, existingDeclResp.Data.ID, &asc.IDFADeclarationUpdateRequestAttributes{
			AttributesActionWithPreviousAd:        &config.AttributesActionWithPreviousAd,
			AttributesAppInstallationToPreviousAd: &config.AttributesAppInstallationToPreviousAd,
			HonorsLimitedAdTracking:               &config.HonorsLimitedAdTracking,
			ServesAds:                             &config.ServesAds,
		})
	}
	return err
}

func (c *ascClient) UpdateReviewDetails(ctx *context.Context, versionID string, config config.ReviewDetails) error {
	detailsResp, _, err := c.client.Submission.GetReviewDetailsForAppStoreVersion(ctx, versionID, nil)
	var reviewDetails *asc.AppStoreReviewDetail
	if err == nil {
		reviewDetails, err = c.UpdateReviewDetail(ctx, detailsResp.Data.ID, config)
	} else {
		reviewDetails, err = c.CreateReviewDetail(ctx, versionID, config)
	}
	if err != nil {
		return err
	}

	if err := c.UploadReviewAttachments(ctx, reviewDetails.ID, config.Attachments); err != nil {
		return err
	}

	return nil
}

func (c *ascClient) CreateReviewDetail(ctx *context.Context, versionID string, config config.ReviewDetails) (*asc.AppStoreReviewDetail, error) {
	attributes := asc.AppStoreReviewDetailCreateRequestAttributes{}
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
	resp, _, err := c.client.Submission.CreateReviewDetail(ctx, &attributes, versionID)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *ascClient) UpdateReviewDetail(ctx *context.Context, reviewDetailID string, config config.ReviewDetails) (*asc.AppStoreReviewDetail, error) {
	attributes := asc.AppStoreReviewDetailUpdateRequestAttributes{}
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
	resp, _, err := c.client.Submission.UpdateReviewDetail(ctx, reviewDetailID, &attributes)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *ascClient) EnablePhasedRelease(ctx *context.Context, versionID string) error {
	activePhasedReleaseState := asc.PhasedReleaseStateActive
	phasedResp, _, err := c.client.Publishing.GetAppStoreVersionPhasedReleaseForAppStoreVersion(ctx, versionID, nil)
	if err == nil && phasedResp.Data.ID != "" {
		_, _, err = c.client.Publishing.UpdatePhasedRelease(ctx, phasedResp.Data.ID, &activePhasedReleaseState)
	} else {
		_, _, err = c.client.Publishing.CreatePhasedRelease(ctx, &activePhasedReleaseState, versionID)
	}
	return err
}

func (c *ascClient) SubmitApp(ctx *context.Context, versionID string) error {
	_, _, err := c.client.Submission.CreateSubmission(ctx, versionID)
	return err
}
