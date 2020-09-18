// Package config contains types and helpers to describe the configuration of an Cider project
package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"gopkg.in/yaml.v2"
)

// Platform represents a supported platform type from App Store Connect.
type Platform string

const (
	// PlatformiOS refers to the iOS platform.
	PlatformiOS Platform = "iOS"
	// PlatformMacOS refers to the macOS platform.
	PlatformMacOS Platform = "macOS"
	// PlatformTvOS refers to the tvOS platform.
	PlatformTvOS Platform = "tvOS"
)

type contentIntensity string

const (
	// ContentIntensityNone refers to the NONE content warning.
	ContentIntensityNone contentIntensity = "none"
	// ContentIntensityInfrequentOrMild refers to the INFREQUENT_OR_MILD content warning.
	ContentIntensityInfrequentOrMild contentIntensity = "infrequentOrMild"
	// ContentIntensityFrequentOrIntense refers to the FREQUENT_OR_INTENSE content warning.
	ContentIntensityFrequentOrIntense contentIntensity = "frequentOrIntense"
)

type kidsAgeBand string

const (
	// KidsAgeBandFiveAndUnder refers to the FIVE_AND_UNDER kids age band.
	KidsAgeBandFiveAndUnder kidsAgeBand = "5 and under"
	// KidsAgeBandSixToEight refers to the SIX_TO_EIGHT kids age band.
	KidsAgeBandSixToEight kidsAgeBand = "6-8"
	// KidsAgeBandNineToEleven refers to the NINE_TO_ELEVEN kids age band.
	KidsAgeBandNineToEleven kidsAgeBand = "9-11"
)

type releaseType string

const (
	// ReleaseTypeManual refers to a manual release type.
	ReleaseTypeManual releaseType = "manual"
	// ReleaseTypeAfterApproval refers to an automatic release type to be completed after approval.
	ReleaseTypeAfterApproval releaseType = "afterApproval"
	// ReleaseTypeScheduled refers to an automatic release type to be completed after a scheduled date.
	ReleaseTypeScheduled releaseType = "scheduled"
)

// File refers to a file on disk by name.
type File struct {
	// Path to a file on-disk. Templated.
	Path string `yaml:"path"`
}

// Preview is an expansion of File that defines a new app preview asset.
type Preview struct {
	File                 `yaml:",inline"`
	MIMEType             string `yaml:"mimeType,omitempty"`
	PreviewFrameTimeCode string `yaml:"previewFrameTimeCode,omitempty"`
}

type previewType string

const (
	// PreviewTypeAppleTV is a preview type for Apple TV.
	PreviewTypeAppleTV previewType = "appleTV"
	// PreviewTypeDesktop is a preview type for Desktop.
	PreviewTypeDesktop previewType = "desktop"
	// PreviewTypeiPad105 is a preview type for iPad 10.5.
	PreviewTypeiPad105 previewType = "ipad105"
	// PreviewTypeiPad97 is a preview type for iPad 9.7.
	PreviewTypeiPad97 previewType = "ipad97"
	// PreviewTypeiPadPro129 is a preview type for iPad Pro 12.9.
	PreviewTypeiPadPro129 previewType = "ipadPro129"
	// PreviewTypeiPadPro3Gen11 is a preview type for third-generation iPad Pro 11.
	PreviewTypeiPadPro3Gen11 previewType = "ipadPro3Gen11"
	// PreviewTypeiPadPro3Gen129 is a preview type for third-generation iPad Pro 12.9.
	PreviewTypeiPadPro3Gen129 previewType = "ipadPro3Gen129"
	// PreviewTypeiPhone35 is a preview type for iPhone 3.5.
	PreviewTypeiPhone35 previewType = "iphone35"
	// PreviewTypeiPhone40 is a preview type for iPhone 4.0.
	PreviewTypeiPhone40 previewType = "iphone40"
	// PreviewTypeiPhone47 is a preview type for iPhone 4.7.
	PreviewTypeiPhone47 previewType = "iphone47"
	// PreviewTypeiPhone55 is a preview type for iPhone 5.5.
	PreviewTypeiPhone55 previewType = "iphone55"
	// PreviewTypeiPhone58 is a preview type for iPhone 5.8.
	PreviewTypeiPhone58 previewType = "iphone58"
	// PreviewTypeiPhone65 is a preview type for iPhone 6.5.
	PreviewTypeiPhone65 previewType = "iphone65"
	// PreviewTypeWatchSeries3 is a preview type for Watch Series 3.
	PreviewTypeWatchSeries3 previewType = "watchSeries3"
	// PreviewTypeWatchSeries4 is a preview type for Watch Series 4.
	PreviewTypeWatchSeries4 previewType = "watchSeries4"
)

