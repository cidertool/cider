# Config

Package config contains types and helpers to describe the configuration of an Cider project


### App

App outlines general information about your app, primarily for querying purposes.  

- [x] **id: string** - Bundle ID of the app.  
- [ ] **primaryLocale: string** - Primary locale of the app.  
- [ ] **usesThirdPartyContent: bool** - Whether or not the app uses third party content. Omit to avoid declarting content rights.  
- [ ] **availability: [Availability](#availability)** - Availability of the app, including pricing and supported territories.  
- [ ] **categories: [Categories](#categories)** - Categories to list under in the App Store.  
- [ ] **ageRatings: [AgeRatingDeclaration](#ageratingdeclaration)** - Content warnings that are used to declare the age rating.  
- [x] **localizations: [AppLocalizations](#applocalizations)** - App info localizations.  
- [x] **versions: [Version](#version)** - App version metadata.  
- [x] **testflight: [Testflight](#testflight)** - TestFlight metadata.  

### Availability

Availability wraps aspects of app availability, such as territories and pricing.  

- [ ] **availableInNewTerritories: bool** - AvailableInNewTerritories refers to whether or not the app should be made automaticaly available in new App Store territories, as Apple makes new ones available.  
- [ ] **priceTiers: [PriceSchedule]** - Pricing is a list of PriceSchedules that describe the pricing details of your app.  
- [ ] **territories: [string]** - Territories corresponds to the IDs of territories as they're referred to in App Store Connect. 

https://help.apple.com/app-store-connect/#/dev997f9cf7c  

### Categories

Categories describes the categories used for classificiation in the App Store.  

- [x] **primary: string** -  
- [ ] **primarySubcategories: [string]** -  
- [ ] **secondary: string** -  
- [ ] **secondarySubcategories: [string]** -  

### AgeRatingDeclaration

AgeRatingDeclaration describes the various content warnings you can provide or apply to your applications.  

- [ ] **gamblingAndContests: bool** -  
- [ ] **unrestrictedWebAccess: bool** -  
- [ ] **kidsAgeBand: [KidsAgeBand](#kidsageband)** -  
- [ ] **alcoholTobaccoOrDrugUseOrReferences: [ContentIntensity](#contentintensity)** -  
- [ ] **medicalOrTreatmentInformation: [ContentIntensity](#contentintensity)** -  
- [ ] **profanityOrCrudeHumor: [ContentIntensity](#contentintensity)** -  
- [ ] **sexualContentOrNudity: [ContentIntensity](#contentintensity)** -  
- [ ] **gamblingSimulated: [ContentIntensity](#contentintensity)** -  
- [ ] **horrorOrFearThemes: [ContentIntensity](#contentintensity)** -  
- [ ] **matureOrSuggestiveThemes: [ContentIntensity](#contentintensity)** -  
- [ ] **sexualContentGraphicAndNudity: [ContentIntensity](#contentintensity)** -  
- [ ] **violenceCartoonOrFantasy: [ContentIntensity](#contentintensity)** -  
- [ ] **violenceRealistic: [ContentIntensity](#contentintensity)** -  
- [ ] **violenceRealisticProlongedGraphicOrSadistic: [ContentIntensity](#contentintensity)** -  


### Version

Version outlines the general details of your app store version as it will be represented on the App Store.  

- [x] **platform: [Platform](#platform)** -  
- [x] **localizations: [VersionLocalizations](#versionlocalizations)** -  
- [ ] **copyright: string** - Copyright information to display on the listing. Templated.  
- [ ] **earliestReleaseDate: Time** -  
- [ ] **releaseType: [ReleaseType](#releasetype)** -  
- [ ] **enablePhasedRelease: bool** -  
- [ ] **idfaDeclaration: [IDFADeclaration](#idfadeclaration)** -  
- [ ] **routingCoverage: [File](#file)** -  
- [ ] **reviewDetails: [ReviewDetails](#reviewdetails)** -  

### Testflight

Testflight represents configuration for beta distribution of apps.  

- [x] **enableAutoNotify: bool** -  
- [x] **licenseAgreement: string** - Beta license agreement text. Templated.  
- [x] **localizations: [TestflightLocalizations](#testflightlocalizations)** -  
- [ ] **betaGroups: [BetaGroup]** -  
- [ ] **betaTesters: [BetaTester]** -  
- [ ] **reviewDetails: [ReviewDetails](#reviewdetails)** -  

### PriceSchedule

PriceSchedule represents pricing availability information that an app should be immediately configured to.  

- [x] **tier: string** - Tier corresponds to a representation of a tier on the App Store Pricing Matrix. For example, Tier 1 should be represented as "1" and the Free tier should be represented as "0". 

https://appstoreconnect.apple.com/apps/pricingmatrix  
- [ ] **startDate: Time** - StartDate is the start date a price schedule should take effect. Set to nil to have it take effect immediately.  
- [ ] **endDate: Time** - EndDate is the end date a price schedule should be in effect until. Field is currently a no-op.  



### AppLocalization

AppLocalization contains localized details for your App Store listing.  

- [x] **name: string** - Name of the app in this locale. Templated.  
- [ ] **subtitle: string** - Subtitle of the app in this locale. Templated.  
- [ ] **privacyPolicyText: string** - Privacy policy text if not using a URL. Templated.  
- [ ] **privacyPolicyURL: string** - Privacy policy URL if not using a text body. Templated.  




### IDFADeclaration

IDFADeclaration outlines regulatory information for Apple to use to handle your apps' use of tracking identifiers. Implicitly enables `usesIdfa` when creating an app store version.  

- [x] **attributesActionWithPreviousAd: bool** -  
- [x] **attributesAppInstallationToPreviousAd: bool** -  
- [x] **honorsLimitedAdTracking: bool** -  
- [x] **servesAds: bool** -  

### File

File refers to a file on disk by name.  

- [x] **path: string** - Path to a file on-disk. Templated.  

### ReviewDetails

ReviewDetails contains information for App Store reviewers to use in their assessment.  

- [ ] **contact: [ContactPerson](#contactperson)** -  
- [ ] **demoAccount: [DemoAccount](#demoaccount)** -  
- [ ] **notes: string** - Notes for the reviewer. Templated.  
- [ ] **attachments: [File]** -  


### BetaGroup

BetaGroup describes a beta group in Testflight that should be kept in sync and used with this app.  

- [x] **group: string** - Beta group name.  
- [ ] **publicLinkEnabled: bool** -  
- [ ] **publicLinkLimitEnabled: bool** -  
- [ ] **feedbackEnabled: bool** -  
- [ ] **publicLinkLimit: int** -  
- [ ] **testers: [BetaTester]** -  

### BetaTester

BetaTester describes an individual beta tester that should have access to this app.  

- [x] **email: string** - Beta tester email.  
- [ ] **firstName: string** - Beta tester first name.  
- [ ] **lastName: string** - Beta tester last name.  

### VersionLocalization

VersionLocalization contains localized details for the listing of a specific version on the App Store.  

- [x] **description: string** - App description in this locale. Templated.  
- [ ] **keywords: string** - App keywords in this locale. Templated.  
- [ ] **marketingURL: string** - Marketing URL to use in this locale. Templated.  
- [ ] **promotionalText: string** - Promotional text to use in this locale. Can be updated without a requiring a new build. Templated.  
- [ ] **supportURL: string** - Support URL to use in this locale. Templated.  
- [ ] **whatsNew: string** - "Whats New" release note text to use in this locale. Templated.  
- [ ] **previewSets: [PreviewSets](#previewsets)** -  
- [ ] **screenshotSets: [ScreenshotSets](#screenshotsets)** -  

### ContactPerson

ContactPerson is a point of contact for App Store reviewers to reach out to in case of an issue.  

- [x] **email: string** - Contact email. Required. Templated.  
- [x] **firstName: string** - Contact first name. Required. Templated.  
- [x] **lastName: string** - Contact last name. Required. Templated.  
- [x] **phone: string** - Contact phone number. Required. Templated.  

### DemoAccount

DemoAccount contains account credentials for App Store reviewers to assess your apps.  

- [x] **isRequired: bool** -  
- [ ] **name: string** - Demo account name or login. Templated.  
- [ ] **password: string** - Demo account password. Templated.  

### TestflightLocalization

TestflightLocalization contains localized details for the listing of a specific build in the Testflight app.  

- [x] **description: string** - Beta build description in this locale. Templated.  
- [ ] **feedbackEmail: string** - Email for testers to provide feedback to in this locale. Templated.  
- [ ] **marketingURL: string** - Marketing URL to use in this locale. Templated.  
- [ ] **privacyPolicyURL: string** - Privacy policy URL to use in this locale. Templated.  
- [ ] **tvOSPrivacyPolicy: string** - Privacy policy text to use on tvOS in this locale. Templated.  
- [ ] **whatsNew: string** - "Whats New" release note text to use in this locale. Templated.  




### Preview

Preview is an expansion of File that defines a new app preview asset.  

- [ ] **mimeType: string** -  
- [ ] **previewFrameTimeCode: string** -  


