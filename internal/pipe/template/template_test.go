package template

import (
	"errors"
	"testing"
	"time"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

const (
	goodTemplatePattern = "{{ .version }}"
	badTemplatePattern  = "{{ .version "
)

func TestTemplatePipeHeader(t *testing.T) {
	pipe := Pipe{}
	assert.Equal(t, "applying template values", pipe.String())
}

func TestTemplateEmptyProjectPasses(t *testing.T) {
	ctx := context.New(config.Project{})
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Equal(t, ctx.RawConfig, ctx.Config)
}

func TestTemplateFullyDefinedProjectWithGoodTemplatesPasses(t *testing.T) {
	ctx := context.New(fullyPopulatedProject(true))
	expected := ctx.Version

	pipe := Pipe{}
	err := pipe.Run(ctx)

	assert.NoError(t, err)
	assert.NotEqual(t, ctx.RawConfig, ctx.Config)

	for _, app := range ctx.Config {
		for _, loc := range app.Localizations {
			assert.Equal(t, expected, loc.Name)
			assert.Equal(t, expected, loc.Subtitle)
			assert.Equal(t, expected, loc.PrivacyPolicyText)
			assert.Equal(t, expected, loc.PrivacyPolicyURL)
		}

		assert.Equal(t, expected, app.Testflight.LicenseAgreement)

		for _, loc := range app.Testflight.Localizations {
			assert.Equal(t, expected, loc.Description)
			assert.Equal(t, expected, loc.FeedbackEmail)
			assert.Equal(t, expected, loc.MarketingURL)
			assert.Equal(t, expected, loc.PrivacyPolicyURL)
			assert.Equal(t, expected, loc.TVOSPrivacyPolicy)
			assert.Equal(t, expected, loc.WhatsNew)
		}

		assert.Equal(t, expected, app.Testflight.ReviewDetails.Contact.Email)
		assert.Equal(t, expected, app.Testflight.ReviewDetails.Contact.FirstName)
		assert.Equal(t, expected, app.Testflight.ReviewDetails.Contact.LastName)
		assert.Equal(t, expected, app.Testflight.ReviewDetails.Contact.Phone)
		assert.Equal(t, expected, app.Testflight.ReviewDetails.DemoAccount.Name)
		assert.Equal(t, expected, app.Testflight.ReviewDetails.DemoAccount.Password)
		assert.Equal(t, expected, app.Testflight.ReviewDetails.Notes)

		for _, file := range app.Testflight.ReviewDetails.Attachments {
			assert.Equal(t, expected, file.Path)
		}

		assert.Equal(t, expected, app.Versions.Copyright)

		for _, loc := range app.Versions.Localizations {
			assert.Equal(t, expected, loc.Description)
			assert.Equal(t, expected, loc.Keywords)
			assert.Equal(t, expected, loc.MarketingURL)
			assert.Equal(t, expected, loc.PromotionalText)
			assert.Equal(t, expected, loc.SupportURL)
			assert.Equal(t, expected, loc.WhatsNewText)

			for _, set := range loc.PreviewSets {
				for _, file := range set {
					assert.Equal(t, expected, file.Path)
				}
			}

			for _, set := range loc.ScreenshotSets {
				for _, file := range set {
					assert.Equal(t, expected, file.Path)
				}
			}
		}

		assert.Equal(t, expected, app.Versions.ReviewDetails.Contact.Email)
		assert.Equal(t, expected, app.Versions.ReviewDetails.Contact.FirstName)
		assert.Equal(t, expected, app.Versions.ReviewDetails.Contact.LastName)
		assert.Equal(t, expected, app.Versions.ReviewDetails.Contact.Phone)
		assert.Equal(t, expected, app.Versions.ReviewDetails.DemoAccount.Name)
		assert.Equal(t, expected, app.Versions.ReviewDetails.DemoAccount.Password)
		assert.Equal(t, expected, app.Versions.ReviewDetails.Notes)

		for _, file := range app.Versions.ReviewDetails.Attachments {
			assert.Equal(t, expected, file.Path)
		}

		assert.Equal(t, expected, app.Versions.RoutingCoverage.Path)
	}
}

func TestTemplateWithBadPatterns(t *testing.T) {
	proj := fullyPopulatedProject(false)
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)

	var merr *multierror.Error
	ok := errors.As(err, &merr)
	assert.True(t, ok)
	assert.NotNil(t, merr)
	assert.Equal(t, 55, merr.Len())
}

