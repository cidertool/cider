---
layout: page
nav_order: 4
---

# Configuration

The project configuration can be written in YAML, or JSON inside a YAML file. Map merges are allowed.

```yaml
# This file contains all available configuration options and demonstrates
# the sorts of values that they can accept.

# Map of app names (used for logging) to app configurations. (required)
'My App':
  # App bundle ID. (required)
  id: com.myproject.MyApp
  # App primary locale. Primary language is used if this is not set.
  primaryLocale: en-US
  # Disclosure whether the app uses third-party content that you are
  # permitted to use. Set to null to unset the value from App Store Connect.
  usesThirdPartyContent: false
  # Information about pricing and territory availability
  availability:
    # Indicates whether or not the app should be made automaticaly available
    # in new App Store territories, as Apple makes new ones available.
    availableInNewTerritories: false
    # Array of price schedules that describe the pricing details of your app.
    priceTiers:
      - # Refers to an app pricing tier on the App Store Pricing Matrix. For example,
        # Tier 1 should be represented as "1" and the Free tier should be represented
        # as "0". (required)
        #
        # https://appstoreconnect.apple.com/apps/pricingmatrix
        tier: '0'
        # The start date a price schedule should take effect. Set to null to
        # have it take effect immediately.
        startDate: null
        # The end date a price schedule should be in effect until. Field is currently a no-op.
        endDate: null
    # Array of ISO 3166-1 Alpha-3 country codes corresponding to territories to make your app available in.
    territories:
      - USA
  # The App Store categories your app belongs to. A primary category is required, and a secondary category
  # is encouraged.
  #
  # Some categories have optional subcategories you can use to improve the specificity of your categorization.
  # Up to two subcategories can provided each for the primary and secondary categories.
  #
  # See the [App Categories](#app-categories) section below for more information on app categories
  categories:
    # ID for the primary category. (required)
    primary: SOCIAL_NETWORKING
    # IDs of any subcategories to apply to the primary category. Only up to two will be accepted.
    primarySubcategories: []
    # ID for the secondary category.
    secondary: GAMES
    # IDs of any subcategories to apply to the secondary category. Only up to two will be accepted.
    secondarySubcategories:
      - GAMES_SIMULATION
      - GAMES_RACING
  # Declaration of age rating qualifiers or content warnings. These are used in your App Store listing,
  # as well as in categorizing your app in lists when users browse.
  #
  # https://help.apple.com/app-store-connect/#/dev599d50efb
  ageRatings:
    # Whether your app enables legally and guideline-compliant gambling.
    gamblingAndContests: false
    # Whether your app enables generalized usage of the internet, such as an internet browser.
    unrestrictedWebAccess: false
    # Age band to use in categorizing your app for lists aimed at kids.
    # Valid values: "5 and under", "6-8", "9-11"
    kidsAgeBand: 9-11
    # Whether your app makes references to alcohol, tobacco, or drug use and/or paraphernalia.
    alcoholTobaccoOrDrugUseOrReferences: none
    # Whether your app offers medical advice or treatment information.
    medicalOrTreatmentInformation: none
    # Whether your app contains or enables profanity and/or crude humor.
    profanityOrCrudeHumor: none
    # Whether your app contains or enables sexual content or nudity.
    sexualContentOrNudity: none
    # Whether your app enables simulated gambling with either real or simulated currency.
    gamblingSimulated: none
    # Whether your app contains horror or fear-inducing themes.
    horrorOrFearThemes: infrequentOrMild
    # Whether your app contains mature or suggestive themes.
    matureOrSuggestiveThemes: none
    # Whether your app contains or enables sexual content or nudity that is graphic in nature.
    sexualContentGraphicAndNudity: none
    # Whether your app contains cartoon or fantasy violence.
    violenceCartoonOrFantasy: frequentOrIntense
    # Whether your app contains realistic violence.
    violenceRealistic: none
    # Whether your app contains prolonged, realistic violence that is graphic or sadistic in nature.
    violenceRealisticProlongedGraphicOrSadistic: none
  # Map of locale IDs to localization configurations for app information. (required)
  localizations:
    en-US:
      # App name in this locale. (required)
      name: ''
      # App subtitle.
      subtitle: ''
      # Privacy Policy textual content, if you don't have a URL.
      privacyPolicyText: ''
      # Privacy Policy URL, if you'd rather host your privacy policy text content.
      privacyPolicyURL: ''
  # Information to configure new App Store versions (required)
  versions:
    # Platform the app is to be released on. (required)
    # Valid values: "iOS", "tvOS", "macOS"
    platform: iOS
    # Map of locale IDs to localization configurations for App Store version information. (required)
    localizations:
      en-US:
        # App description used in App Store listing. (required)
        description: ''
        # Comma separated keywords used in App Store listing.
        keywords: ''
        # URL to point back to app's company or service website.
        marketingURL: ''
        # Promotional text used in App Store listing.
        promotionalText: ''
        # URL to point back to app's support website.
        supportURL: ''
        # Release note text used to inform users of new changes in updates.
        # Not required for the initial release.
        whatsNew: ''
        # Map of app preview types to arrays of app screenshot resources
        previewSets:
          iphone58:
            - # Path to an asset, relative to the current directory. (required)
              path: assets/store/iphone58/preview.mp4
              # MIME type of the asset. Overriding this is usually unnecessary.
              mimeType: video/mp4
              # Time code to a frame to show as a preview of the video, if not the beginning.
              previewFrameTimeCode: '0'
        # Map of screenshot display types to arrays of app screenshot resources
        screenshotSets:
          iphone58:
            - # Path to an asset, relative to the current directory. (required)
              path: assets/store/iphone58/home.png
            - # Path to an asset, relative to the current directory. (required)
              path: assets/store/iphone58/videochat.png
    # Copyright text
    copyright: ''
    # Earliest release date, in Go's RFC3339 format. Set to null to release
    # as soon as is permitted by the release type.
    earliestReleaseDate: null
    # Release type. Valid values: "manual", "afterApproval", "scheduled"
    releaseType: afterApproval
    # Indicates whether phased release is enabled
    enablePhasedRelease: true
    # Information about an app's IDFA declaration. Omit or set to null to declare to
    # Apple that your app does not use an IDFA.
    idfaDeclaration:
      # Indicates that the app attributes user action with previous ads
      attributesActionWithPreviousAd: false
      # Indicates that the app attributes user installation with previous ads
      attributesAppInstallationToPreviousAd: false
      # Indicates that the app developer will honor Apple's guidelines around tracking when
      # the user has chosen to limit ad tracking.
      honorsLimitedAdTracking: false
      # Indicates that the app serves ads
      servesAds: false
    # Routing coverage resource
    routingCoverage:
      # Path to an asset, relative to the current directory. (required)
      path: assets/store/coverage.json
    # Details about an app to share with the App Store reviewer.
    reviewDetails:
      # Point of contact for the App Store reviewer.
      contact:
        # Contact email (required)
        email: ''
        # Contact first (given) name (required)
        firstName: ''
        # Contact last (family) name (required)
        lastName: ''
        # Contact phone number (required)
        phone: ''
      # A demo account the reviewer can use to evaluate functionality
      demoAccount:
        # Whether or not a demo account is required. Other fields can be
        # omitted if this is set to false. (required)
        isRequired: false
        # Account name
        name: ''
        # Account password
        password: ''
      # Notes that the reviewer should be aware of
      notes: ''
      # Attachment resources the reviewer should be aware of or use in evaluation.
      attachments:
        - path: assets/review/test_image.png
  # Information to configure new Testflight beta releases (required)
  testflight:
    # Indicates whether to auto-notify existing beta testers of a new Testflight update. (required)
    enableAutoNotify: true
    # Beta license agreement content. (required)
    licenseAgreement: ''
    # Map of locale IDs to localization configurations for beta app and beta build information. (required)
    localizations:
      en-US:
        # App description used in App Store listing. (required)
        description: ''
        # Email address testers can use to provide feedback.
        feedbackEmail: ''
        # URL to point back to app's company or service website.
        marketingURL: ''
        # Privacy Policy URL.
        privacyPolicyURL: ''
        # Privacy Policy textual content for tvOS.
        tvOSPrivacyPolicy: ''
        # Release note text used to inform users of new changes in updates.
        # Not required for the initial release.
        whatsNew: ''
    # Array of beta group names. If you want to refer to beta groups defined
    # in this configuration file, use the value provided for the group field
    # on the corresponding beta group.
    # Beta groups to add or update in App Store Connect
    betaGroups:
      - # Name of the beta group (required)
        group: QA Team
      - # Name of the beta group (required)
        group: External Beta
        # Indicates whether a public link is enabled.
        publicLinkEnabled: true
        # Maximum number of testers that can join the beta group using the public link
        publicLinkLimit: 100
        # Indicates whether a public link limit is enabled.
        publicLinkLimitEnabled: true
        # Indicates whether a tester feedback is enabled within Testflight.
        feedbackEnabled: true
        # Array of beta testers to explicitly add to the beta group
        testers:
          - # Tester email (required)
            email: tester@email.com
            # Tester first (given) name
            firstName: Tester
            # Tester last (family) name
            lastName: Tester
    # Individual beta testers to add or update in App Store Connect
    betaTesters:
      - # Tester email (required)
        email: tester2@email.com
        # Tester first (given) name
        firstName: Tester
        # Tester last (family) name
        lastName: Two
    # Details about an app to share with the App Store reviewer.
    reviewDetails:
      # Point of contact for the App Store reviewer.
      contact:
        # Contact email (required)
        email: ''
        # Contact first (given) name (required)
        firstName: ''
        # Contact last (family) name (required)
        lastName: ''
        # Contact phone number (required)
        phone: ''
      # A demo account the reviewer can use to evaluate functionality
      demoAccount:
        # Whether or not a demo account is required. Other fields can be
        # omitted if this is set to false. (required)
        isRequired: false
        # Account name
        name: ''
        # Account password
        password: ''
      # Notes that the reviewer should be aware of
      notes: ''
      # Attachment resources the reviewer should be aware of or use in evaluation.
      attachments:
        - path: assets/review/test_image.png
```