type screenshotType string

const (
	// ScreenshotTypeAppleTV is a screenshot type for Apple TV.
	ScreenshotTypeAppleTV screenshotType = "appleTV"
	// ScreenshotTypeDesktop is a screenshot type for Desktop.
	ScreenshotTypeDesktop screenshotType = "desktop"
	// ScreenshotTypeiPad105 is a screenshot type for iPad 10.5.
	ScreenshotTypeiPad105 screenshotType = "ipad105"
	// ScreenshotTypeiPad97 is a screenshot type for iPad 9.7.
	ScreenshotTypeiPad97 screenshotType = "ipad97"
	// ScreenshotTypeiPadPro129 is a screenshot type for iPad Pro 12.9.
	ScreenshotTypeiPadPro129 screenshotType = "ipadPro129"
	// ScreenshotTypeiPadPro3Gen11 is a screenshot type for third-generation iPad Pro 11.
	ScreenshotTypeiPadPro3Gen11 screenshotType = "ipadPro3Gen11"
	// ScreenshotTypeiPadPro3Gen129 is a screenshot type for third-generation iPad Pro 12.9.
	ScreenshotTypeiPadPro3Gen129 screenshotType = "ipadPro3Gen129"
	// ScreenshotTypeiPhone35 is a screenshot type for iPhone 3.5.
	ScreenshotTypeiPhone35 screenshotType = "iphone35"
	// ScreenshotTypeiPhone40 is a screenshot type for iPhone 4.0.
	ScreenshotTypeiPhone40 screenshotType = "iphone40"
	// ScreenshotTypeiPhone47 is a screenshot type for iPhone 4.7.
	ScreenshotTypeiPhone47 screenshotType = "iphone47"
	// ScreenshotTypeiPhone55 is a screenshot type for iPhone 5.5.
	ScreenshotTypeiPhone55 screenshotType = "iphone55"
	// ScreenshotTypeiPhone58 is a screenshot type for iPhone 5.8.
	ScreenshotTypeiPhone58 screenshotType = "iphone58"
	// ScreenshotTypeiPhone65 is a screenshot type for iPhone 6.5.
	ScreenshotTypeiPhone65 screenshotType = "iphone65"
	// ScreenshotTypeWatchSeries3 is a screenshot type for Watch Series 3.
	ScreenshotTypeWatchSeries3 screenshotType = "watchSeries3"
	// ScreenshotTypeWatchSeries4 is a screenshot type for Watch Series 4.
	ScreenshotTypeWatchSeries4 screenshotType = "watchSeries4"
	// ScreenshotTypeiMessageiPad105 is a screenshot type for iMessage apps on iPad 10.5.
	ScreenshotTypeiMessageiPad105 screenshotType = "ipad105imessage"
	// ScreenshotTypeiMessageiPad97 is a screenshot type for iMessage apps on iPad 9.7.
	ScreenshotTypeiMessageiPad97 screenshotType = "ipad97imessage"
	// ScreenshotTypeiMessageiPadPro129 is a screenshot type for iMessage apps on iPad Pro 12.9.
	ScreenshotTypeiMessageiPadPro129 screenshotType = "ipadPro129imessage"
	// ScreenshotTypeiMessageiPadPro3Gen11 is a screenshot type for iMessage apps on third-generation iPad Pro 11.
	ScreenshotTypeiMessageiPadPro3Gen11 screenshotType = "ipadPro3Gen11imessage"
	// ScreenshotTypeiMessageiPadPro3Gen129 is a screenshot type for iMessage apps on third-generation iPad Pro 12.9.
	ScreenshotTypeiMessageiPadPro3Gen129 screenshotType = "ipadPro3Gen129imessage"
	// ScreenshotTypeiMessageiPhone40 is a screenshot type for iMessage apps on iPhone 4.0.
	ScreenshotTypeiMessageiPhone40 screenshotType = "iphone40imessage"
	// ScreenshotTypeiMessageiPhone47 is a screenshot type for iMessage apps on iPhone 4.7.
	ScreenshotTypeiMessageiPhone47 screenshotType = "iphone47imessage"
	// ScreenshotTypeiMessageiPhone55 is a screenshot type for iMessage apps on iPhone 5.5.
	ScreenshotTypeiMessageiPhone55 screenshotType = "iphone55imessage"
	// ScreenshotTypeiMessageiPhone58 is a screenshot type for iMessage apps on iPhone 5.8.
	ScreenshotTypeiMessageiPhone58 screenshotType = "iphone58imessage"
	// ScreenshotTypeiMessageiPhone65 is a screenshot type for iMessage apps on iPhone 6.5.
	ScreenshotTypeiMessageiPhone65 screenshotType = "iphone65imessage"
)