func fullyPopulatedProject(good bool) config.Project {
	var pattern string
	if good {
		pattern = goodTemplatePattern
	} else {
		pattern = badTemplatePattern
	}

	return config.Project{
		"First": {
			BundleID:              "com.app.bundleid",
			PrimaryLocale:         "en-US",
			UsesThirdPartyContent: asc.Bool(true),
			Availability: &config.Availability{
				AvailableInNewTerritories: asc.Bool(true),
				Pricing: []config.PriceSchedule{
					{
						Tier:      "1",
						StartDate: &time.Time{},
						EndDate:   &time.Time{},
					},
				},
				Territories: []string{
					"USA",
					"JPN",
				},
			},
			Localizations: map[string]config.AppLocalization{
				"en-US": {
					Name:              pattern,
					Subtitle:          pattern,
					PrivacyPolicyText: pattern,
					PrivacyPolicyURL:  pattern,
				},
				"ja": {
					Name:              pattern,
					Subtitle:          pattern,
					PrivacyPolicyText: pattern,
					PrivacyPolicyURL:  pattern,
				},
			},
			Versions: config.Version{
				Platform: "",
				Localizations: map[string]config.VersionLocalization{
					"en-US": {
						Description:     pattern,
						Keywords:        pattern,
						MarketingURL:    pattern,
						PromotionalText: pattern,
						SupportURL:      pattern,
						WhatsNewText:    pattern,
						PreviewSets: config.PreviewSets{
							config.PreviewTypeDesktop: {
								{
									File: config.File{
										Path: pattern,
									},
									MIMEType:             "image/jpg",
									PreviewFrameTimeCode: "0",
								},
							},
						},
						ScreenshotSets: config.ScreenshotSets{
							config.ScreenshotTypeDesktop: {
								{
									Path: pattern,
								},
							},
						},
					},
					"ja": {
						Description:     pattern,
						Keywords:        pattern,
						MarketingURL:    pattern,
						PromotionalText: pattern,
						SupportURL:      pattern,
						WhatsNewText:    pattern,
						PreviewSets: config.PreviewSets{
							config.PreviewTypeDesktop: {
								{
									File: config.File{
										Path: pattern,
									},
									MIMEType:             "image/jpg",
									PreviewFrameTimeCode: "0",
								},
							},
						},
						ScreenshotSets: config.ScreenshotSets{
							config.ScreenshotTypeDesktop: {
								{
									Path: pattern,
								},
							},
						},
					},
				},
				Copyright:            pattern,
				EarliestReleaseDate:  &time.Time{},
				ReleaseType:          config.ReleaseTypeAfterApproval,
				PhasedReleaseEnabled: false,
				IDFADeclaration: &config.IDFADeclaration{
					AttributesActionWithPreviousAd:        false,
					AttributesAppInstallationToPreviousAd: false,
					HonorsLimitedAdTracking:               false,
					ServesAds:                             false,
				},
				RoutingCoverage: &config.File{
					Path: pattern,
				},
				ReviewDetails: &config.ReviewDetails{
					Contact: &config.ContactPerson{
						Email:     pattern,
						FirstName: pattern,
						LastName:  pattern,
						Phone:     pattern,
					},
					DemoAccount: &config.DemoAccount{
						Required: false,
						Name:     pattern,
						Password: pattern,
					},
					Notes: pattern,
					Attachments: []config.File{
						{
							Path: pattern,
						},
					},
				},
			},
			Testflight: config.Testflight{
				EnableAutoNotify: false,
				LicenseAgreement: pattern,
				Localizations: map[string]config.TestflightLocalization{
					"en-US": {
						Description:       pattern,
						FeedbackEmail:     pattern,
						MarketingURL:      pattern,
						PrivacyPolicyURL:  pattern,
						TVOSPrivacyPolicy: pattern,
						WhatsNew:          pattern,
					},
					"ja": {
						Description:       pattern,
						FeedbackEmail:     pattern,
						MarketingURL:      pattern,
						PrivacyPolicyURL:  pattern,
						TVOSPrivacyPolicy: pattern,
						WhatsNew:          pattern,
					},
				},
				BetaGroups: []config.BetaGroup{
					{
						Name:                  "Jeff's Team",
						EnablePublicLink:      true,
						EnablePublicLinkLimit: true,
						FeedbackEnabled:       true,
						PublicLinkLimit:       100,
						Testers: []config.BetaTester{
							{
								Email:     "jeff@jeff.com",
								FirstName: "Jeff",
								LastName:  "Jefferson",
							},
						},
					},
				},
				BetaTesters: []config.BetaTester{
					{
						Email:     "jeff@jeff.com",
						FirstName: "Jeff",
						LastName:  "Jefferson",
					},
				},
				ReviewDetails: &config.ReviewDetails{
					Contact: &config.ContactPerson{
						Email:     pattern,
						FirstName: pattern,
						LastName:  pattern,
						Phone:     pattern,
					},
					DemoAccount: &config.DemoAccount{
						Required: false,
						Name:     pattern,
						Password: pattern,
					},
					Notes: pattern,
					Attachments: []config.File{
						{
							Path: pattern,
						},
					},
				},
			},
		},
	}
}