### App Categories

App categories provided and supported by the App Store Connect API are fluid and difficult to create a consistent format for. The App Store adds categories regularly, and it represents a challenge for both metadata maintainers and maintainers of Cider to support. Therefore, the choice has been made to accept any string as a category ID, and let the API respond with whether or not it's valid.

Here are some known category IDs, with subcategories broken out where applicable, that you can use in your configuration:

- `BOOKS`
- `BUSINESS`
- `DEVELOPER_TOOLS`
- `EDUCATION`
- `ENTERTAINMENT`
- `FINANCE`
- `FOOD_AND_DRINK`
- `GAMES`
  - `GAMES_SPORTS`
  - `GAMES_WORD`
  - `GAMES_MUSIC`
  - `GAMES_ADVENTURE`
  - `GAMES_ACTION`
  - `GAMES_ROLE_PLAYING`
  - `GAMES_CASUAL`
  - `GAMES_BOARD`
  - `GAMES_TRIVIA`
  - `GAMES_CARD`
  - `GAMES_PUZZLE`
  - `GAMES_CASINO`
  - `GAMES_STRATEGY`
  - `GAMES_SIMULATION`
  - `GAMES_RACING`
  - `GAMES_FAMILY`
- `HEALTH_AND_FITNESS`
- `LIFESTYLE`
- `MAGAZINES_AND_NEWSPAPERS`
- `MEDICAL`
- `PRODUCTIVITY`
- `REFERENCE`
- `SHOPPING`
- `SOCIAL_NETWORKING`
- `SPORTS`
- `STICKERS`
  - `STICKERS_PLACES_AND_OBJECTS`
  - `STICKERS_EMOJI_AND_EXPRESSIONS`
  - `STICKERS_CELEBRATIONS`
  - `STICKERS_CELEBRITIES`
  - `STICKERS_MOVIES_AND_TV`
  - `STICKERS_SPORTS_AND_ACTIVITIES`
  - `STICKERS_EATING_AND_DRINKING`
  - `STICKERS_CHARACTERS`
  - `STICKERS_ANIMALS`
  - `STICKERS_FASHION`
  - `STICKERS_ART`
  - `STICKERS_GAMING`
  - `STICKERS_KIDS_AND_FAMILY`
  - `STICKERS_PEOPLE`
  - `STICKERS_MUSIC`
- `MUSIC`
- `TRAVEL`
- `UTILITIES`
- `WEATHER`

For more information on categories, see [Choosing a category](https://developer.apple.com/app-store/categories/) on the Apple Developer Portal.