// Repo represents any kind of repo (github, gitlab, etc).
type Repo struct {
	Owner string `yaml:",omitempty"`
	Name  string `yaml:",omitempty"`
}

// Project is the top level configuration type.
type Project struct {
	// Name of the project, used for logging.
	Name       string         `yaml:"name"`
	Testflight Testflight     `yaml:"testflight"`
	Apps       map[string]App `yaml:"apps"`
}

// Testflight represents information about a Testflight configuration for an entire App Store Connect team.
type Testflight struct {
	BetaGroups  []BetaGroup  `yaml:"betaGroups"`
	BetaTesters []BetaTester `yaml:"betaTesters"`
}

// App outlines general information about your app, primarily for querying purposes.
type App struct {
	// Bundle ID of the app. Required.
	BundleID string `yaml:"id"`
	// Primary locale of the app.
	PrimaryLocale         string                `yaml:"primaryLocale,omitempty"`
	UsesThirdPartyContent *bool                 `yaml:"usesThirdPartyContent,omitempty"`
	Availability          *Availability         `yaml:"availability,omitempty"`
	Categories            *Categories           `yaml:"categories,omitempty"`
	AgeRatingDeclaration  *AgeRatingDeclaration `yaml:"ageRatings,omitempty"`
	Localizations         AppLocalizations      `yaml:"localizations"`
	Versions              Version               `yaml:"versions"`
	Testflight            TestflightForApp      `yaml:"testflight"`
}

// Categories describes the categories used for classificiation in the App Store.
type Categories struct {
	Primary                string    `yaml:"primary"`
	PrimarySubcategories   [2]string `yaml:"primarySubcategories"`
	Secondary              string    `yaml:"secondary,omitempty"`
	SecondarySubcategories [2]string `yaml:"secondarySubcategories"`
}

// AgeRatingDeclaration describes the various content warnings you can provide or apply to your applications.
type AgeRatingDeclaration struct {
	GamblingAndContests                         *bool             `yaml:"gamblingAndContests,omitempty"`
	UnrestrictedWebAccess                       *bool             `yaml:"unrestrictedWebAccess,omitempty"`
	KidsAgeBand                                 *kidsAgeBand      `yaml:"kidsAgeBand,omitempty"`
	AlcoholTobaccoOrDrugUseOrReferences         *contentIntensity `yaml:"alcoholTobaccoOrDrugUseOrReferences,omitempty"`
	MedicalOrTreatmentInformation               *contentIntensity `yaml:"medicalOrTreatmentInformation,omitempty"`
	ProfanityOrCrudeHumor                       *contentIntensity `yaml:"profanityOrCrudeHumor,omitempty"`
	SexualContentOrNudity                       *contentIntensity `yaml:"sexualContentOrNudity,omitempty"`
	GamblingSimulated                           *contentIntensity `yaml:"gamblingSimulated,omitempty"`
	HorrorOrFearThemes                          *contentIntensity `yaml:"horrorOrFearThemes,omitempty"`
	MatureOrSuggestiveThemes                    *contentIntensity `yaml:"matureOrSuggestiveThemes,omitempty"`
	SexualContentGraphicAndNudity               *contentIntensity `yaml:"sexualContentGraphicAndNudity,omitempty"`
	ViolenceCartoonOrFantasy                    *contentIntensity `yaml:"violenceCartoonOrFantasy,omitempty"`
	ViolenceRealistic                           *contentIntensity `yaml:"violenceRealistic,omitempty"`
	ViolenceRealisticProlongedGraphicOrSadistic *contentIntensity `yaml:"violenceRealisticProlongedGraphicOrSadistic,omitempty"`
}

