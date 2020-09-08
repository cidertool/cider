package template

import (
	"testing"
	"time"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
	"github.com/stretchr/testify/assert"
)

const badTemplatePattern = "{{ .projectName "

func TestTemplatePipeHeader(t *testing.T) {
	pipe := Pipe{}
	assert.Equal(t, "applying template values", pipe.String())
}

func TestTemplateEmptyProjectPasses(t *testing.T) {
	ctx := context.New(config.Project{})
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, ctx.RawConfig, ctx.Config)
}

func TestTemplateFullyDefinedProjectWithGoodTemplatesPasses(t *testing.T) {
	ctx := context.New(fullyPopulatedProject())
	expected := ctx.Config.Name
	pipe := Pipe{}
	err := pipe.Run(ctx)

	assert.NoError(t, err)
	assert.NotEqual(t, ctx.RawConfig, ctx.Config)
	assert.Equal(t, expected, ctx.Config.Name)
	for _, app := range ctx.Config.Apps {
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

func TestTemplateWithBadAppLocalizationName(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Localizations["en-US"]
	loc.Name = badTemplatePattern
	proj.Apps["First"].Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppLocalizationSubtitle(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Localizations["en-US"]
	loc.Subtitle = badTemplatePattern
	proj.Apps["First"].Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppLocalizationPrivacyPolicyText(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Localizations["en-US"]
	loc.PrivacyPolicyText = badTemplatePattern
	proj.Apps["First"].Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithAppLocalizationPrivacyPolicyURL(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Localizations["en-US"]
	loc.PrivacyPolicyURL = badTemplatePattern
	proj.Apps["First"].Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionsCopyright(t *testing.T) {
	proj := fullyPopulatedProject()
	app := proj.Apps["First"]
	app.Versions.Copyright = badTemplatePattern
	proj.Apps["First"] = app
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionLocalizationDescription(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Versions.Localizations["en-US"]
	loc.Description = badTemplatePattern
	proj.Apps["First"].Versions.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionLocalizationKeywords(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Versions.Localizations["en-US"]
	loc.Keywords = badTemplatePattern
	proj.Apps["First"].Versions.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithAppVersionLocalizationMarketingURL(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Versions.Localizations["en-US"]
	loc.MarketingURL = badTemplatePattern
	proj.Apps["First"].Versions.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionLocalizationPromotionalText(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Versions.Localizations["en-US"]
	loc.PromotionalText = badTemplatePattern
	proj.Apps["First"].Versions.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionLocalizationSupportURL(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Versions.Localizations["en-US"]
	loc.SupportURL = badTemplatePattern
	proj.Apps["First"].Versions.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionLocalizationWhatsNewText(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Versions.Localizations["en-US"]
	loc.WhatsNewText = badTemplatePattern
	proj.Apps["First"].Versions.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionLocalizationPreviewSets(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.Localizations["en-US"].PreviewSets[config.PreviewTypeDesktop][0].Path = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionLocalizationScreenshotSets(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.Localizations["en-US"].ScreenshotSets[config.ScreenshotTypeDesktop][0].Path = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionReviewDetailsEmail(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.Contact.Email = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionReviewDetailsFirstName(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.Contact.FirstName = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionReviewDetailsLastName(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.Contact.LastName = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithAppVersionReviewDetailsPhone(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.Contact.Phone = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionReviewDetailsAccountName(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.DemoAccount.Name = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionReviewDetailsAccountPassword(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.DemoAccount.Password = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionReviewDetailsNotes(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.Notes = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionReviewDetailsAttachment(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.ReviewDetails.Attachments[0].Path = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppVersionRoutingCoveragePath(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Versions.RoutingCoverage.Path = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightLicenseAgreement(t *testing.T) {
	proj := fullyPopulatedProject()
	app := proj.Apps["First"]
	app.Testflight.LicenseAgreement = badTemplatePattern
	proj.Apps["First"] = app
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightLocalizationDescription(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Testflight.Localizations["en-US"]
	loc.Description = badTemplatePattern
	proj.Apps["First"].Testflight.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightLocalizationFeedbackEmail(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Testflight.Localizations["en-US"]
	loc.FeedbackEmail = badTemplatePattern
	proj.Apps["First"].Testflight.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightLocalizationMarketingURL(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Testflight.Localizations["en-US"]
	loc.MarketingURL = badTemplatePattern
	proj.Apps["First"].Testflight.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightLocalizationPrivacyPolicyURL(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Testflight.Localizations["en-US"]
	loc.PrivacyPolicyURL = badTemplatePattern
	proj.Apps["First"].Testflight.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightLocalizationTVOSPrivacyPolicy(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Testflight.Localizations["en-US"]
	loc.TVOSPrivacyPolicy = badTemplatePattern
	proj.Apps["First"].Testflight.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightLocalizationWhatsNewText(t *testing.T) {
	proj := fullyPopulatedProject()
	loc := proj.Apps["First"].Testflight.Localizations["en-US"]
	loc.WhatsNew = badTemplatePattern
	proj.Apps["First"].Testflight.Localizations["en-US"] = loc
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsEmail(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.Contact.Email = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsFirstName(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.Contact.FirstName = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsLastName(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.Contact.LastName = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsPhoneNumber(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.Contact.Phone = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsAccountName(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.DemoAccount.Name = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsAccountPassword(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.DemoAccount.Password = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsNotes(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.Notes = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func TestTemplateWithBadAppTestflightReviewDetailsAttachments(t *testing.T) {
	proj := fullyPopulatedProject()
	proj.Apps["First"].Testflight.ReviewDetails.Attachments[0].Path = badTemplatePattern
	ctx := context.New(proj)
	pipe := Pipe{}
	err := pipe.Run(ctx)
	assert.Error(t, err)
}

func fullyPopulatedProject() config.Project {
	return config.Project{
		Name: "My Project",
		Testflight: config.Testflight{
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
		},
		Apps: map[string]config.App{
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
						Name:              "{{ .projectName }}",
						Subtitle:          "{{ .projectName }}",
						PrivacyPolicyText: "{{ .projectName }}",
						PrivacyPolicyURL:  "{{ .projectName }}",
					},
					"ja": {
						Name:              "{{ .projectName }}",
						Subtitle:          "{{ .projectName }}",
						PrivacyPolicyText: "{{ .projectName }}",
						PrivacyPolicyURL:  "{{ .projectName }}",
					},
				},
				Versions: config.Version{
					Platform: "",
					Localizations: map[string]config.VersionLocalization{
						"en-US": {
							Description:     "{{ .projectName }}",
							Keywords:        "{{ .projectName }}",
							MarketingURL:    "{{ .projectName }}",
							PromotionalText: "{{ .projectName }}",
							SupportURL:      "{{ .projectName }}",
							WhatsNewText:    "{{ .projectName }}",
							PreviewSets: config.PreviewSets{
								config.PreviewTypeDesktop: {
									{
										File: config.File{
											Path: "{{ .projectName }}",
										},
										MIMEType:             "image/jpg",
										PreviewFrameTimeCode: "0",
									},
								},
							},
							ScreenshotSets: config.ScreenshotSets{
								config.ScreenshotTypeDesktop: {
									{
										Path: "{{ .projectName }}",
									},
								},
							},
						},
						"ja": {
							Description:     "{{ .projectName }}",
							Keywords:        "{{ .projectName }}",
							MarketingURL:    "{{ .projectName }}",
							PromotionalText: "{{ .projectName }}",
							SupportURL:      "{{ .projectName }}",
							WhatsNewText:    "{{ .projectName }}",
							PreviewSets: config.PreviewSets{
								config.PreviewTypeDesktop: {
									{
										File: config.File{
											Path: "{{ .projectName }}",
										},
										MIMEType:             "image/jpg",
										PreviewFrameTimeCode: "0",
									},
								},
							},
							ScreenshotSets: config.ScreenshotSets{
								config.ScreenshotTypeDesktop: {
									{
										Path: "{{ .projectName }}",
									},
								},
							},
						},
					},
					Copyright:            "{{ .projectName }}",
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
						Path: "{{ .projectName }}",
					},
					ReviewDetails: &config.ReviewDetails{
						Contact: &config.ContactPerson{
							Email:     "{{ .projectName }}",
							FirstName: "{{ .projectName }}",
							LastName:  "{{ .projectName }}",
							Phone:     "{{ .projectName }}",
						},
						DemoAccount: &config.DemoAccount{
							Required: false,
							Name:     "{{ .projectName }}",
							Password: "{{ .projectName }}",
						},
						Notes: "{{ .projectName }}",
						Attachments: []config.File{
							{
								Path: "{{ .projectName }}",
							},
						},
					},
				},
				Testflight: config.TestflightForApp{
					EnableAutoNotify: false,
					LicenseAgreement: "{{ .projectName }}",
					Localizations: map[string]config.TestflightLocalization{
						"en-US": {
							Description:       "{{ .projectName }}",
							FeedbackEmail:     "{{ .projectName }}",
							MarketingURL:      "{{ .projectName }}",
							PrivacyPolicyURL:  "{{ .projectName }}",
							TVOSPrivacyPolicy: "{{ .projectName }}",
							WhatsNew:          "{{ .projectName }}",
						},
						"ja": {
							Description:       "{{ .projectName }}",
							FeedbackEmail:     "{{ .projectName }}",
							MarketingURL:      "{{ .projectName }}",
							PrivacyPolicyURL:  "{{ .projectName }}",
							TVOSPrivacyPolicy: "{{ .projectName }}",
							WhatsNew:          "{{ .projectName }}",
						},
					},
					BetaGroups: []string{
						"Jeff's Team",
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
							Email:     "{{ .projectName }}",
							FirstName: "{{ .projectName }}",
							LastName:  "{{ .projectName }}",
							Phone:     "{{ .projectName }}",
						},
						DemoAccount: &config.DemoAccount{
							Required: false,
							Name:     "{{ .projectName }}",
							Password: "{{ .projectName }}",
						},
						Notes: "{{ .projectName }}",
						Attachments: []config.File{
							{
								Path: "{{ .projectName }}",
							},
						},
					},
				},
			},
		},
	}
}
