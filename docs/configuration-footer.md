## App Categories

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