// Availability wraps aspects of app availability, such as territories and pricing.
type Availability struct {
	// AvailableInNewTerritories refers to whether or not the app should be made automaticaly available
	// in new App Store territories, as Apple makes new ones available.
	AvailableInNewTerritories *bool `yaml:"availableInNewTerritories,omitempty"`
	// Pricing is a list of PriceSchedules that describe the pricing details of your app.
	Pricing []PriceSchedule `yaml:"priceTiers,omitempty"`
	// Territories corresponds to the IDs of territories as they're referred to in App Store Connect.
	//
	// https://help.apple.com/app-store-connect/#/dev997f9cf7c
	Territories []string `yaml:"territories,omitempty"`
}

// PriceSchedule represents pricing availability information that an app should be immediately
// configured to.
type PriceSchedule struct {
	// Tier corresponds to a representation of a tier on the App Store Pricing Matrix.
	// For example, Tier 1 should be represented as "1" and the Free tier should be
	// represented as "0".
	//
	// https://appstoreconnect.apple.com/apps/pricingmatrix
	Tier string `yaml:"tier"`
	// StartDate is the start date a price schedule should take effect. Set to nil to have it take
	// effect immediately.
	StartDate *time.Time `yaml:"startDate,omitempty"`
	// EndDate is the end date a price schedule should be in effect until. Field is currently a no-op.
	EndDate *time.Time `yaml:"endDate,omitempty"`
}

// AppLocalizations is a map of locales to AppLocalization objects.
type AppLocalizations map[string]AppLocalization

// AppLocalization contains localized details for your App Store listing.
type AppLocalization struct {
	// Name of the app in this locale. Templated.
	Name string `yaml:"name"`
	// Subtitle of the app in this locale. Templated.
	Subtitle string `yaml:"subtitle,omitempty"`
	// Privacy policy text if not using a URL. Templated.
	PrivacyPolicyText string `yaml:"privacyPolicyText,omitempty"`
	// Privacy policy URL if not using a text body. Templated.
	PrivacyPolicyURL string `yaml:"privacyPolicyURL,omitempty"`
}

// Version outlines the general details of your app store version as it will be represented
// on the App Store.
type Version struct {
	Platform      Platform             `yaml:"platform"`
	Localizations VersionLocalizations `yaml:"localizations"`
	// Copyright information to display on the listing. Templated.
	Copyright            string           `yaml:"copyright,omitempty"`
	EarliestReleaseDate  *time.Time       `yaml:"earliestReleaseDate,omitempty"`
	ReleaseType          releaseType      `yaml:"releaseType,omitempty"`
	PhasedReleaseEnabled bool             `yaml:"enablePhasedRelease,omitempty"`
	IDFADeclaration      *IDFADeclaration `yaml:"idfaDeclaration,omitempty"`
	RoutingCoverage      *File            `yaml:"routingCoverage,omitempty"`
	ReviewDetails        *ReviewDetails   `yaml:"reviewDetails,omitempty"`
}

// VersionLocalizations is a map of locales to VersionLocalization objects.
type VersionLocalizations map[string]VersionLocalization

// VersionLocalization contains localized details for the listing of a specific version on the App Store.
type VersionLocalization struct {
	// App description in this locale. Templated.
	Description string `yaml:"description"`
	// App keywords in this locale. Templated.
	Keywords string `yaml:"keywords,omitempty"`
	// Marketing URL to use in this locale. Templated.
	MarketingURL string `yaml:"marketingURL,omitempty"`
	// Promotional text to use in this locale. Can be updated without a requiring a new build. Templated.
	PromotionalText string `yaml:"promotionalText,omitempty"`
	// Support URL to use in this locale. Templated.
	SupportURL string `yaml:"supportURL,omitempty"`
	// "Whats New" release note text to use in this locale. Templated.
	WhatsNewText   string         `yaml:"whatsNew,omitempty"`
	PreviewSets    PreviewSets    `yaml:"previewSets,omitempty"`
	ScreenshotSets ScreenshotSets `yaml:"screenshotSets,omitempty"`
}

// PreviewSets is a map of preview types to []Preview slices.
type PreviewSets map[previewType][]Preview

// ScreenshotSets is a map of screenshot types to []File slices.
type ScreenshotSets map[screenshotType][]File

