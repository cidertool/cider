package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Test UpdateApp

func TestUpdateApp_Happy(t *testing.T) {
	ctx, client := newTestContext(
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateApp(ctx.Context, testID, testID, config.App{
		UsesThirdPartyContent: asc.Bool(true),
		Availability: &config.Availability{
			AvailableInNewTerritories: asc.Bool(true),
			Pricing: []config.PriceSchedule{
				{Tier: "0"},
			},
			Territories: []string{"USA", "JPN"},
		},
		PrimaryLocale: "en-US",
		BundleID:      "com.app.bundleid",
	})
	assert.NoError(t, err)
}

// Test UpdateAppLocalizations

func TestUpdateAppLocalizations_Happy(t *testing.T) {
	prepareForSubmission := asc.AppStoreVersionStatePrepareForSubmission
	readyForSale := asc.AppStoreVersionStateReadyForSale
	ctx, client := newTestContext(
		response{
			Response: asc.AppInfosResponse{
				Data: []asc.AppInfo{
					{
						Attributes: &asc.AppInfoAttributes{
							AppStoreState: &readyForSale,
						},
						ID: "TEST",
					},
					{
						Attributes: &asc.AppInfoAttributes{
							AppStoreState: &prepareForSubmission,
						},
						ID: "TEST",
					},
				},
			},
		},
		response{
			Response: asc.AppInfoLocalizationsResponse{
				Data: []asc.AppInfoLocalization{
					{
						Attributes: &asc.AppInfoLocalizationAttributes{
							Locale: asc.String("en-US"),
							Name:   asc.String("My App"),
						},
					},
					{
						Attributes: &asc.AppInfoLocalizationAttributes{
							Locale: asc.String("en-GB"),
							Name:   asc.String("No App"),
						},
					},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateAppLocalizations(ctx.Context, testID, config.AppLocalizations{
		"en-US": {
			Name: "My App",
		},
		"ja": {
			Name: "Your App",
		},
	})
	assert.NoError(t, err)
}

// Test CreateVersionIfNeeded

func TestCreateVersionIfNeeded_Happy(t *testing.T) {
	now := time.Now()
	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":[{}]}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	version, err := client.CreateVersionIfNeeded(ctx.Context, testID, testID, config.Version{
		Platform:            config.PlatformiOS,
		ReleaseType:         config.ReleaseTypeAfterApproval,
		EarliestReleaseDate: &now,
	})
	assert.NoError(t, err)
	assert.NotNil(t, version)
}

// Test UpdateVersionLocalizations

func TestUpdateVersionLocalizations_Happy(t *testing.T) {

}

// Test UpdateIDFADeclaration

func TestUpdateIDFADeclaration_Happy(t *testing.T) {
	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateIDFADeclaration(ctx.Context, testID, config.IDFADeclaration{})
	assert.NoError(t, err)
}

func TestUpdateIDFADeclaration_ErrUpdate(t *testing.T) {
	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateIDFADeclaration(ctx.Context, testID, config.IDFADeclaration{})
	assert.Error(t, err)
}

func TestUpdateIDFADeclaration_ErrCreate(t *testing.T) {
	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateIDFADeclaration(ctx.Context, testID, config.IDFADeclaration{})
	assert.Error(t, err)
}

// Test UploadRoutingCoverage

func TestUploadRoutingCoverage_Happy(t *testing.T) {

}

// Test UpdateReviewDetails

func TestUpdateReviewDetails_Happy(t *testing.T) {
	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateReviewDetails(ctx.Context, testID, config.ReviewDetails{
		Contact:     &config.ContactPerson{},
		DemoAccount: &config.DemoAccount{},
	})
	assert.NoError(t, err)
}

func TestUpdateReviewDetails_ErrUpdate(t *testing.T) {
	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateReviewDetails(ctx.Context, testID, config.ReviewDetails{})
	assert.Error(t, err)
}

func TestUpdateReviewDetails_ErrCreate(t *testing.T) {
	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateReviewDetails(ctx.Context, testID, config.ReviewDetails{
		Contact:     &config.ContactPerson{},
		DemoAccount: &config.DemoAccount{},
	})
	assert.Error(t, err)
}

// Test SubmitApp

func TestSubmitApp_Happy(t *testing.T) {
	ctx, client := newTestContext(
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.SubmitApp(ctx.Context, testID)
	assert.NoError(t, err)
}

func TestSubmitApp_Err(t *testing.T) {
	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.SubmitApp(ctx.Context, testID)
	assert.Error(t, err)
}
