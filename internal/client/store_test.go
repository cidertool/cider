/**
Copyright (C) 2020 Aaron Sky.

This file is part of Cider, a tool for automating submission
of apps to Apple's App Stores.

Cider is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cider is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cider.  If not, see <http://www.gnu.org/licenses/>.
*/

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
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateApp(ctx.Context, testID, testID, testID, config.App{
		AgeRatingDeclaration:  &config.AgeRatingDeclaration{},
		Categories:            &config.Categories{},
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