// IDFADeclaration outlines regulatory information for Apple to use to handle your apps' use
// of tracking identifiers. Implicitly enables `usesIdfa` when creating an app store version.
type IDFADeclaration struct {
	AttributesActionWithPreviousAd        bool `yaml:"attributesActionWithPreviousAd"`
	AttributesAppInstallationToPreviousAd bool `yaml:"attributesAppInstallationToPreviousAd"`
	HonorsLimitedAdTracking               bool `yaml:"honorsLimitedAdTracking"`
	ServesAds                             bool `yaml:"servesAds"`
}

// ReviewDetails contains information for App Store reviewers to use in their assessment.
type ReviewDetails struct {
	Contact     *ContactPerson `yaml:"contact,omitempty"`
	DemoAccount *DemoAccount   `yaml:"demoAccount,omitempty"`
	// Notes for the reviewer. Templated.
	Notes       string `yaml:"notes,omitempty"`
	Attachments []File `yaml:"attachments,omitempty"`
}

// ContactPerson is a point of contact for App Store reviewers to reach out to in case of an
// issue.
type ContactPerson struct {
	// Contact email. Required. Templated.
	Email string `yaml:"email"`
	// Contact first name. Required. Templated.
	FirstName string `yaml:"firstName"`
	// Contact last name. Required. Templated.
	LastName string `yaml:"lastName"`
	// Contact phone number. Required. Templated.
	Phone string `yaml:"phone"`
}

// DemoAccount contains account credentials for App Store reviewers to assess your apps.
type DemoAccount struct {
	Required bool `yaml:"isRequired"`
	// Demo account name or login. Templated.
	Name string `yaml:"name,omitempty"`
	// Demo account password. Templated.
	Password string `yaml:"password,omitempty"`
}

// TestflightForApp represents configuration for beta distribution of apps.
type TestflightForApp struct {
	EnableAutoNotify bool `yaml:"enableAutoNotify"`
	// Beta license agreement text. Templated.
	LicenseAgreement string                  `yaml:"licenseAgreement"`
	Localizations    TestflightLocalizations `yaml:"localizations"`
	BetaGroups       []string                `yaml:"betaGroups,omitempty"`
	BetaTesters      []BetaTester            `yaml:"betaTesters,omitempty"`
	ReviewDetails    *ReviewDetails          `yaml:"reviewDetails,omitempty"`
}

// TestflightLocalizations is a map of locales to TestflightLocalization objects.
type TestflightLocalizations map[string]TestflightLocalization

// BetaGroup describes a beta group in Testflight that should be kept in sync and used with this app.
type BetaGroup struct {
	// Beta group name.
	Name                  string       `yaml:"group"`
	EnablePublicLink      bool         `yaml:"publicLinkEnabled,omitempty"`
	EnablePublicLinkLimit bool         `yaml:"publicLinkLimitEnabled,omitempty"`
	FeedbackEnabled       bool         `yaml:"feedbackEnabled,omitempty"`
	PublicLinkLimit       int          `yaml:"publicLinkLimit,omitempty"`
	Testers               []BetaTester `yaml:"testers"`
}

// BetaTester describes an individual beta tester that should have access to this app.
type BetaTester struct {
	// Beta tester email.
	Email string `yaml:"email"`
	// Beta tester first name.
	FirstName string `yaml:"firstName,omitempty"`
	// Beta tester last name.
	LastName string `yaml:"lastName,omitempty"`
}

// TestflightLocalization contains localized details for the listing of a specific build in the Testflight app.
type TestflightLocalization struct {
	// Beta build description in this locale. Templated.
	Description string `yaml:"description"`
	// Email for testers to provide feedback to in this locale. Templated.
	FeedbackEmail string `yaml:"feedbackEmail,omitempty"`
	// Marketing URL to use in this locale. Templated.
	MarketingURL string `yaml:"marketingURL,omitempty"`
	// Privacy policy URL to use in this locale. Templated.
	PrivacyPolicyURL string `yaml:"privacyPolicyURL,omitempty"`
	// Privacy policy text to use on tvOS in this locale. Templated.
	TVOSPrivacyPolicy string `yaml:"tvOSPrivacyPolicy,omitempty"`
	// "Whats New" release note text to use in this locale. Templated.
	WhatsNew string `yaml:"whatsNew,omitempty"`
}

