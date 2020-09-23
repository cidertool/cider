---
layout: page
nav_order: 4
---

# Configuration

Package config contains types and helpers to describe the configuration of an Cider project

## Specification

### Project

Project is the top level configuration type.  



### App

App outlines general information about your app, primarily for querying purposes.  

- [x] **id: string** – Bundle ID of the app.  
- [ ] **primaryLocale: string** – Primary locale of the app.  
- [ ] **usesThirdPartyContent: bool** – Whether or not the app uses third party content. Omit to avoid declarting content rights.  
- [ ] **availability: [Availability](#availability)** – Availability of the app, including pricing and supported territories.  
- [ ] **categories: [Categories](#categories)** – Categories to list under in the App Store.  
- [ ] **ageRatings: [AgeRatingDeclaration](#ageratingdeclaration)** – Content warnings that are used to declare the age rating.  
- [x] **localizations: [AppLocalizations](#applocalizations)** – App info localizations.  
- [x] **versions: [Version](#version)** – Metadata to configure new App Store versions.  
- [x] **testflight: [Testflight](#testflight)** – Metadata to configure new Testflight beta releases.  

### Availability

Availability wraps aspects of app availability, such as territories and pricing.  

- [ ] **availableInNewTerritories: bool** – Indicates whether or not the app should be made automaticaly available in new App Store territories, as Apple makes new ones available.  
- [ ] **priceTiers: [PriceSchedule]** – List of PriceSchedules that describe the pricing details of your app.  
- [ ] **territories: [string]** – Array of ISO 3166-1 Alpha-3 country codes corresponding to territories to make your app available in.  

### Categories

Categories describes the categories your app belongs to. A primary category is required, and a secondary category is encouraged. 

Some categories have optional subcategories you can use to improve the specificity of your categorization. Up to two subcategories can provided each for the primary and secondary categories. 

See the [App Categories](#app-categories) section below for more information on app categories.  

- [x] **primary: string** – ID for the primary category.  
- [ ] **primarySubcategories: [string]** – IDs of any subcategories to apply to the primary category. Only up to two will be accepted.  
- [ ] **secondary: string** – ID for the secondary category.  
- [ ] **secondarySubcategories: [string]** – IDs of any subcategories to apply to the secondary category. Only up to two will be accepted.  

### AgeRatingDeclaration

AgeRatingDeclaration describes the various content warnings you can provide or apply to your applications.  

- [ ] **gamblingAndContests: bool** – Whether your app enables legally and guideline-compliant gambling.  
- [ ] **unrestrictedWebAccess: bool** – Whether your app enables generalized usage of the internet, such as an internet browser.  
- [ ] **kidsAgeBand: string** – Age band to use in categorizing your app for lists aimed at kids.   Valid options: `"5 and under"`, `"6-8"`, `"9-11"`.
- [ ] **alcoholTobaccoOrDrugUseOrReferences: string** – Whether your app makes references to alcohol, tobacco, or drug use and/or paraphernalia.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **medicalOrTreatmentInformation: string** – Whether your app offers medical advice or treatment information.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **profanityOrCrudeHumor: string** – Whether your app contains or enables profanity and/or crude humor.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **sexualContentOrNudity: string** – Whether your app contains or enables sexual content or nudity.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **gamblingSimulated: string** – Whether your app enables simulated gambling with either real or simulated currency.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **horrorOrFearThemes: string** – Whether your app contains horror or fear-inducing themes.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **matureOrSuggestiveThemes: string** – Whether your app contains mature or suggestive themes.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **sexualContentGraphicAndNudity: string** – Whether your app contains or enables sexual content or nudity that is graphic in nature.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **violenceCartoonOrFantasy: string** – Whether your app contains cartoon or fantasy violence.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **violenceRealistic: string** – Whether your app contains realistic violence.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.
- [ ] **violenceRealisticProlongedGraphicOrSadistic: string** – Whether your app contains prolonged, realistic violence that is graphic or sadistic in nature.   Valid options: `"none"`, `"infrequentOrMild"`, `"frequentOrIntense"`.

### AppLocalizations

AppLocalizations is a map of locales to AppLocalization objects.  



### Version

Version outlines the general details of your app store version as it will be represented on the App Store.  

- [x] **platform: [Platform](#platform)** – Platform the app is to be released on.  
- [x] **localizations: [VersionLocalizations](#versionlocalizations)** – Map of locale IDs to localization configurations for App Store version information.  
- [ ] **copyright: string** – Copyright information to display on the listing. Templated.  
- [ ] **earliestReleaseDate: Time** – Earliest release date, in Go's RFC3339 format. Set to null to release as soon as is permitted by the release type.  
- [ ] **releaseType: string** – Release type.   Valid options: `"manual"`, `"afterApproval"`, `"scheduled"`.
- [ ] **enablePhasedRelease: bool** – Indicates whether phased release should be enabled for updates.  
- [ ] **idfaDeclaration: [IDFADeclaration](#idfadeclaration)** – Information about an app's IDFA declaration. Omit or set to null to declare to Apple that your app does not use the IDFA.  
- [ ] **routingCoverage: [File](#file)** – Routing coverage resource.  
- [ ] **reviewDetails: [ReviewDetails](#reviewdetails)** – Details about an app to share with the App Store reviewer.  

### Testflight

Testflight represents configuration for beta distribution of apps.  

- [x] **enableAutoNotify: bool** – Indicates whether to auto-notify existing beta testers of a new Testflight update.  
- [x] **licenseAgreement: string** – Beta license agreement content. Templated.  
- [x] **localizations: [TestflightLocalizations](#testflightlocalizations)** – Map of locale IDs to localization configurations for beta app and beta build information.  
- [ ] **betaGroups: [BetaGroup]** – Array of beta group names. If you want to refer to beta groups defined in this configuration file, use the value provided for the group field on the corresponding beta group. Beta groups to add or update in App Store Connect.  
- [ ] **betaTesters: [BetaTester]** – Individual beta testers to add or update in App Store Connect.  
- [ ] **reviewDetails: [ReviewDetails](#reviewdetails)** – Details about an app to share with the App Store reviewer.  

### PriceSchedule

PriceSchedule represents pricing availability information that an app should be immediately configured to.  

- [x] **tier: string** – Tier corresponds to a representation of a tier on the App Store Pricing Matrix. For example, Tier 1 should be represented as "1" and the Free tier should be represented as "0". 

https://appstoreconnect.apple.com/apps/pricingmatrix  
- [ ] **startDate: Time** – StartDate is the start date a price schedule should take effect. Set to nil to have it take effect immediately.  
- [ ] **endDate: Time** – EndDate is the end date a price schedule should be in effect until. Field is currently a no-op.  

### AppLocalization

AppLocalization contains localized details for your App Store listing.  

- [x] **name: string** – Name of the app in this locale. Templated.  
- [ ] **subtitle: string** – Subtitle of the app in this locale. Templated.  
- [ ] **privacyPolicyText: string** – Privacy policy text if not using a URL. Templated.  
- [ ] **privacyPolicyURL: string** – Privacy policy URL if not using a text body. Templated.  


### VersionLocalizations

VersionLocalizations is a map of locales to VersionLocalization objects.  



### IDFADeclaration

IDFADeclaration outlines regulatory information for Apple to use to handle your apps' use of tracking identifiers. Implicitly enables `usesIdfa` when creating an app store version.  

- [x] **attributesActionWithPreviousAd: bool** – Indicates that the app attributes user action with previous ads.  
- [x] **attributesAppInstallationToPreviousAd: bool** – Indicates that the app attributes user installation with previous ads.  
- [x] **honorsLimitedAdTracking: bool** – Indicates that the app developer will honor Apple's guidelines around tracking when the user has chosen to limit ad tracking.  
- [x] **servesAds: bool** – Indicates that the app serves ads  

### File

File refers to a file on disk by name.  

- [x] **path: string** – Path to a file on-disk. Templated.  

### ReviewDetails

ReviewDetails contains information for App Store reviewers to use in their assessment.  

- [ ] **contact: [ContactPerson](#contactperson)** – Point of contact for the App Store reviewer.  
- [ ] **demoAccount: [DemoAccount](#demoaccount)** – A demo account the reviewer can use to evaluate functionality  
- [ ] **notes: string** – Notes that the reviewer should be aware of. Templated.  
- [ ] **attachments: [File]** – Attachment resources the reviewer should be aware of or use in evaluation.  

### TestflightLocalizations

TestflightLocalizations is a map of locales to TestflightLocalization objects.  



### BetaGroup

BetaGroup describes a beta group in Testflight that should be kept in sync and used with this app.  

- [x] **group: string** – Name of the beta group.  
- [ ] **publicLinkEnabled: bool** – Indicates whether to enable the public link.  
- [ ] **publicLinkLimitEnabled: bool** – Indicates whether a limit on the number of testers who can use the public link is enabled.  
- [ ] **feedbackEnabled: bool** – Indicates whether tester feedback is enabled within TestFlight  
- [ ] **publicLinkLimit: int** – Maximum number of testers that can join the beta group using the public link.  
- [ ] **testers: [BetaTester]** – Array of beta testers to explicitly assign to the beta group.  

### BetaTester

BetaTester describes an individual beta tester that should have access to this app.  

- [x] **email: string** – Beta tester email.  
- [ ] **firstName: string** – Beta tester first (given) name.  
- [ ] **lastName: string** – Beta tester last (family) name.  

### VersionLocalization

VersionLocalization contains localized details for the listing of a specific version on the App Store.  

- [x] **description: string** – App description in this locale. Templated.  
- [ ] **keywords: string** – App keywords in this locale. Templated.  
- [ ] **marketingURL: string** – Marketing URL to use in this locale. Templated.  
- [ ] **promotionalText: string** – Promotional text to use in this locale. Can be updated without a requiring a new build. Templated.  
- [ ] **supportURL: string** – Support URL to use in this locale. Templated.  
- [ ] **whatsNew: string** – "Whats New" release note text to use in this locale. Templated.  
- [ ] **previewSets: [PreviewSets](#previewsets)** – Map of app preview types to arrays of app screenshot resources.  
- [ ] **screenshotSets: [ScreenshotSets](#screenshotsets)** – Map of screenshot display types to arrays of app screenshot resources.  

### ContactPerson

ContactPerson is a point of contact for App Store reviewers to reach out to in case of an issue.  

- [x] **email: string** – Contact email. Templated.  
- [x] **firstName: string** – Contact first (given) name. Templated.  
- [x] **lastName: string** – Contact last (family) name. Templated.  
- [x] **phone: string** – Contact phone number. Templated.  

### DemoAccount

DemoAccount contains account credentials for App Store reviewers to assess your apps.  

- [x] **isRequired: bool** – Whether or not a demo account is required. Other fields can be omitted if this is set to false.  
- [ ] **name: string** – Demo account name or login. Templated.  
- [ ] **password: string** – Demo account password. Templated.  

### TestflightLocalization

TestflightLocalization contains localized details for the listing of a specific build in the Testflight app.  

- [x] **description: string** – Beta build description in this locale. Templated.  
- [ ] **feedbackEmail: string** – Email for testers to provide feedback to in this locale. Templated.  
- [ ] **marketingURL: string** – Marketing URL to use in this locale. Templated.  
- [ ] **privacyPolicyURL: string** – Privacy policy URL to use in this locale. Templated.  
- [ ] **tvOSPrivacyPolicy: string** – Privacy policy text to use on tvOS in this locale. Templated.  
- [ ] **whatsNew: string** – "Whats New" release note text to use in this locale. Templated.  

### PreviewSets

PreviewSets is a map of preview types to []Preview slices.  

 Valid previewTypes:

- `"appleTV"`
- `"desktop"`
- `"ipad105"`
- `"ipad97"`
- `"ipadPro129"`
- `"ipadPro3Gen11"`
- `"ipadPro3Gen129"`
- `"iphone35"`
- `"iphone40"`
- `"iphone47"`
- `"iphone55"`
- `"iphone58"`
- `"iphone65"`
- `"watchSeries3"`
- `"watchSeries4"`

### ScreenshotSets

ScreenshotSets is a map of screenshot types to []File slices.  

 Valid screenshotTypes:

- `"appleTV"`
- `"desktop"`
- `"ipad105"`
- `"ipad97"`
- `"ipadPro129"`
- `"ipadPro3Gen11"`
- `"ipadPro3Gen129"`
- `"iphone35"`
- `"iphone40"`
- `"iphone47"`
- `"iphone55"`
- `"iphone58"`
- `"iphone65"`
- `"watchSeries3"`
- `"watchSeries4"`
- `"ipad105imessage"`
- `"ipad97imessage"`
- `"ipadPro129imessage"`
- `"ipadPro3Gen11imessage"`
- `"ipadPro3Gen129imessage"`
- `"iphone40imessage"`
- `"iphone47imessage"`
- `"iphone55imessage"`
- `"iphone58imessage"`
- `"iphone65imessage"`

### Preview

Preview is an expansion of File that defines a new app preview asset.  

- [x] **path: string** – Path to a file on-disk. Templated.  
- [ ] **mimeType: string** – MIME type of the asset. Overriding this is usually unnecessary.  
- [ ] **previewFrameTimeCode: string** – Time code to a frame to show as a preview of the video, if not the beginning.  

## App Categories

App categories provided and supported by the App Store Connect API are fluid and difficult to create a consistent format for. The App Store adds categories regularly, and it represents a challenge for both metadata maintainers and maintainers of Cider to support. Therefore, the choice has been made to accept any string as a category ID, and let the API respond with whether or not it's valid.

Here are some known category IDs, with subcategories broken out where applicable, that you can use in your configuration:

- `"BOOKS"`
- `"BUSINESS"`
- `"DEVELOPER_TOOLS"`
- `"EDUCATION"`
- `"ENTERTAINMENT"`
- `"FINANCE"`
- `"FOOD_AND_DRINK"`
- `"GAMES"`
  - `"GAMES_SPORTS"`
  - `"GAMES_WORD"`
  - `"GAMES_MUSIC"`
  - `"GAMES_ADVENTURE"`
  - `"GAMES_ACTION"`
  - `"GAMES_ROLE_PLAYING"`
  - `"GAMES_CASUAL"`
  - `"GAMES_BOARD"`
  - `"GAMES_TRIVIA"`
  - `"GAMES_CARD"`
  - `"GAMES_PUZZLE"`
  - `"GAMES_CASINO"`
  - `"GAMES_STRATEGY"`
  - `"GAMES_SIMULATION"`
  - `"GAMES_RACING"`
  - `"GAMES_FAMILY"`
- `"HEALTH_AND_FITNESS"`
- `"LIFESTYLE"`
- `"MAGAZINES_AND_NEWSPAPERS"`
- `"MEDICAL"`
- `"PRODUCTIVITY"`
- `"REFERENCE"`
- `"SHOPPING"`
- `"SOCIAL_NETWORKING"`
- `"SPORTS"`
- `"STICKERS"`
  - `"STICKERS_PLACES_AND_OBJECTS"`
  - `"STICKERS_EMOJI_AND_EXPRESSIONS"`
  - `"STICKERS_CELEBRATIONS"`
  - `"STICKERS_CELEBRITIES"`
  - `"STICKERS_MOVIES_AND_TV"`
  - `"STICKERS_SPORTS_AND_ACTIVITIES"`
  - `"STICKERS_EATING_AND_DRINKING"`
  - `"STICKERS_CHARACTERS"`
  - `"STICKERS_ANIMALS"`
  - `"STICKERS_FASHION"`
  - `"STICKERS_ART"`
  - `"STICKERS_GAMING"`
  - `"STICKERS_KIDS_AND_FAMILY"`
  - `"STICKERS_PEOPLE"`
  - `"STICKERS_MUSIC"`
- `"MUSIC"`
- `"TRAVEL"`
- `"UTILITIES"`
- `"WEATHER"`

For more information on categories, see [Choosing a category](https://developer.apple.com/app-store/categories/) on the Apple Developer Portal.
