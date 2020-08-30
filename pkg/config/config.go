package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/aaronsky/asc-go/asc"
	"gopkg.in/yaml.v2"
)

type platform string

const (
	// PlatformiOS refers to the iOS platform
	PlatformiOS platform = "iOS"
	// PlatformMacOS refers to the macOS platform
	PlatformMacOS platform = "macOS"
	// PlatformTvOS refers to the tvOS platform
	PlatformTvOS platform = "tvOS"
)

type releaseType string

const (
	// ReleaseTypeManual refers to a manual release type
	ReleaseTypeManual releaseType = "manual"
	// ReleaseTypeAfterApproval refers to an automatic release type to be completed after approval
	ReleaseTypeAfterApproval releaseType = "afterApproval"
	// ReleaseTypeScheduled refers to an automatic release type to be completed after a scheduled date
	ReleaseTypeScheduled releaseType = "scheduled"
)

type locale string

// File refers to a file on disk by name
type File struct {
	Path string `yaml:"path"`
}

// Preview is an expansion of File that defines a new app preview asset
type Preview struct {
	File                 `yaml:",inline"`
	MIMEType             string `yaml:"mimeType,omitempty"`
	PreviewFrameTimeCode string `yaml:"previewFrameTimeCode,omitempty"`
}

type previewType string

const (
	// PreviewTypeAppleTV is a preview type for Apple TV
	PreviewTypeAppleTV previewType = "appleTV"
	// PreviewTypeDesktop is a preview type for Desktop
	PreviewTypeDesktop previewType = "desktop"
	// PreviewTypeiPad105 is a preview type for iPad 10.5
	PreviewTypeiPad105 previewType = "ipad105"
	// PreviewTypeiPad97 is a preview type for iPad 9.7
	PreviewTypeiPad97 previewType = "ipad97"
	// PreviewTypeiPadPro129 is a preview type for iPad Pro 12.9
	PreviewTypeiPadPro129 previewType = "ipadPro129"
	// PreviewTypeiPadPro3Gen11 is a preview type for third-generation iPad Pro 11
	PreviewTypeiPadPro3Gen11 previewType = "ipadPro3Gen11"
	// PreviewTypeiPadPro3Gen129 is a preview type for third-generation iPad Pro 12.9
	PreviewTypeiPadPro3Gen129 previewType = "ipadPro3Gen129"
	// PreviewTypeiPhone35 is a preview type for iPhone 3.5
	PreviewTypeiPhone35 previewType = "iphone35"
	// PreviewTypeiPhone40 is a preview type for iPhone 4.0
	PreviewTypeiPhone40 previewType = "iphone40"
	// PreviewTypeiPhone47 is a preview type for iPhone 4.7
	PreviewTypeiPhone47 previewType = "iphone47"
	// PreviewTypeiPhone55 is a preview type for iPhone 5.5
	PreviewTypeiPhone55 previewType = "iphone55"
	// PreviewTypeiPhone58 is a preview type for iPhone 5.8
	PreviewTypeiPhone58 previewType = "iphone58"
	// PreviewTypeiPhone65 is a preview type for iPhone 6.5
	PreviewTypeiPhone65 previewType = "iphone65"
	// PreviewTypeWatchSeries3 is a preview type for Watch Series 3
	PreviewTypeWatchSeries3 previewType = "watchSeries3"
	// PreviewTypeWatchSeries4 is a preview type for Watch Series 4
	PreviewTypeWatchSeries4 previewType = "watchSeries4"
)

type screenshotType string