// Load config file.
func Load(file string) (config Project, err error) {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				log.Fatal(closeErr.Error())
			}
		}
	}()
	return LoadReader(f)
}

// LoadReader config via io.Reader.
func LoadReader(fd io.Reader) (config Project, err error) {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return config, err
	}
	err = yaml.UnmarshalStrict(data, &config)
	return config, err
}

func (p Project) String() (string, error) {
	b, err := yaml.Marshal(p)
	return string(b), err
}

// Copy performs a costly deep-copy of the entire project structure.
func (p *Project) Copy() (copy Project, err error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return copy, err
	}
	err = json.Unmarshal(bytes, &copy)
	return copy, err
}

// AppsMatching returns an array of keys in the Project matching the app names, or all names if the flag is set.
func (p *Project) AppsMatching(keys []string, shouldIncludeAll bool) []string {
	if shouldIncludeAll || len(p.Apps) == 1 {
		appNamesMatching := make([]string, len(p.Apps))
		i := 0
		for key := range p.Apps {
			appNamesMatching[i] = key
			i++
		}
		return appNamesMatching
	}
	appNamesMatching := make([]string, 0, len(keys))
	for _, key := range keys {
		if _, ok := p.Apps[key]; ok {
			appNamesMatching = append(appNamesMatching, key)
		}
	}
	return appNamesMatching
}

// APIValue returns the corresponding API value type for this config type.
func (p *Platform) APIValue() *asc.Platform {
	if p == nil {
		return nil
	}
	var value asc.Platform
	switch *p {
	case PlatformiOS:
		value = asc.PlatformIOS
	case PlatformMacOS:
		value = asc.PlatformMACOS
	case PlatformTvOS:
		value = asc.PlatformTVOS
	default:
		return nil
	}
	return &value
}

func (c *contentIntensity) APIValue() *string {
	if c == nil {
		return nil
	}
	var value string
	switch *c {
	case ContentIntensityNone:
		value = "NONE"
	case ContentIntensityInfrequentOrMild:
		value = "INFREQUENT_OR_MILD"
	case ContentIntensityFrequentOrIntense:
		value = "FREQUENT_OR_INTENSE"
	default:
		return nil
	}
	return &value
}

func (b *kidsAgeBand) APIValue() *asc.KidsAgeBand {
	if b == nil {
		return nil
	}
	var value asc.KidsAgeBand
	switch *b {
	case KidsAgeBandFiveAndUnder:
		value = asc.KidsAgeBandFiveAndUnder
	case KidsAgeBandSixToEight:
		value = asc.KidsAgeBandSixToEight
	case KidsAgeBandNineToEleven:
		value = asc.KidsAgeBandNineToEleven
	default:
		return nil
	}
	return &value
}

func (t *releaseType) APIValue() *string {
	if t == nil {
		return nil
	}
	var value string
	switch *t {
	case ReleaseTypeManual:
		value = "MANUAL"
	case ReleaseTypeAfterApproval:
		value = "AFTER_APPROVAL"
	case ReleaseTypeScheduled:
		value = "SCHEDULED"
	default:
		return nil
	}
	return &value
}

func (t *previewType) APIValue() *asc.PreviewType {
	if t == nil {
		return nil
	}
	var value asc.PreviewType
	switch *t {
	case PreviewTypeAppleTV:
		value = asc.PreviewTypeAppleTV
	case PreviewTypeDesktop:
		value = asc.PreviewTypeDesktop
	case PreviewTypeiPad105:
		value = asc.PreviewTypeiPad105
	case PreviewTypeiPad97:
		value = asc.PreviewTypeiPad97
	case PreviewTypeiPadPro129:
		value = asc.PreviewTypeiPadPro129
	case PreviewTypeiPadPro3Gen11:
		value = asc.PreviewTypeiPadPro3Gen11
	case PreviewTypeiPadPro3Gen129:
		value = asc.PreviewTypeiPadPro3Gen129
	case PreviewTypeiPhone35:
		value = asc.PreviewTypeiPhone35
	case PreviewTypeiPhone40:
		value = asc.PreviewTypeiPhone40
	case PreviewTypeiPhone47:
		value = asc.PreviewTypeiPhone47
	case PreviewTypeiPhone55:
		value = asc.PreviewTypeiPhone55
	case PreviewTypeiPhone58:
		value = asc.PreviewTypeiPhone58
	case PreviewTypeiPhone65:
		value = asc.PreviewTypeiPhone65
	case PreviewTypeWatchSeries3:
		value = asc.PreviewTypeWatchSeries3
	case PreviewTypeWatchSeries4:
		value = asc.PreviewTypeWatchSeries4
	default:
		return nil
	}
	return &value
}

