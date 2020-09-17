package client

import (
	"os"
	"path/filepath"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/closer"
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
		if found[id] {
			availableTerritoryIDs = append(availableTerritoryIDs, id)
		}
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
	appInfosResp, _, err := c.client.Apps.ListAppInfosForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	for _, appInfo := range appInfosResp.Data {
		if *appInfo.Attributes.AppStoreState != asc.AppStoreVersionStatePrepareForSubmission {
			continue
		}
		appLocResp, _, err := c.client.Apps.ListAppInfoLocalizationsForAppInfo(ctx, appInfo.ID, nil)
		if err != nil {
			return err
		}

		found := make(map[string]bool)
		for _, loc := range appLocResp.Data {
			locale := *loc.Attributes.Locale
			locConfig, ok := config[locale]
			if !ok {
				continue
			}
			found[locale] = true

			if _, _, err := c.client.Apps.UpdateAppInfoLocalization(ctx, loc.ID, &asc.AppInfoLocalizationUpdateRequestAttributes{
				Name:              &locConfig.Name,
				Subtitle:          &locConfig.Subtitle,
				PrivacyPolicyText: &locConfig.PrivacyPolicyText,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
			}); err != nil {
				return err
			}
		}

		for locale, locConfig := range config {
			if found[locale] {
				continue
			}

			if _, _, err := c.client.Apps.CreateAppInfoLocalization(ctx.Context, asc.AppInfoLocalizationCreateRequestAttributes{
				Locale:            locale,
				Name:              &locConfig.Name,
				Subtitle:          &locConfig.Subtitle,
				PrivacyPolicyText: &locConfig.PrivacyPolicyText,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
			}, appInfo.ID); err != nil {
				return err
			}
		}
	}

	return nil
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
	locListResp, _, err := c.client.Apps.ListLocalizationsForAppStoreVersion(ctx, versionID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for _, loc := range locListResp.Data {
		locale := *loc.Attributes.Locale
		locConfig, ok := config[locale]
		if !ok {
			continue
		}
		found[locale] = true

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
		updatedLocResp, _, err := c.client.Apps.UpdateAppStoreVersionLocalization(ctx, loc.ID, &attrs)
		if err != nil {
			return err
		}
		if err := c.UpdatePreviewsAndScreenshotsIfNeeded(ctx, &updatedLocResp.Data, locConfig); err != nil {
			return err
		}
	}

	for locale, locConfig := range config {
		if found[locale] {
			continue
		}

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
		locResp, _, err := c.client.Apps.CreateAppStoreVersionLocalization(ctx.Context, attrs, versionID)
		if err != nil {
			return err
		}
		if err := c.UpdatePreviewsAndScreenshotsIfNeeded(ctx, &locResp.Data, locConfig); err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UpdatePreviewsAndScreenshotsIfNeeded(ctx *context.Context, loc *asc.AppStoreVersionLocalization, config config.VersionLocalization) error {
	if loc.Relationships.AppPreviewSets != nil {
		var previewSets asc.AppPreviewSetsResponse
		_, err := c.client.FollowReference(ctx, loc.Relationships.AppPreviewSets.Links.Related, &previewSets)
		if err != nil {
			return err
		}
		err = c.UpdatePreviewSets(ctx, previewSets.Data, loc.ID, config.PreviewSets)
		if err != nil {
			return err
		}
	}

	if loc.Relationships.AppScreenshotSets != nil {
		var screenshotSets asc.AppScreenshotSetsResponse
		_, err := c.client.FollowReference(ctx, loc.Relationships.AppScreenshotSets.Links.Related, &screenshotSets)
		if err != nil {
			return err
		}
		err = c.UpdateScreenshotSets(ctx, screenshotSets.Data, loc.ID, config.ScreenshotSets)
		if err != nil {
			return err
		}
	}
	return nil
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

func (c *ascClient) UploadRoutingCoverage(ctx *context.Context, versionID string, config config.File) error {
	// TODO: I'm silencing an error here
	if covResp, _, _ := c.client.Apps.GetRoutingAppCoverageForAppStoreVersion(ctx, versionID, nil); covResp != nil {
		if _, err := c.client.Apps.DeleteRoutingAppCoverage(ctx, covResp.Data.ID); err != nil {
			return err
		}
	}
	create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
		resp, _, err := c.client.Apps.CreateRoutingAppCoverage(ctx, name, size, versionID)
		if err != nil {
			return "", nil, err
		}
		return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
	}
	commit := func(id string, checksum string) error {
		_, _, err := c.client.Apps.CommitRoutingAppCoverage(ctx, id, asc.Bool(true), &checksum)
		return err
	}
	return c.uploadFile(ctx, config.Path, create, commit)
}

func (c *ascClient) UpdatePreviewSets(ctx *context.Context, previewSets []asc.AppPreviewSet, appStoreVersionLocalizationID string, config config.PreviewSets) error {
	found := make(map[asc.PreviewType]bool)
	for i := range previewSets {
		previewSet := previewSets[i]
		previewType := *previewSet.Attributes.PreviewType
		found[previewType] = true
		previewsConfig := config.GetPreviews(previewType)
		if err := c.UploadPreviews(ctx, &previewSet, previewsConfig); err != nil {
			return err
		}
	}
	for previewType, previews := range config {
		t := previewType.APIValue()
		if found[t] {
			continue
		}
		previewSetResp, _, err := c.client.Apps.CreateAppPreviewSet(ctx, t, appStoreVersionLocalizationID)
		if err != nil {
			return err
		}
		if err = c.UploadPreviews(ctx, &previewSetResp.Data, previews); err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) UploadPreviews(ctx *context.Context, previewSet *asc.AppPreviewSet, previewConfigs []config.Preview) error {
	var g = parallel.New(ctx.MaxProcesses)
	for _, previewConfig := range previewConfigs {
		path := previewConfig.Path
		create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
			resp, _, err := c.client.Apps.CreateAppPreview(ctx, name, size, previewSet.ID)
			if err != nil {
				return "", nil, err
			}
			return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
		}
		commit := func(con config.Preview) commitFunc {
			return func(id string, checksum string) error {
				_, _, err := c.client.Apps.CommitAppPreview(ctx, id, asc.Bool(true), &checksum, &con.PreviewFrameTimeCode)
				return err
			}
		}(previewConfig)
		g.Go(func() error {
			return c.uploadFile(ctx, path, create, commit)
		})
	}
	return g.Wait()
}

func (c *ascClient) UpdateScreenshotSets(ctx *context.Context, screenshotSets []asc.AppScreenshotSet, appStoreVersionLocalizationID string, config config.ScreenshotSets) error {
	found := make(map[asc.ScreenshotDisplayType]bool)
	for i := range screenshotSets {
		screenshotSet := screenshotSets[i]
		screenshotType := *screenshotSet.Attributes.ScreenshotDisplayType
		found[screenshotType] = true
		screenshotConfig := config.GetScreenshots(screenshotType)
		if err := c.UploadScreenshots(ctx, &screenshotSet, screenshotConfig); err != nil {
			return err
		}
	}
	for screenshotType, screenshots := range config {
		t := screenshotType.APIValue()
		if found[t] {
			continue
		}
		screenshotSetResp, _, err := c.client.Apps.CreateAppScreenshotSet(ctx, t, appStoreVersionLocalizationID)
		if err != nil {
			return err
		}
		if err = c.UploadScreenshots(ctx, &screenshotSetResp.Data, screenshots); err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) UploadScreenshots(ctx *context.Context, screenshotSet *asc.AppScreenshotSet, config []config.File) error {
	var g = parallel.New(ctx.MaxProcesses)
	for _, screenshotConfig := range config {
		path := screenshotConfig.Path
		create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
			resp, _, err := c.client.Apps.CreateAppScreenshot(ctx, name, size, screenshotSet.ID)
			if err != nil {
				return "", nil, err
			}
			return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
		}
		commit := func(id string, checksum string) error {
			_, _, err := c.client.Apps.CommitAppScreenshot(ctx, id, asc.Bool(true), &checksum)
			return err
		}
		g.Go(func() error {
			return c.uploadFile(ctx, path, create, commit)
		})
	}
	return g.Wait()
}

func (c *ascClient) UpdateReviewDetails(ctx *context.Context, versionID string, config config.ReviewDetails) error {
	detailsResp, _, err := c.client.Submission.GetReviewDetailsForAppStoreVersion(ctx, versionID, nil)

	if err != nil {
		return c.CreateReviewDetail(ctx, versionID, config)
	}
	return c.UpdateReviewDetail(ctx, detailsResp.Data.ID, config)
}

func (c *ascClient) CreateReviewDetail(ctx *context.Context, versionID string, config config.ReviewDetails) error {
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
	_, _, err := c.client.Submission.CreateReviewDetail(ctx, &attributes, versionID)
	return err
}

func (c *ascClient) UpdateReviewDetail(ctx *context.Context, reviewDetailID string, config config.ReviewDetails) error {
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
	_, _, err := c.client.Submission.UpdateReviewDetail(ctx, reviewDetailID, &attributes)
	return err
}

func (c *ascClient) SubmitApp(ctx *context.Context, versionID string) error {
	_, _, err := c.client.Submission.CreateSubmission(ctx, versionID)
	return err
}

type createFunc func(name string, size int64) (id string, ops []asc.UploadOperation, err error)
type commitFunc func(id string, checksum string) error

func (c *ascClient) uploadFile(ctx *context.Context, path string, create createFunc, commit commitFunc) error {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer closer.Close(f)

	fstat, err := os.Stat(path)
	if err != nil {
		return err
	}
	checksum, err := md5Checksum(f)
	if err != nil {
		return err
	}

	id, ops, err := create(fstat.Name(), fstat.Size())
	if err != nil {
		return err
	}
	err = c.client.Upload(ctx, ops, f)
	if err != nil {
		return err
	}
	return commit(id, checksum)
}