const (
	// ScreenshotTypeAppleTV is a screenshot type for Apple TV
	ScreenshotTypeAppleTV screenshotType = "appleTV"
	// ScreenshotTypeDesktop is a screenshot type for Desktop
	ScreenshotTypeDesktop screenshotType = "desktop"
	// ScreenshotTypeiPad105 is a screenshot type for iPad 10.5
	ScreenshotTypeiPad105 screenshotType = "ipad105"
	// ScreenshotTypeiPad97 is a screenshot type for iPad 9.7
	ScreenshotTypeiPad97 screenshotType = "ipad97"
	// ScreenshotTypeiPadPro129 is a screenshot type for iPad Pro 12.9
	ScreenshotTypeiPadPro129 screenshotType = "ipadPro129"
	// ScreenshotTypeiPadPro3Gen11 is a screenshot type for third-generation iPad Pro 11
	ScreenshotTypeiPadPro3Gen11 screenshotType = "ipadPro3Gen11"
	// ScreenshotTypeiPadPro3Gen129 is a screenshot type for third-generation iPad Pro 12.9
	ScreenshotTypeiPadPro3Gen129 screenshotType = "ipadPro3Gen129"
	// ScreenshotTypeiPhone35 is a screenshot type for iPhone 3.5
	ScreenshotTypeiPhone35 screenshotType = "iphone35"
	// ScreenshotTypeiPhone40 is a screenshot type for iPhone 4.0
	ScreenshotTypeiPhone40 screenshotType = "iphone40"
	// ScreenshotTypeiPhone47 is a screenshot type for iPhone 4.7
	ScreenshotTypeiPhone47 screenshotType = "iphone47"
	// ScreenshotTypeiPhone55 is a screenshot type for iPhone 5.5
	ScreenshotTypeiPhone55 screenshotType = "iphone55"
	// ScreenshotTypeiPhone58 is a screenshot type for iPhone 5.8
	ScreenshotTypeiPhone58 screenshotType = "iphone58"
	// ScreenshotTypeiPhone65 is a screenshot type for iPhone 6.5
	ScreenshotTypeiPhone65 screenshotType = "iphone65"
	// ScreenshotTypeWatchSeries3 is a screenshot type for Watch Series 3
	ScreenshotTypeWatchSeries3 screenshotType = "watchSeries3"
	// ScreenshotTypeWatchSeries4 is a screenshot type for Watch Series 4
	ScreenshotTypeWatchSeries4 screenshotType = "watchSeries4"
	// ScreenshotTypeiMessageiPad105 is a screenshot type for iMessage apps on iPad 10.5
	ScreenshotTypeiMessageiPad105 screenshotType = "ipad105imessage"
	// ScreenshotTypeiMessageiPad97 is a screenshot type for iMessage apps on iPad 9.7
	ScreenshotTypeiMessageiPad97 screenshotType = "ipad97imessage"
	// ScreenshotTypeiMessageiPadPro129 is a screenshot type for iMessage apps on iPad Pro 12.9
	ScreenshotTypeiMessageiPadPro129 screenshotType = "ipadPro129imessage"
	// ScreenshotTypeiMessageiPadPro3Gen11 is a screenshot type for iMessage apps on third-generation iPad Pro 11
	ScreenshotTypeiMessageiPadPro3Gen11 screenshotType = "ipadPro3Gen11imessage"
	// ScreenshotTypeiMessageiPadPro3Gen129 is a screenshot type for iMessage apps on third-generation iPad Pro 12.9
	ScreenshotTypeiMessageiPadPro3Gen129 screenshotType = "ipadPro3Gen129imessage"
	// ScreenshotTypeiMessageiPhone40 is a screenshot type for iMessage apps on iPhone 4.0
	ScreenshotTypeiMessageiPhone40 screenshotType = "iphone40imessage"
	// ScreenshotTypeiMessageiPhone47 is a screenshot type for iMessage apps on iPhone 4.7
	ScreenshotTypeiMessageiPhone47 screenshotType = "iphone47imessage"
	// ScreenshotTypeiMessageiPhone55 is a screenshot type for iMessage apps on iPhone 5.5
	ScreenshotTypeiMessageiPhone55 screenshotType = "iphone55imessage"
	// ScreenshotTypeiMessageiPhone58 is a screenshot type for iMessage apps on iPhone 5.8
	ScreenshotTypeiMessageiPhone58 screenshotType = "iphone58imessage"
	// ScreenshotTypeiMessageiPhone65 is a screenshot type for iMessage apps on iPhone 6.5
	ScreenshotTypeiMessageiPhone65 screenshotType = "iphone65imessage"
)

