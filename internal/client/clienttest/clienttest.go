// Package clienttest provides utilities for mocking the App Store Connect client.
package clienttest

import (
	"net/http"
	"time"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

// Client is a mock type that implements client.Client.
type Client struct{}

// Credentials mocks the credentials store used by the client.
type Credentials struct {
	client *http.Client
}

// GetAppForBundleID mocks returning an app for a given bundle ID.
func (c *Client) GetAppForBundleID(ctx *context.Context, bundleID string) (*asc.App, error) {
	return &asc.App{
		Attributes: &asc.AppAttributes{
			AvailableInNewTerritories: asc.Bool(false),
			BundleID:                  &bundleID,
			ContentRightsDeclaration:  asc.String("DOES_NOT_USE_THIRD_PARTY_CONTENT"),
			IsOrEverWasMadeForKids:    asc.Bool(false),
			Name:                      asc.String("TEST"),
			PrimaryLocale:             asc.String("en-US"),
			Sku:                       asc.String("TEST"),
		},
		ID: "TEST",
	}, nil
}

// GetAppInfo mocks returning an app info corresponding to an app.
func (c *Client) GetAppInfo(ctx *context.Context, appID string) (*asc.AppInfo, error) {
	appStoreAgeRating := asc.AppStoreAgeRatingFourPlus
	appStoreState := asc.AppStoreVersionStatePrepareForSubmission
	brazilAgeRating := asc.BrazilAgeRatingL
	kidsAgeBand := asc.KidsAgeBandFiveAndUnder

	return &asc.AppInfo{
		Attributes: &asc.AppInfoAttributes{
			AppStoreAgeRating: &appStoreAgeRating,
			AppStoreState:     &appStoreState,
			BrazilAgeRating:   &brazilAgeRating,
			KidsAgeBand:       &kidsAgeBand,
		},
		ID: "TEST",
	}, nil
}

// GetBuild mocks returning the latest valid build corresponding to an app.
func (c *Client) GetBuild(ctx *context.Context, app *asc.App) (*asc.Build, error) {
	var testImageSize = 140

	return &asc.Build{
		Attributes: &asc.BuildAttributes{
			ExpirationDate: &asc.DateTime{Time: time.Now()},
			Expired:        asc.Bool(false),
			IconAssetToken: &asc.ImageAsset{
				Height:      &testImageSize,
				Width:       &testImageSize,
				TemplateURL: asc.String(""),
			},
			MinOsVersion:            asc.String("9.0"),
			ProcessingState:         asc.String("VALID"),
			UploadedDate:            &asc.DateTime{Time: time.Now()},
			UsesNonExemptEncryption: asc.Bool(false),
			Version:                 asc.String("99"),
		},
		ID: "TEST",
	}, nil
}

// ReleaseForAppIsInitial mocks returning whether or not an app has released on the App Store before.
func (c *Client) ReleaseForAppIsInitial(ctx *context.Context, appID string) (bool, error) {
	return false, nil
}

// UpdateBetaAppLocalizations mocks updating localized properties for a beta app.
func (c *Client) UpdateBetaAppLocalizations(ctx *context.Context, appID string, config config.TestflightLocalizations) error {
	return nil
}

// UpdateBetaBuildDetails mocks updating beta build details for a beta build.
func (c *Client) UpdateBetaBuildDetails(ctx *context.Context, buildID string, config config.Testflight) error {
	return nil
}

// UpdateBetaBuildLocalizations mocks updating localized properties for a beta build.
func (c *Client) UpdateBetaBuildLocalizations(ctx *context.Context, buildID string, config config.TestflightLocalizations) error {
	return nil
}

// UpdateBetaLicenseAgreement mocks updating the beta license agreement for an app.
func (c *Client) UpdateBetaLicenseAgreement(ctx *context.Context, appID string, config config.Testflight) error {
	return nil
}

// AssignBetaGroups mocks assigning beta groups to a beta build.
func (c *Client) AssignBetaGroups(ctx *context.Context, appID string, buildID string, groups []config.BetaGroup) error {
	return nil
}

// AssignBetaTesters mocks assigning beta testers to a beta build.
func (c *Client) AssignBetaTesters(ctx *context.Context, appID string, buildID string, testers []config.BetaTester) error {
	return nil
}

// UpdateBetaReviewDetails mocks updating review details for a beta app.
func (c *Client) UpdateBetaReviewDetails(ctx *context.Context, appID string, config config.ReviewDetails) error {
	return nil
}

// SubmitBetaApp mocks submitting an app to Testflight.
func (c *Client) SubmitBetaApp(ctx *context.Context, buildID string) error {
	return nil
}

// UpdateApp mocks updating properties for an app.
func (c *Client) UpdateApp(ctx *context.Context, appID string, appInfoID string, versionID string, config config.App) error {
	return nil
}

// UpdateAppLocalizations mocks updating localized properties for an app info.
func (c *Client) UpdateAppLocalizations(ctx *context.Context, appID string, config config.AppLocalizations) error {
	return nil
}

// CreateVersionIfNeeded mocks creating a version, if one does not exist for the tag, and returning its model.
func (c *Client) CreateVersionIfNeeded(ctx *context.Context, appID string, buildID string, config config.Version) (*asc.AppStoreVersion, error) {
	appStoreState := asc.AppStoreVersionStatePrepareForSubmission
	platform := asc.PlatformIOS

	return &asc.AppStoreVersion{
		Attributes: &asc.AppStoreVersionAttributes{
			AppStoreState:       &appStoreState,
			Copyright:           asc.String("2020"),
			CreatedDate:         &asc.DateTime{Time: time.Now()},
			Downloadable:        asc.Bool(true),
			EarliestReleaseDate: &asc.DateTime{Time: time.Now()},
			Platform:            &platform,
			ReleaseType:         asc.String("AFTER_APPROVAL"),
			UsesIDFA:            asc.Bool(true),
			VersionString:       asc.String("1.0"),
		},
		ID: "TEST",
	}, nil
}

// UpdateVersionLocalizations mocks updating localized properties for a version.
func (c *Client) UpdateVersionLocalizations(ctx *context.Context, versionID string, config config.VersionLocalizations) error {
	return nil
}

// UpdateIDFADeclaration mocks updating the IDFA declaration for a version.
func (c *Client) UpdateIDFADeclaration(ctx *context.Context, versionID string, config config.IDFADeclaration) error {
	return nil
}

// UploadRoutingCoverage mocks uploading a routing coverage file for a version.
func (c *Client) UploadRoutingCoverage(ctx *context.Context, versionID string, config config.File) error {
	return nil
}

// UpdateReviewDetails mocks updating review details for a version.
func (c *Client) UpdateReviewDetails(ctx *context.Context, versionID string, config config.ReviewDetails) error {
	return nil
}

// EnablePhasedRelease mocks enabling phased release for a version.
func (c *Client) EnablePhasedRelease(ctx *context.Context, versionID string) error {
	return nil
}

// SubmitApp mocks submitting a version to the App Store.
func (c *Client) SubmitApp(ctx *context.Context, versionID string) error {
	return nil
}

// Project mocks returning a project built from API values.
func (c *Client) Project() (*config.Project, error) {
	return &config.Project{}, nil
}

// Client returns an http.Client for the mock Credentials instance.
func (c *Credentials) Client() *http.Client {
	if c.client == nil {
		c.client = &http.Client{}
	}

	return c.client
}
