package client

import (
	"os"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

func (c *ascClient) UpdateAppLocalizations(ctx *context.Context, app *asc.App, config config.AppLocalizations) error {
	appInfosResp, _, err := c.client.Apps.ListAppInfosForApp(ctx, app.ID, nil)
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
			found[locale] = true
			locConfig := config[locale]
			_, _, err := c.client.Apps.UpdateAppInfoLocalization(ctx, loc.ID, &asc.AppInfoLocalizationUpdateRequestAttributes{
				Name:              &locConfig.Name,
				PrivacyPolicyText: &locConfig.PrivacyPolicyText,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
				Subtitle:          &locConfig.Subtitle,
			})
			if err != nil {
				return err
			}
		}

		for locale, locConfig := range config {
			if found[locale] {
				continue
			}
			_, _, err := c.client.Apps.CreateAppInfoLocalization(ctx.Context, asc.AppInfoLocalizationCreateRequestAttributes{
				Locale:            locale,
				Name:              &locConfig.Name,
				PrivacyPolicyText: &locConfig.PrivacyPolicyText,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
				Subtitle:          &locConfig.Subtitle,
			}, appInfo.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ascClient) CreateVersionIfNeeded(ctx *context.Context, app *asc.App, build *asc.Build, config config.Version) (*asc.AppStoreVersion, error) {
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
	versionsResp, _, err := c.client.Apps.ListAppStoreVersionsForApp(ctx, app.ID, &asc.ListAppStoreVersionsQuery{
		FilterVersionString: []string{ctx.Version},
		FilterPlatform:      []string{string(platform)},
	})
	if len(versionsResp.Data) == 0 {
		versionResp, _, err = c.client.Apps.CreateAppStoreVersion(ctx, asc.AppStoreVersionCreateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: earliestReleaseDate,
			Platform:            platform,
			ReleaseType:         releaseTypeP,
			UsesIDFA:            asc.Bool(config.IDFADeclaration != nil),
			VersionString:       ctx.Version,
		}, app.ID, &build.ID)
	} else {
		latestVersion := versionsResp.Data[0]
		versionResp, _, err = c.client.Apps.UpdateAppStoreVersion(ctx, latestVersion.ID, &asc.AppStoreVersionUpdateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: earliestReleaseDate,
			ReleaseType:         releaseTypeP,
			UsesIDFA:            asc.Bool(config.IDFADeclaration != nil),
			VersionString:       &ctx.Version,
		}, &build.ID)
	}
	return &versionResp.Data, err
}

