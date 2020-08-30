# Functions

## sync org settings

- beta groups
- beta testers
- app info localizations

## publish beta release

1. create or update beta app localizations
1. create or update app encryption details
1. select build matching pattern
1. update beta build details
1. create or update beta build localizations and assign to build
1. create or update beta license agreements
1. assign beta groups and beta testers to build
1. create beta app submission for build

## publish app release

1. create or update app info localizations
1. create (or update, if resuming) app store version
1. create or update app store version localizations
1. select build matching pattern
1. update idfa declaration for app store version
1. upload assets
   - NOTE: If assets already exist on ASC and there is not a local asset present to replace it, warn but do not fail
   - routing coverages
   - previews
   - screenshots
1. create or update app store review details
1. create app store version submission