// Repo represents any kind of repo (github, gitlab, etc).
// to upload releases into.
type Repo struct {
	Owner string `yaml:",omitempty"`
	Name  string `yaml:",omitempty"`
}

// RepoRef represents any kind of repo which may differ
// from the one we are building from and may therefore
// also require separate authentication
// e.g. Homebrew Tap, Scoop bucket.
type RepoRef struct {
	Owner string `yaml:",omitempty"`
	Name  string `yaml:",omitempty"`
	Token string `yaml:",omitempty"`
}

// Project is the top level configuration type
type Project struct {
	Name       string         `yaml:"name"`
	Testflight Testflight     `yaml:"testflight"`
	Apps       map[string]App `yaml:"apps"`
}

// Testflight represents information about a Testflight configuration for an entire App Store Connect team
type Testflight struct {
	BetaGroups  []BetaGroup  `yaml:"betaGroups"`
	BetaTesters []BetaTester `yaml:"betaTesters"`
}

// App outlines general information about your app, primarily for querying purposes
type App struct {
	BundleID      string                     `yaml:"id"`
	Localizations map[locale]AppLocalization `yaml:"localizations"`
	Versions      Version                    `yaml:"versions"`
	Testflight    TestflightForApp           `yaml:"testflight"`
}

// AppLocalization contains localized details for your App Store listing.
type AppLocalization struct {
	Name              string `yaml:"name"`
	Subtitle          string `yaml:"subtitle,omitempty"`
	PrivacyPolicyText string `yaml:"privacyPolicyText,omitempty"`
	PrivacyPolicyURL  string `yaml:"privacyPolicyURL,omitempty"`
}

// Version outlines the general details of your app store version as it will be represented
// on the App Store.
type Version struct {
	Platform             platform                       `yaml:"platform"`
	Localizations        map[locale]VersionLocalization `yaml:"localizations"`
	Copyright            string                         `yaml:"copyright,omitempty"`
	EarliestReleaseDate  *time.Time                     `yaml:"earliestReleaseDate,omitempty"`
	ReleaseType          releaseType                    `yaml:"releaseType,omitempty"`
	PhasedReleaseEnabled bool                           `yaml:"enablePhasedRelease,omitempty"`
	IDFADeclaration      *IDFADeclaration               `yaml:"idfaDeclaration,omitempty"`
	RoutingCoverages     []File                         `yaml:"routingCoverages,omitempty"`
	ReviewDetails        *ReviewDetails                 `yaml:"reviewDetails,omitempty"`
}

// VersionLocalization contains localized details for the listing of a specific version on the App Store.
type VersionLocalization struct {
	Description     string                    `yaml:"description,omitempty"`
	Keywords        string                    `yaml:"keywords,omitempty"`
	MarketingURL    string                    `yaml:"marketingURL,omitempty"`
	PromotionalText string                    `yaml:"promotionalText,omitempty"`
	SupportURL      string                    `yaml:"supportURL,omitempty"`
	WhatsNewText    string                    `yaml:"whatsNew,omitempty"`
	PreviewSets     map[previewType][]Preview `yaml:"previewSets,omitempty"`
	ScreenshotSets  map[screenshotType][]File `yaml:"screenshotSets,omitempty"`
}

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
	Notes       string         `yaml:"notes,omitempty"`
	Attachments []File         `yaml:"attachments,omitempty"`
}

// ContactPerson is a point of contact for App Store reviewers to reach out to in case of an
// issue.
type ContactPerson struct {
	Email     string `yaml:"email,omitempty"`
	FirstName string `yaml:"firstName,omitempty"`
	LastName  string `yaml:"lastName,omitempty"`
	Phone     string `yaml:"phone,omitempty"`
}