func (t *screenshotType) APIValue() *asc.ScreenshotDisplayType {
	if t == nil {
		return nil
	}
	var value asc.ScreenshotDisplayType
	switch *t {
	case ScreenshotTypeAppleTV:
		value = asc.ScreenshotDisplayTypeAppAppleTV
	case ScreenshotTypeDesktop:
		value = asc.ScreenshotDisplayTypeAppDesktop
	case ScreenshotTypeiPad105:
		value = asc.ScreenshotDisplayTypeAppiPad105
	case ScreenshotTypeiPad97:
		value = asc.ScreenshotDisplayTypeAppiPad97
	case ScreenshotTypeiPadPro129:
		value = asc.ScreenshotDisplayTypeAppiPadPro129
	case ScreenshotTypeiPadPro3Gen11:
		value = asc.ScreenshotDisplayTypeAppiPadPro3Gen11
	case ScreenshotTypeiPadPro3Gen129:
		value = asc.ScreenshotDisplayTypeAppiPadPro3Gen129
	case ScreenshotTypeiPhone35:
		value = asc.ScreenshotDisplayTypeAppiPhone35
	case ScreenshotTypeiPhone40:
		value = asc.ScreenshotDisplayTypeAppiPhone40
	case ScreenshotTypeiPhone47:
		value = asc.ScreenshotDisplayTypeAppiPhone47
	case ScreenshotTypeiPhone55:
		value = asc.ScreenshotDisplayTypeAppiPhone55
	case ScreenshotTypeiPhone58:
		value = asc.ScreenshotDisplayTypeAppiPhone58
	case ScreenshotTypeiPhone65:
		value = asc.ScreenshotDisplayTypeAppiPhone65
	case ScreenshotTypeWatchSeries3:
		value = asc.ScreenshotDisplayTypeAppWatchSeries3
	case ScreenshotTypeWatchSeries4:
		value = asc.ScreenshotDisplayTypeAppWatchSeries4
	case ScreenshotTypeiMessageiPad105:
		value = asc.ScreenshotDisplayTypeiMessageAppIPad105
	case ScreenshotTypeiMessageiPad97:
		value = asc.ScreenshotDisplayTypeiMessageAppIPad97
	case ScreenshotTypeiMessageiPadPro129:
		value = asc.ScreenshotDisplayTypeiMessageAppIPadPro129
	case ScreenshotTypeiMessageiPadPro3Gen11:
		value = asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen11
	case ScreenshotTypeiMessageiPadPro3Gen129:
		value = asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen129
	case ScreenshotTypeiMessageiPhone40:
		value = asc.ScreenshotDisplayTypeiMessageAppIPhone40
	case ScreenshotTypeiMessageiPhone47:
		value = asc.ScreenshotDisplayTypeiMessageAppIPhone47
	case ScreenshotTypeiMessageiPhone55:
		value = asc.ScreenshotDisplayTypeiMessageAppIPhone55
	case ScreenshotTypeiMessageiPhone58:
		value = asc.ScreenshotDisplayTypeiMessageAppIPhone58
	case ScreenshotTypeiMessageiPhone65:
		value = asc.ScreenshotDisplayTypeiMessageAppIPhone65
	default:
		return nil
	}
	return &value
}

