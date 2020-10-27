package client

import (
	"fmt"

	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/parallel"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

// errPlatformNotFound happens when the platform string in the configuration file does not match a supported App Store
// platform.
func errPlatformNotFound(plat config.Platform) error {
	return fmt.Errorf(`platform %s could not be matched up with a supported App Store platform. supported values are "iOS", "macOS", or "tvOS"`, plat)
}

func (c *ascClient) UpdateApp(ctx *context.Context, appID string, appInfoID string, versionID string, config config.App) error {
	var g = parallel.New(ctx.MaxProcesses)

	g.Go(func() error {
		attrs := asc.AppUpdateRequestAttributes{
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
			attrs.AvailableInNewTerritories = config.Availability.AvailableInNewTerritories
		}
		if config.PrimaryLocale != "" {
			attrs.PrimaryLocale = &config.PrimaryLocale
		}

		_, _, err := c.client.Apps.UpdateApp(ctx, appID, &attrs, availableTerritoryIDs, prices)
		return err
	})

	g.Go(func() error {
		if config.Categories == nil {
			return nil
		}

		if _, _, err := c.client.Apps.UpdateAppInfo(ctx, appInfoID, categoriesUpdate(*config.Categories)); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		if config.AgeRatingDeclaration == nil {
			return nil
		}

		ageRatingResp, _, err := c.client.Apps.GetAgeRatingDeclarationForAppStoreVersion(ctx, versionID, nil)
		if err != nil {
			return err
		}
		_, _, err = c.client.Apps.UpdateAgeRatingDeclaration(ctx, ageRatingResp.Data.ID, ageRatingDeclaration(*config.AgeRatingDeclaration))
		return err
	})

	return g.Wait()
}

func (c *ascClient) AvailableTerritoryIDsInConfig(ctx *context.Context, config []string) (availableTerritoryIDs []string, err error) {
	availableTerritoryIDs = make([]string, 0)

	if len(config) == 0 {
		return availableTerritoryIDs, nil
	}

	territoriesLimit := 200

	territoriesResp, _, err := c.client.Pricing.ListTerritories(ctx, &asc.ListTerritoriesQuery{Limit: territoriesLimit})
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

func categoriesUpdate(config config.Categories) *asc.AppInfoUpdateRequestRelationships {
	attrs := asc.AppInfoUpdateRequestRelationships{}

	if config.Primary != "" {
		attrs.PrimaryCategoryID = &config.Primary

		if config.PrimarySubcategories[0] != "" {
			attrs.PrimarySubcategoryOneID = &config.PrimarySubcategories[0]
		}

		if config.PrimarySubcategories[1] != "" {
			attrs.PrimarySubcategoryTwoID = &config.PrimarySubcategories[1]
		}
	}

	if config.Secondary != "" {
		attrs.SecondaryCategoryID = &config.Secondary

		if config.SecondarySubcategories[0] != "" {
			attrs.SecondarySubcategoryOneID = &config.SecondarySubcategories[0]
		}

		if config.SecondarySubcategories[1] != "" {
			attrs.SecondarySubcategoryTwoID = &config.SecondarySubcategories[1]
		}
	}

	return &attrs
}

func ageRatingDeclaration(config config.AgeRatingDeclaration) *asc.AgeRatingDeclarationUpdateRequestAttributes {
	attrs := asc.AgeRatingDeclarationUpdateRequestAttributes{}

	attrs.AlcoholTobaccoOrDrugUseOrReferences = config.AlcoholTobaccoOrDrugUseOrReferences.APIValue()
	attrs.GamblingAndContests = config.GamblingAndContests
	attrs.GamblingSimulated = config.GamblingSimulated.APIValue()
	attrs.HorrorOrFearThemes = config.HorrorOrFearThemes.APIValue()
	attrs.KidsAgeBand = config.KidsAgeBand.APIValue()
	attrs.MatureOrSuggestiveThemes = config.MatureOrSuggestiveThemes.APIValue()
	attrs.MedicalOrTreatmentInformation = config.MedicalOrTreatmentInformation.APIValue()
	attrs.ProfanityOrCrudeHumor = config.ProfanityOrCrudeHumor.APIValue()
	attrs.SexualContentGraphicAndNudity = config.SexualContentGraphicAndNudity.APIValue()
	attrs.SexualContentOrNudity = config.SexualContentOrNudity.APIValue()
	attrs.UnrestrictedWebAccess = config.UnrestrictedWebAccess
	attrs.ViolenceCartoonOrFantasy = config.ViolenceCartoonOrFantasy.APIValue()
	attrs.ViolenceRealistic = config.ViolenceRealistic.APIValue()
	attrs.ViolenceRealisticProlongedGraphicOrSadistic = config.ViolenceRealisticProlongedGraphicOrSadistic.APIValue()

	return &attrs
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
	platform := config.Platform.APIValue()
	if platform == nil {
		return nil, errPlatformNotFound(config.Platform)
	}

	releaseType := config.ReleaseType.APIValue()

	var earliestReleaseDate *asc.DateTime
	if config.EarliestReleaseDate != nil {
		earliestReleaseDate = &asc.DateTime{Time: *config.EarliestReleaseDate}
	}

	var versionResp *asc.AppStoreVersionResponse

	versionsResp, _, err := c.client.Apps.ListAppStoreVersionsForApp(ctx, appID, &asc.ListAppStoreVersionsQuery{
		FilterVersionString: []string{ctx.Version},
		FilterPlatform:      []string{string(*platform)},
	})
	if err != nil || len(versionsResp.Data) == 0 {
		versionResp, _, err = c.client.Apps.CreateAppStoreVersion(ctx, asc.AppStoreVersionCreateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: earliestReleaseDate,
			Platform:            *platform,
			ReleaseType:         releaseType,
			UsesIDFA:            asc.Bool(config.IDFADeclaration != nil),
			VersionString:       ctx.Version,
		}, appID, &buildID)
	} else {
		latestVersion := versionsResp.Data[0]
		versionResp, _, err = c.client.Apps.UpdateAppStoreVersion(ctx, latestVersion.ID, &asc.AppStoreVersionUpdateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: earliestReleaseDate,
			ReleaseType:         releaseType,
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
			log.WithField("locale", locale).Debug("update version locale")
			updatedLocResp, _, err := c.client.Apps.UpdateAppStoreVersionLocalization(ctx, loc.ID, appStoreVersionLocalizationUpdateRequestAttributes(ctx, locConfig))
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
			log.WithField("locale", locale).Debug("create version locale")
			locResp, _, err := c.client.Apps.CreateAppStoreVersionLocalization(ctx.Context, appStoreVersionLocalizationCreateRequestAttributes(ctx, locale, locConfig), versionID)
			if err != nil {
				return err
			}
			return c.UpdatePreviewsAndScreenshotsIfNeeded(ctx, g, &locResp.Data, locConfig)
		})
	}

	return g.Wait()
}