// DemoAccount contains account credentials for App Store reviewers to assess your apps.
type DemoAccount struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Required bool   `yaml:"isRequired"`
}

// TestflightForApp represents configuration for beta distribution of apps.
type TestflightForApp struct {
	EnableAutoNotify bool                              `yaml:"enableAutoNotify"`
	BetaGroups       []string                          `yaml:"betaGroups"`
	BetaTesters      []BetaTester                      `yaml:"betaTesters"`
	Localizations    map[locale]TestflightLocalization `yaml:"localizations"`
	ReviewDetails    ReviewDetails                     `yaml:"reviewDetails"`
}

// BetaGroup describes a beta group in Testflight that should be kept in sync and used with this app.
type BetaGroup struct {
	Name                  string       `yaml:"group,omitempty"`
	EnablePublicLink      bool         `yaml:"publicLinkEnabled,omitempty"`
	PublicLinkLimit       bool         `yaml:"publicLinkLimit,omitempty"`
	EnablePublicLinkLimit bool         `yaml:"publicLinkLimitEnabled,omitempty"`
	FeedbackEnabled       bool         `yaml:"feedbackEnabled,omitempty"`
	Testers               []BetaTester `yaml:"testers,omitempty"`
}

// BetaTester describes an individual beta tester that should have access to this app.
type BetaTester struct {
	Email     string `yaml:"email,omitempty"`
	FirstName string `yaml:"firstName,omitempty"`
	LastName  string `yaml:"lastName,omitempty"`
}

// TestflightLocalization contains localized details for the listing of a specific build in the Testflight app.
type TestflightLocalization struct {
	Description       string `yaml:"description,omitempty"`
	FeedbackEmail     string `yaml:"feedbackEmail,omitempty"`
	MarketingURL      string `yaml:"marketingURL,omitempty"`
	PrivacyPolicyURL  string `yaml:"privacyPolicyURL,omitempty"`
	TVOSPrivacyPolicy string `yaml:"tvOSPrivacyPolicy,omitempty"`
	WhatsNew          string `yaml:"whatsNew,omitempty"`
}