// GetPreviews fetches the value from the map corresponding to the API value.
func (s PreviewSets) GetPreviews(previewType asc.PreviewType) []Preview {
	switch previewType {
	case asc.PreviewTypeAppleTV:
		return s[PreviewTypeAppleTV]
	case asc.PreviewTypeDesktop:
		return s[PreviewTypeDesktop]
	case asc.PreviewTypeiPad105:
		return s[PreviewTypeiPad105]
	case asc.PreviewTypeiPad97:
		return s[PreviewTypeiPad97]
	case asc.PreviewTypeiPadPro129:
		return s[PreviewTypeiPadPro129]
	case asc.PreviewTypeiPadPro3Gen11:
		return s[PreviewTypeiPadPro3Gen11]
	case asc.PreviewTypeiPadPro3Gen129:
		return s[PreviewTypeiPadPro3Gen129]
	case asc.PreviewTypeiPhone35:
		return s[PreviewTypeiPhone35]
	case asc.PreviewTypeiPhone40:
		return s[PreviewTypeiPhone40]
	case asc.PreviewTypeiPhone47:
		return s[PreviewTypeiPhone47]
	case asc.PreviewTypeiPhone55:
		return s[PreviewTypeiPhone55]
	case asc.PreviewTypeiPhone58:
		return s[PreviewTypeiPhone58]
	case asc.PreviewTypeiPhone65:
		return s[PreviewTypeiPhone65]
	case asc.PreviewTypeWatchSeries3:
		return s[PreviewTypeWatchSeries3]
	case asc.PreviewTypeWatchSeries4:
		return s[PreviewTypeWatchSeries4]
	}
	return []Preview{}
}

// GetScreenshots fetches the value from the map corresponding to the API value.
func (s ScreenshotSets) GetScreenshots(screenshotType asc.ScreenshotDisplayType) []File {
	switch screenshotType {
	case asc.ScreenshotDisplayTypeAppAppleTV:
		return s[ScreenshotTypeAppleTV]
	case asc.ScreenshotDisplayTypeAppDesktop:
		return s[ScreenshotTypeDesktop]
	case asc.ScreenshotDisplayTypeAppiPad105:
		return s[ScreenshotTypeiPad105]
	case asc.ScreenshotDisplayTypeAppiPad97:
		return s[ScreenshotTypeiPad97]
	case asc.ScreenshotDisplayTypeAppiPadPro129:
		return s[ScreenshotTypeiPadPro129]
	case asc.ScreenshotDisplayTypeAppiPadPro3Gen11:
		return s[ScreenshotTypeiPadPro3Gen11]
	case asc.ScreenshotDisplayTypeAppiPadPro3Gen129:
		return s[ScreenshotTypeiPadPro3Gen129]
	case asc.ScreenshotDisplayTypeAppiPhone35:
		return s[ScreenshotTypeiPhone35]
	case asc.ScreenshotDisplayTypeAppiPhone40:
		return s[ScreenshotTypeiPhone40]
	case asc.ScreenshotDisplayTypeAppiPhone47:
		return s[ScreenshotTypeiPhone47]
	case asc.ScreenshotDisplayTypeAppiPhone55:
		return s[ScreenshotTypeiPhone55]
	case asc.ScreenshotDisplayTypeAppiPhone58:
		return s[ScreenshotTypeiPhone58]
	case asc.ScreenshotDisplayTypeAppiPhone65:
		return s[ScreenshotTypeiPhone65]
	case asc.ScreenshotDisplayTypeAppWatchSeries3:
		return s[ScreenshotTypeWatchSeries3]
	case asc.ScreenshotDisplayTypeAppWatchSeries4:
		return s[ScreenshotTypeWatchSeries4]
	case asc.ScreenshotDisplayTypeiMessageAppIPad105:
		return s[ScreenshotTypeiMessageiPad105]
	case asc.ScreenshotDisplayTypeiMessageAppIPad97:
		return s[ScreenshotTypeiMessageiPad97]
	case asc.ScreenshotDisplayTypeiMessageAppIPadPro129:
		return s[ScreenshotTypeiMessageiPadPro129]
	case asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen11:
		return s[ScreenshotTypeiMessageiPadPro3Gen11]
	case asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen129:
		return s[ScreenshotTypeiMessageiPadPro3Gen129]
	case asc.ScreenshotDisplayTypeiMessageAppIPhone40:
		return s[ScreenshotTypeiMessageiPhone40]
	case asc.ScreenshotDisplayTypeiMessageAppIPhone47:
		return s[ScreenshotTypeiMessageiPhone47]
	case asc.ScreenshotDisplayTypeiMessageAppIPhone55:
		return s[ScreenshotTypeiMessageiPhone55]
	case asc.ScreenshotDisplayTypeiMessageAppIPhone58:
		return s[ScreenshotTypeiMessageiPhone58]
	case asc.ScreenshotDisplayTypeiMessageAppIPhone65:
		return s[ScreenshotTypeiMessageiPhone65]
	}
	return []File{}
}