func appStoreVersionLocalizationUpdateRequestAttributes(ctx *context.Context, config config.VersionLocalization) *asc.AppStoreVersionLocalizationUpdateRequestAttributes {
	attrs := asc.AppStoreVersionLocalizationUpdateRequestAttributes{}
	if config.Description != "" {
		attrs.Description = &config.Description
	}

	if config.Keywords != "" {
		attrs.Keywords = &config.Keywords
	}

	if config.MarketingURL != "" {
		attrs.MarketingURL = &config.MarketingURL
	}

	if config.PromotionalText != "" {
		attrs.PromotionalText = &config.PromotionalText
	}

	if config.SupportURL != "" {
		attrs.SupportURL = &config.SupportURL
	}
	// If WhatsNew is set on an app that has never released before, the API will respond with a 409 Conflict when
	// attempting to set the value.
	if !ctx.VersionIsInitialRelease && config.WhatsNewText != "" {
		attrs.WhatsNew = &config.WhatsNewText
	}

	return &attrs
}

func appStoreVersionLocalizationCreateRequestAttributes(ctx *context.Context, locale string, config config.VersionLocalization) asc.AppStoreVersionLocalizationCreateRequestAttributes {
	attrs := asc.AppStoreVersionLocalizationCreateRequestAttributes{
		Locale: locale,
	}

	if config.Description != "" {
		attrs.Description = &config.Description
	}

	if config.Keywords != "" {
		attrs.Keywords = &config.Keywords
	}

	if config.MarketingURL != "" {
		attrs.MarketingURL = &config.MarketingURL
	}

	if config.PromotionalText != "" {
		attrs.PromotionalText = &config.PromotionalText
	}

	if config.SupportURL != "" {
		attrs.SupportURL = &config.SupportURL
	}

	// If WhatsNew is set on an app that has never released before, the API will respond with a 409 Conflict when
	// attempting to set the value.
	if !ctx.VersionIsInitialRelease && config.WhatsNewText != "" {
		attrs.WhatsNew = &config.WhatsNewText
	}

	return attrs
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
	var reviewDetails *asc.AppStoreReviewDetail

	detailsResp, _, err := c.client.Submission.GetReviewDetailsForAppStoreVersion(ctx, versionID, nil)
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
	attrs := asc.AppStoreReviewDetailCreateRequestAttributes{}

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

	resp, _, err := c.client.Submission.CreateReviewDetail(ctx, &attrs, versionID)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *ascClient) UpdateReviewDetail(ctx *context.Context, reviewDetailID string, config config.ReviewDetails) (*asc.AppStoreReviewDetail, error) {
	attrs := asc.AppStoreReviewDetailUpdateRequestAttributes{}

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

	resp, _, err := c.client.Submission.UpdateReviewDetail(ctx, reviewDetailID, &attrs)
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