// Load config file.
func Load(file string) (config Project, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
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

func (p platform) APIValue() (asc.Platform, error) {
	switch p {
	case PlatformiOS:
		return asc.PlatformIOS, nil
	case PlatformMacOS:
		return asc.PlatformMACOS, nil
	case PlatformTvOS:
		return asc.PlatformTVOS, nil
	}
	return asc.Platform(""), fmt.Errorf("could not convert platform %s to asc.Platform type", p)
}

func (t releaseType) APIValue() (string, error) {
	switch t {
	case ReleaseTypeManual:
		return "MANUAL", nil
	case ReleaseTypeAfterApproval:
		return "AFTER_APPROVAL", nil
	case ReleaseTypeScheduled:
		return "SCHEDULED", nil
	}
	return "", fmt.Errorf("could not convert releaseType %s to valid release type", t)
}

func (t previewType) APIValue() (asc.PreviewType, error) {
	switch t {
	case PreviewTypeAppleTV:
		return asc.PreviewTypeAppleTV, nil
	case PreviewTypeDesktop:
		return asc.PreviewTypeDesktop, nil
	case PreviewTypeiPad105:
		return asc.PreviewTypeiPad105, nil
	case PreviewTypeiPad97:
		return asc.PreviewTypeiPad97, nil
	case PreviewTypeiPadPro129:
		return asc.PreviewTypeiPadPro129, nil
	case PreviewTypeiPadPro3Gen11:
		return asc.PreviewTypeiPadPro3Gen11, nil
	case PreviewTypeiPadPro3Gen129:
		return asc.PreviewTypeiPadPro3Gen129, nil
	case PreviewTypeiPhone35:
		return asc.PreviewTypeiPhone35, nil
	case PreviewTypeiPhone40:
		return asc.PreviewTypeiPhone40, nil
	case PreviewTypeiPhone47:
		return asc.PreviewTypeiPhone47, nil
	case PreviewTypeiPhone55:
		return asc.PreviewTypeiPhone55, nil
	case PreviewTypeiPhone58:
		return asc.PreviewTypeiPhone58, nil
	case PreviewTypeiPhone65:
		return asc.PreviewTypeiPhone65, nil
	case PreviewTypeWatchSeries3:
		return asc.PreviewTypeWatchSeries3, nil
	case PreviewTypeWatchSeries4:
		return asc.PreviewTypeWatchSeries4, nil
	}
	return "", fmt.Errorf("could not convert previewType %s to valid asc.PreviewType", t)
}

func (t screenshotType) APIValue() (asc.ScreenshotDisplayType, error) {
	switch t {
	case ScreenshotTypeAppleTV:
		return asc.ScreenshotDisplayTypeAppAppleTV, nil
	case ScreenshotTypeDesktop:
		return asc.ScreenshotDisplayTypeAppDesktop, nil
	case ScreenshotTypeiPad105:
		return asc.ScreenshotDisplayTypeAppiPad105, nil
	case ScreenshotTypeiPad97:
		return asc.ScreenshotDisplayTypeAppiPad97, nil
	case ScreenshotTypeiPadPro129:
		return asc.ScreenshotDisplayTypeAppiPadPro129, nil
	case ScreenshotTypeiPadPro3Gen11:
		return asc.ScreenshotDisplayTypeAppiPadPro3Gen11, nil
	case ScreenshotTypeiPadPro3Gen129:
		return asc.ScreenshotDisplayTypeAppiPadPro3Gen129, nil
	case ScreenshotTypeiPhone35:
		return asc.ScreenshotDisplayTypeAppiPhone35, nil
	case ScreenshotTypeiPhone40:
		return asc.ScreenshotDisplayTypeAppiPhone40, nil
	case ScreenshotTypeiPhone47:
		return asc.ScreenshotDisplayTypeAppiPhone47, nil
	case ScreenshotTypeiPhone55:
		return asc.ScreenshotDisplayTypeAppiPhone55, nil
	case ScreenshotTypeiPhone58:
		return asc.ScreenshotDisplayTypeAppiPhone58, nil
	case ScreenshotTypeiPhone65:
		return asc.ScreenshotDisplayTypeAppiPhone65, nil
	case ScreenshotTypeWatchSeries3:
		return asc.ScreenshotDisplayTypeAppWatchSeries3, nil
	case ScreenshotTypeWatchSeries4:
		return asc.ScreenshotDisplayTypeAppWatchSeries4, nil
	case ScreenshotTypeiMessageiPad105:
		return asc.ScreenshotDisplayTypeiMessageAppIPad105, nil
	case ScreenshotTypeiMessageiPad97:
		return asc.ScreenshotDisplayTypeiMessageAppIPad97, nil
	case ScreenshotTypeiMessageiPadPro129:
		return asc.ScreenshotDisplayTypeiMessageAppIPadPro129, nil
	case ScreenshotTypeiMessageiPadPro3Gen11:
		return asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen11, nil
	case ScreenshotTypeiMessageiPadPro3Gen129:
		return asc.ScreenshotDisplayTypeiMessageAppIPadPro3Gen129, nil
	case ScreenshotTypeiMessageiPhone40:
		return asc.ScreenshotDisplayTypeiMessageAppIPhone40, nil
	case ScreenshotTypeiMessageiPhone47:
		return asc.ScreenshotDisplayTypeiMessageAppIPhone47, nil
	case ScreenshotTypeiMessageiPhone55:
		return asc.ScreenshotDisplayTypeiMessageAppIPhone55, nil
	case ScreenshotTypeiMessageiPhone58:
		return asc.ScreenshotDisplayTypeiMessageAppIPhone58, nil
	case ScreenshotTypeiMessageiPhone65:
		return asc.ScreenshotDisplayTypeiMessageAppIPhone65, nil
	}
	return "", fmt.Errorf("could not convert screenshotType %s to valid asc.ScreenshotDisplayType", t)
}
