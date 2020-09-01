package client

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

// Client is an abstraction of an App Store Connect API client's functionality
type Client interface {
	GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error)
	GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error)
	UpdateBetaAppLocalizations(ctx *context.Context, app *asc.App, config config.TestflightLocalizations) error
	UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config *config.TestflightForApp) error
	UpdateBetaBuildLocalizations(ctx *context.Context, build *asc.Build, config config.TestflightLocalizations) error
	UpdateBetaLicenseAgreement(ctx *context.Context, app *asc.App, config *config.TestflightForApp) error
	AssignBetaGroups(ctx *context.Context, build *asc.Build, groups []string) error
	AssignBetaTesters(ctx *context.Context, build *asc.Build, testers []config.BetaTester) error
	UpdateBetaReviewDetails(ctx *context.Context, app *asc.App, config *config.ReviewDetails) error
	SubmitBetaApp(ctx *context.Context, build *asc.Build) error
	UpdateAppLocalizations(ctx *context.Context, app *asc.App, config config.AppLocalizations) error
	CreateVersionIfNeeded(ctx *context.Context, app *asc.App, build *asc.Build, config *config.Version) (*asc.AppStoreVersion, error)
	UpdateVersionLocalizations(ctx *context.Context, version *asc.AppStoreVersion, config config.VersionLocalizations) error
	UpdateIDFADeclaration(ctx *context.Context, version *asc.AppStoreVersion, config *config.IDFADeclaration) error
	UploadRoutingCoverage(ctx *context.Context, version *asc.AppStoreVersion, config *config.File) error
	UpdatePreviewSets(ctx *context.Context, previewSets []asc.AppPreviewSet, appStoreVersionLocalizationID string, config config.PreviewSets) error
	UpdateScreenshotSets(ctx *context.Context, screenshotSets []asc.AppScreenshotSet, appStoreVersionLocalizationID string, config config.ScreenshotSets) error
	UpdateReviewDetails(ctx *context.Context, version *asc.AppStoreVersion, config *config.ReviewDetails) error
	SubmitApp(ctx *context.Context, version *asc.AppStoreVersion) error
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
		FilterApp:                      []string{app.ID},
		FilterPreReleaseVersionVersion: []string{ctx.Version},
	})
	if err != nil {
		return nil, fmt.Errorf("build not found matching app %s and version %s: %w", *app.Attributes.BundleID, ctx.Version, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("build not found matching app %s and version %s", *app.Attributes.BundleID, ctx.Version)
	}
	return &resp.Data[0], nil
}

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

func (c *ascClient) UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config *config.TestflightForApp) error {
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

func (c *ascClient) UpdateBetaLicenseAgreement(ctx *context.Context, app *asc.App, config *config.TestflightForApp) error {
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

func (c *ascClient) UpdateBetaReviewDetails(ctx *context.Context, app *asc.App, config *config.ReviewDetails) error {
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

func (c *ascClient) CreateVersionIfNeeded(ctx *context.Context, app *asc.App, build *asc.Build, config *config.Version) (*asc.AppStoreVersion, error) {
	platform, err := config.Platform.APIValue()
	if err != nil {
		return nil, err
	}
	releaseType, err := config.ReleaseType.APIValue()
	if err != nil {
		return nil, err
	}
	var versionResp *asc.AppStoreVersionResponse
	versionsResp, _, err := c.client.Apps.ListAppStoreVersionsForApp(ctx, app.ID, &asc.ListAppStoreVersionsQuery{FilterVersionString: []string{ctx.Version}, FilterPlatform: []string{string(platform)}})
	if len(versionsResp.Data) == 0 {
		versionResp, _, err = c.client.Apps.CreateAppStoreVersion(ctx, asc.AppStoreVersionCreateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: config.EarliestReleaseDate,
			Platform:            platform,
			ReleaseType:         &releaseType,
			UsesIDFA:            asc.Bool(config.IDFADeclaration != nil),
			VersionString:       ctx.Version,
		}, app.ID, &build.ID)
	} else {
		latestVersion := versionsResp.Data[0]
		versionResp, _, err = c.client.Apps.UpdateAppStoreVersion(ctx, latestVersion.ID, &asc.AppStoreVersionUpdateRequestAttributes{
			Copyright:           &config.Copyright,
			EarliestReleaseDate: config.EarliestReleaseDate,
			ReleaseType:         &releaseType,
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
			WhatsNew:        &locConfig.WhatsNewText,
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

func (c *ascClient) UpdateIDFADeclaration(ctx *context.Context, version *asc.AppStoreVersion, config *config.IDFADeclaration) error {
	existingDeclResp, _, err := c.client.Submission.GetIDFADeclarationForAppStoreVersion(ctx, version.ID, nil)
	if err != nil {
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

func (c *ascClient) UploadRoutingCoverage(ctx *context.Context, version *asc.AppStoreVersion, config *config.File) error {
	f, err := os.Open(config.Path)
	if err != nil {
		return err
	}
	fstat, err := os.Stat(config.Path)
	if err != nil {
		return err
	}
	resp, _, err := c.client.Apps.CreateRoutingAppCoverage(ctx, fstat.Name(), fstat.Size(), version.ID)
	if err != nil {
		return err
	}
	ops := resp.Data.Attributes.UploadOperations
	err = c.client.Upload(ctx, ops, f)
	if err != nil {
		return err
	}
	checksum, err := md5Checksum(config.Path)
	if err != nil {
		return err
	}
	_, _, err = c.client.Apps.CommitRoutingAppCoverage(ctx, resp.Data.ID, asc.Bool(true), &checksum)
	return err
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
		f, err := os.Open(previewConfig.Path)
		if err != nil {
			return err
		}
		fstat, err := os.Stat(previewConfig.Path)
		if err != nil {
			return err
		}
		resp, _, err := c.client.Apps.CreateAppPreview(ctx, fstat.Name(), fstat.Size(), previewSet.ID)
		if err != nil {
			return err
		}
		ops := resp.Data.Attributes.UploadOperations
		err = c.client.Upload(ctx, ops, f)
		if err != nil {
			return err
		}
		checksum, err := md5Checksum(previewConfig.Path)
		if err != nil {
			return err
		}
		_, _, err = c.client.Apps.CommitAppPreview(ctx, resp.Data.ID, asc.Bool(true), &checksum, &previewConfig.PreviewFrameTimeCode)
		if err != nil {
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
		err = c.UploadScreenshots(ctx, &screenshotSetResp.Data, screenshots)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) UploadScreenshots(ctx *context.Context, screenshotSet *asc.AppScreenshotSet, config []config.File) error {
	for _, screenshotConfig := range config {
		f, err := os.Open(screenshotConfig.Path)
		if err != nil {
			return err
		}
		fstat, err := os.Stat(screenshotConfig.Path)
		if err != nil {
			return err
		}
		resp, _, err := c.client.Apps.CreateAppScreenshot(ctx, fstat.Name(), fstat.Size(), screenshotSet.ID)
		if err != nil {
			return err
		}
		ops := resp.Data.Attributes.UploadOperations
		err = c.client.Upload(ctx, ops, f)
		if err != nil {
			return err
		}
		checksum, err := md5Checksum(screenshotConfig.Path)
		if err != nil {
			return err
		}
		_, _, err = c.client.Apps.CommitAppScreenshot(ctx, resp.Data.ID, asc.Bool(true), &checksum)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ascClient) UpdateReviewDetails(ctx *context.Context, version *asc.AppStoreVersion, config *config.ReviewDetails) error {
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

func md5Checksum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
