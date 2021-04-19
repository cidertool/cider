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
	t.Parallel()

	now := time.Now()

	ctx, client := newTestContext(
		response{
			Response: asc.TerritoriesResponse{
				Data: []asc.Territory{
					{ID: "USA"},
					{ID: "GBP"},
					{ID: "JPN"},
				},
			},
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
		AgeRatingDeclaration: &config.AgeRatingDeclaration{},
		Categories: &config.Categories{
			Primary:                "TEST",
			PrimarySubcategories:   [2]string{"TEST", "TEST"},
			Secondary:              "TEST",
			SecondarySubcategories: [2]string{"TEST", "TEST"},
		},
		UsesThirdPartyContent: asc.Bool(true),
		Availability: &config.Availability{
			AvailableInNewTerritories: asc.Bool(true),
			Pricing: []config.PriceSchedule{
				{Tier: "0", StartDate: &now},
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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	asset := newTestAsset(t, "TEST")
	previewID := "TEST-Preview-en_US"
	screenshotID := "TEST-Screenshot-en_US"

	localizations := config.VersionLocalizations{
		"en-US": {
			Description: "TEST",
			Keywords:    "TEST,TEST",
			PreviewSets: config.PreviewSets{
				config.PreviewTypeWatchSeries3: []config.Preview{
					{
						File: config.File{
							Path: asset.Name,
						},
					},
				},
			},
			ScreenshotSets: config.ScreenshotSets{
				config.ScreenshotTypeWatchSeries3: []config.File{
					{Path: asset.Name},
				},
			},
			MarketingURL:    "TEST",
			PromotionalText: "TEST",
			SupportURL:      "TEST",
			WhatsNewText:    "Going away",
		},
		"ja": {
			Description:     "TEST",
			Keywords:        "TEST,TEST",
			PreviewSets:     config.PreviewSets{},
			ScreenshotSets:  config.ScreenshotSets{},
			MarketingURL:    "TEST",
			PromotionalText: "TEST",
			SupportURL:      "TEST",
			WhatsNewText:    "Going away",
		},
	}

	ctx, client := newTestContext()

	previewSetsURL, err := ctx.URL(previewID)
	assert.NoError(t, err)

	screenshotSetsURL, err := ctx.URL(screenshotID)
	assert.NoError(t, err)

	previewURL, err := asset.URL(ctx, previewID)
	assert.NoError(t, err)

	screenshotURL, err := asset.URL(ctx, screenshotID)
	assert.NoError(t, err)

	ctx.SetResponses(
		// List app store version localizations
		response{
			Response: asc.AppStoreVersionLocalizationsResponse{
				Data: []asc.AppStoreVersionLocalization{
					{
						Attributes: &asc.AppStoreVersionLocalizationAttributes{
							Locale: asc.String("en-US"),
						},
						ID: "TEST",
						Relationships: &asc.AppStoreVersionLocalizationRelationships{
							AppPreviewSets: &asc.PagedRelationship{
								Links: &asc.RelationshipLinks{
									Related: &asc.Reference{
										URL: *previewSetsURL,
									},
								},
							},
							AppScreenshotSets: &asc.PagedRelationship{
								Links: &asc.RelationshipLinks{
									Related: &asc.Reference{
										URL: *screenshotSetsURL,
									},
								},
							},
						},
					},
					{
						Attributes: &asc.AppStoreVersionLocalizationAttributes{
							Locale: asc.String("en-GB"),
						},
						ID: "TEST",
					},
				},
			},
		},
		// Update app store version localization - #1
		response{
			Response: asc.AppStoreVersionLocalizationResponse{
				Data: asc.AppStoreVersionLocalization{
					Attributes: &asc.AppStoreVersionLocalizationAttributes{
						Locale: asc.String("en-US"),
					},
					ID: "TEST",
					Relationships: &asc.AppStoreVersionLocalizationRelationships{
						AppPreviewSets: &asc.PagedRelationship{
							Links: &asc.RelationshipLinks{
								Related: &asc.Reference{
									URL: *previewSetsURL,
								},
							},
						},
						AppScreenshotSets: &asc.PagedRelationship{
							Links: &asc.RelationshipLinks{
								Related: &asc.Reference{
									URL: *screenshotSetsURL,
								},
							},
						},
					},
				},
			},
		},
		// Follow app preview sets relationship reference
		response{
			Response: asc.AppPreviewSetsResponse{},
		},
		// Create missing app preview set
		response{
			Response: asc.AppPreviewSetResponse{
				Data: asc.AppPreviewSet{
					ID: "TEST",
				},
			},
		},
		// List app previews for preview set
		response{
			Response: asc.AppPreviewsResponse{
				Data: []asc.AppPreview{},
			},
		},
		// Create app preview
		response{
			Response: asc.AppPreviewResponse{
				Data: asc.AppPreview{
					Attributes: &asc.AppPreviewAttributes{
						AssetDeliveryState: nil,
						FileName:           &asset.Name,
						FileSize:           &asset.Size,
						MimeType:           asc.String("text/plain"),
						UploadOperations: []asc.UploadOperation{
							{
								Length: asc.Int(int(asset.Size)),
								Method: asc.String("PATCH"),
								Offset: asc.Int(0),
								URL:    asc.String(previewURL.String()),
							},
						},
					},
					ID: screenshotID,
				},
			},
		},
		// Upload operation #1 response
		response{
			RawResponse: `{}`,
		},
		// Commit app preview
		response{
			Response: asc.AppPreviewResponse{},
		},
		// Follow app screenshot sets relationship reference
		response{
			Response: asc.AppScreenshotSetsResponse{},
		},
		// Create missing app screenshot set
		response{
			Response: asc.AppScreenshotSetResponse{},
		},
		// Get app screenshots for app screenshot set
		response{
			Response: asc.AppScreenshotsResponse{},
		},
		// Create app screenshot
		response{
			Response: asc.AppScreenshotResponse{
				Data: asc.AppScreenshot{
					Attributes: &asc.AppScreenshotAttributes{
						AssetDeliveryState: nil,
						FileName:           &asset.Name,
						FileSize:           &asset.Size,
						UploadOperations: []asc.UploadOperation{
							{
								Length: asc.Int(int(asset.Size)),
								Method: asc.String("PATCH"),
								Offset: asc.Int(0),
								URL:    asc.String(screenshotURL.String()),
							},
						},
					},
					ID: screenshotID,
				},
			},
		},
		// Upload operation #1 response
		response{
			RawResponse: `{}`,
		},
		// Commit app screenshot
		response{
			Response: asc.AppScreenshotResponse{},
		},
		// Update app store version localization - #2
		response{
			Response: asc.AppStoreVersionLocalizationResponse{
				Data: asc.AppStoreVersionLocalization{
					Attributes: &asc.AppStoreVersionLocalizationAttributes{
						Locale: asc.String("en-GB"),
					},
					ID: "TEST",
				},
			},
		},
	)

	defer ctx.Close()

	ctx.Context.MaxProcesses = 1

	err = client.UpdateVersionLocalizations(ctx.Context, testID, localizations)

	assert.NoError(t, err)
}

// Test UpdateIDFADeclaration

func TestUpdateIDFADeclaration_Happy(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			Response: asc.IDFADeclarationResponse{
				Data: asc.IDFADeclaration{
					ID: testID,
				},
			},
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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	asset := newTestAsset(t, "TEST")
	attachmentID := "TEST-Attachment-1"

	ctx, client := newTestContext()

	defer ctx.Close()

	attachmentURL, err := asset.URL(ctx, attachmentID)
	assert.NoError(t, err)

	ctx.SetResponses(
		// GetReviewDetailsForAppStoreVersion
		response{
			Response: asc.AppStoreReviewDetailResponse{
				Data: asc.AppStoreReviewDetail{
					ID: "TEST",
				},
			},
		},
		// UpdateReviewDetail
		response{
			Response: asc.AppStoreReviewDetailResponse{
				Data: asc.AppStoreReviewDetail{
					ID: "TEST",
				},
			},
		},
		// ListAttachmentsForReviewDetail
		response{
			Response: asc.AppStoreReviewAttachmentsResponse{
				Data: []asc.AppStoreReviewAttachment{},
			},
		},
		// CreateAttachment
		response{
			Response: asc.AppStoreReviewAttachmentResponse{
				Data: asc.AppStoreReviewAttachment{
					Attributes: &asc.AppStoreReviewAttachmentAttributes{
						FileName: &asset.Name,
						FileSize: &asset.Size,
						UploadOperations: []asc.UploadOperation{
							{
								Length: asc.Int(int(asset.Size)),
								Method: asc.String("PATCH"),
								Offset: asc.Int(0),
								URL:    asc.String(attachmentURL.String()),
							},
						},
					},
					ID: attachmentID,
				},
			},
		},
		// UploadOperation #1
		response{
			RawResponse: `{}`,
		},
		// CommitAttachment
		response{
			Response: asc.AppStoreReviewAttachmentResponse{},
		},
	)

	err = client.UpdateReviewDetails(ctx.Context, testID, config.ReviewDetails{
		Contact:     &config.ContactPerson{},
		DemoAccount: &config.DemoAccount{},
		Attachments: []config.File{
			{Path: asset.Name},
		},
	})
	assert.NoError(t, err)
}

func TestUpdateReviewDetails_ErrUpdate(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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

// Test EnablePhasedRelease

func TestEnablePhasedRelease_Update(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			Response: asc.AppStoreVersionPhasedReleaseResponse{
				Data: asc.AppStoreVersionPhasedRelease{
					ID: testID,
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.EnablePhasedRelease(ctx.Context, testID)
	assert.NoError(t, err)
}

func TestEnablePhasedRelease_Create(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.EnablePhasedRelease(ctx.Context, testID)
	assert.NoError(t, err)
}

// Test SubmitApp

func TestSubmitApp_Happy(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
