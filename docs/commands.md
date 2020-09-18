---
layout: page
nav_order: 7
---

# Command Reference

## Global Options

You can provide the `--debug` flag on any of the `cider` commands in order to glean more logging information.

## `init`

Use to initialize a new Cider project. This will create a new configuration file in the current directory that should be checked into source control. Use the `--config` flag to provide a customized configuration file path. For example:

```shell
cider check
```

## `release`

Run with the appropriate flags and environment variables to process one or more apps defined in your configuration apps. For example:

```shell
env \
    ASC_KEY_ID="..." \
    ASC_ISSUER_ID="..." \
    ASC_PRIVATE_KEY_PATH="..." \
  cider release --mode="appstore" --set-version="1.0"
```

You can provide a path to a project directory as an argument to be the root directory of all relative path expansions in the program, such as the Git repository, preview sets, and screenshot resources. An exception to this is if you set a custom configuration file path with the `--config` flag.

Use multiple `--app` flags, each one set to a key in the `apps` map in your configuration file corresponding to an app you wish to process. You can also use `--all-apps` or `-A` to select all apps. You can omit this flag if your `apps` map has only one app defined.

The `--mode` flag is used to declare the publishing mode for submission. The default is `"testflight"` for submitting to Testflight, and the other alternative option is `"appstore"` for submitting to the App Store.

The `--set-version` is used to override the version string used when creating new App Store versions or querying the API. It should correspond to a version string used in your app builds, and it should follow semver. If this flag is omitted, Git will be leveraged to determine the latest tag. The tag will be used to calculate the version string under the same constraints. 

You can also use the `--set-build` flag to override the specific build you want to operate on. The default behavior without this flag is to select the latest build. In both cases, if the selected build has an invalid processing state, Cider will abort with an error to ensure your release is handled safely. 

The `--set-beta-group` and `--set-beta-tester` flags can be invoked repeatedly for any number of beta group names or beta tester emails you want to include. Use of these flags will totally override the beta group and beta tester configurations used for all apps included during a single process. For `--set-beta-group`, provide a name of a beta group that already exists in App Store Connect, or that is declared in the top-level TestFlight settings of your configuration file. For `--set-beta-tester`, provide the email address of a beta tester that has already been added as a beta tester on your App Store Connect team, or that has been declared in the top-level TestFlight settings of your configuration file. 

The release command provides a variety of skip flags:

- `--skip-git`: Skips deriving version information from Git. Must only be used in conjunction with `--set-version`.
- `--skip-submit`: Skips submitting for review
- `--skip-update-metadata`: Skips updating metadata (app info, localizations, assets, review details, etc.)
- `--skip-update-pricing`: Skips updating app pricing

Finally, you can provide a `--timeout` flag with a duration value to place a limit on the runtime of the release process. While Cider is intended to be fast, apps with a lot of localization fields or many apps will take a while longer to run through all of their metadata. Also relevant is the `--max-processes` flag, which will run certain metadata syncing and asset uploading logic in parallel, limited by the number provided with the flag. 

### Environment

In addition to flags, Cider requires a few environment variables to be set in order to operate. They each correspond to an element of authorization described by [Creating API Keys for App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi/creating_api_keys_for_app_store_connect_api) from the Apple Developer Documentation.

- `ASC_KEY_ID` – Your key ID
- `ASC_ISSUER_ID` – Your team's issuer ID
- `ASC_PRIVATE_KEY` – Your .p8 private key issued by Apple
  - You can alternatively use `ASC_PRIVATE_KEY_PATH` if you'd rather supply a path to a key instead.

These three values each have varying degrees of sensetivity and should be treated as secrets. Store them securely in your environment so Cider can leverage them safely.

## `check`

Use to validate your configuration file. Use the `--config` flag to provide a customized configuration file path. For example:

```shell
cider check
```

## `help`

Invoke to get help information for the entire program, or `cider help [command]` to get help information for a specific command. For example:

```shell
cider help release
```