func (c *ascClient) UpdateVersionLocalizations(ctx *context.Context, version *asc.AppStoreVersion, config config.VersionLocalizations) error {
	locListResp, _, err := c.client.Apps.ListLocalizationsForAppStoreVersion(ctx, version.ID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for _, loc := range locListResp.Data {
		locale := *loc.Attributes.Locale
		found[locale] = true
		locConfig := config[locale]
		updatedLocResp, _, err := c.client.Apps.UpdateAppStoreVersionLocalization(ctx, loc.ID, &asc.AppStoreVersionLocalizationUpdateRequestAttributes{
			Description:     &locConfig.Description,
			Keywords:        &locConfig.Keywords,
			MarketingURL:    &locConfig.MarketingURL,
			PromotionalText: &locConfig.PromotionalText,
			SupportURL:      &locConfig.SupportURL,
			// TODO: Modifying this field results in a 409 Conflict error
			// WhatsNew: &locConfig.WhatsNewText,
		})
		if err != nil {
			return err
		}
		loc = updatedLocResp.Data
		if loc.Relationships.AppPreviewSets != nil {
			var previewSets asc.AppPreviewSetsResponse
			_, err = c.client.FollowReference(ctx, loc.Relationships.AppPreviewSets.Links.Related, &previewSets)
			if err != nil {
				return err
			}
			err = c.UpdatePreviewSets(ctx, previewSets.Data, loc.ID, locConfig.PreviewSets)
			if err != nil {
				return err
			}
		}
		if loc.Relationships.AppScreenshotSets != nil {
			var screenshotSets asc.AppScreenshotSetsResponse
			_, err = c.client.FollowReference(ctx, loc.Relationships.AppScreenshotSets.Links.Related, &screenshotSets)
			if err != nil {
				return err
			}
			err = c.UpdateScreenshotSets(ctx, screenshotSets.Data, loc.ID, locConfig.ScreenshotSets)
			if err != nil {
				return err
			}
		}
	}

	for locale, locConfig := range config {
		if found[locale] {
			continue
		}
		locResp, _, err := c.client.Apps.CreateAppStoreVersionLocalization(ctx.Context, asc.AppStoreVersionLocalizationCreateRequestAttributes{
			Description:     &locConfig.Description,
			Keywords:        &locConfig.Keywords,
			Locale:          locale,
			MarketingURL:    &locConfig.MarketingURL,
			PromotionalText: &locConfig.PromotionalText,
			SupportURL:      &locConfig.SupportURL,
			WhatsNew:        &locConfig.WhatsNewText,
		}, version.ID)
		if err != nil {
			return err
		}
		loc := locResp.Data
		if loc.Relationships.AppPreviewSets != nil {
			var previewSets asc.AppPreviewSetsResponse
			_, err = c.client.FollowReference(ctx, loc.Relationships.AppPreviewSets.Links.Related, &previewSets)
			if err != nil {
				return err
			}
			err = c.UpdatePreviewSets(ctx, previewSets.Data, loc.ID, locConfig.PreviewSets)
			if err != nil {
				return err
			}
		}
		if loc.Relationships.AppScreenshotSets != nil {
			var screenshotSets asc.AppScreenshotSetsResponse
			_, err = c.client.FollowReference(ctx, loc.Relationships.AppScreenshotSets.Links.Related, &screenshotSets)
			if err != nil {
				return err
			}
			err = c.UpdateScreenshotSets(ctx, screenshotSets.Data, loc.ID, locConfig.ScreenshotSets)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ascClient) UpdateIDFADeclaration(ctx *context.Context, version *asc.AppStoreVersion, config config.IDFADeclaration) error {
	existingDeclResp, _, err := c.client.Submission.GetIDFADeclarationForAppStoreVersion(ctx, version.ID, nil)
	if err != nil || existingDeclResp.Data.ID == "" {
		_, _, err = c.client.Submission.CreateIDFADeclaration(ctx, asc.IDFADeclarationCreateRequestAttributes{
			AttributesActionWithPreviousAd:        config.AttributesActionWithPreviousAd,
			AttributesAppInstallationToPreviousAd: config.AttributesAppInstallationToPreviousAd,
			HonorsLimitedAdTracking:               config.HonorsLimitedAdTracking,
			ServesAds:                             config.ServesAds,
		}, version.ID)
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

func (c *ascClient) UploadRoutingCoverage(ctx *context.Context, version *asc.AppStoreVersion, config config.File) error {
	// TODO: I'm silencing an error here
	if covResp, _, _ := c.client.Apps.GetRoutingAppCoverageForAppStoreVersion(ctx, version.ID, nil); covResp != nil {
		if _, err := c.client.Apps.DeleteRoutingAppCoverage(ctx, covResp.Data.ID); err != nil {
			return err
		}
	}
	create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
		resp, _, err := c.client.Apps.CreateRoutingAppCoverage(ctx, name, size, version.ID)
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
	for _, previewSet := range previewSets {
		previewType := *previewSet.Attributes.PreviewType
		found[previewType] = true
		previewsConfig := config.GetPreviews(previewType)
		c.UploadPreviews(ctx, &previewSet, previewsConfig)
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
		err = c.UploadPreviews(ctx, &previewSetResp.Data, previews)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) UploadPreviews(ctx *context.Context, previewSet *asc.AppPreviewSet, config []config.Preview) error {
	for _, previewConfig := range config {
		create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
			resp, _, err := c.client.Apps.CreateAppPreview(ctx, name, size, previewSet.ID)
			if err != nil {
				return "", nil, err
			}
			return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
		}
		commit := func(id string, checksum string) error {
			_, _, err := c.client.Apps.CommitAppPreview(ctx, id, asc.Bool(true), &checksum, &previewConfig.PreviewFrameTimeCode)
			return err
		}
		if err := c.uploadFile(ctx, previewConfig.Path, create, commit); err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) UpdateScreenshotSets(ctx *context.Context, screenshotSets []asc.AppScreenshotSet, appStoreVersionLocalizationID string, config config.ScreenshotSets) error {
	found := make(map[asc.ScreenshotDisplayType]bool)
	for _, screenshotSet := range screenshotSets {
		screenshotType := *screenshotSet.Attributes.ScreenshotDisplayType
		found[screenshotType] = true
		screenshotConfig := config.GetScreenshots(screenshotType)
		c.UploadScreenshots(ctx, &screenshotSet, screenshotConfig)
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
	for _, screenshotConfig := range config {
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
		if err := c.uploadFile(ctx, screenshotConfig.Path, create, commit); err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) UpdateReviewDetails(ctx *context.Context, version *asc.AppStoreVersion, config config.ReviewDetails) error {
	detailsResp, _, err := c.client.Submission.GetReviewDetailsForAppStoreVersion(ctx, version.ID, nil)
	if err != nil {
		_, _, err = c.client.Submission.CreateReviewDetail(ctx, &asc.AppStoreReviewDetailCreateRequestAttributes{
			ContactEmail:        &config.Contact.Email,
			ContactFirstName:    &config.Contact.FirstName,
			ContactLastName:     &config.Contact.LastName,
			ContactPhone:        &config.Contact.Phone,
			DemoAccountName:     &config.DemoAccount.Name,
			DemoAccountPassword: &config.DemoAccount.Password,
			DemoAccountRequired: &config.DemoAccount.Required,
			Notes:               &config.Notes,
		}, version.ID)
	} else {
		_, _, err = c.client.Submission.UpdateReviewDetail(ctx, detailsResp.Data.ID, &asc.AppStoreReviewDetailUpdateRequestAttributes{
			ContactEmail:        &config.Contact.Email,
			ContactFirstName:    &config.Contact.FirstName,
			ContactLastName:     &config.Contact.LastName,
			ContactPhone:        &config.Contact.Phone,
			DemoAccountName:     &config.DemoAccount.Name,
			DemoAccountPassword: &config.DemoAccount.Password,
			DemoAccountRequired: &config.DemoAccount.Required,
			Notes:               &config.Notes,
		})
	}
	return err
}

func (c *ascClient) SubmitApp(ctx *context.Context, version *asc.AppStoreVersion) error {
	_, _, err := c.client.Submission.CreateSubmission(ctx, version.ID)
	return err
}

type createFunc func(name string, size int64) (id string, ops []asc.UploadOperation, err error)
type commitFunc func(id string, checksum string) error

func (c *ascClient) uploadFile(ctx *context.Context, path string, create createFunc, commit commitFunc) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

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
