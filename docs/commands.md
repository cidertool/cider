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
  cider release --mode="appstore" --version="1.0"
```

You can provide a path to a project directory as an argument to be the root directory of all relative path expansions in the program, such as the Git repository, preview sets, and screenshot resources. An exception to this is if you set a custom configuration file path with the `--config` flag.

Use multiple `--app` flags, each one set to a key in the `apps` map in your configuration file corresponding to an app you wish to process. You can also use `--all-apps` or `-A` to select all apps. You can omit this flag if your `apps` map has only one app defined.

The `--mode` flag is used to declare the publishing mode for submission. The default is `"testflight"` for submitting to Testflight, and the other alternative option is `"appstore"` for submitting to the App Store.

The `--version` is used to override the version string used when creating new App Store versions or querying the API. It should correspond to a version string used in your app builds, and it should follow semver. If this flag is omitted, Git will be leveraged to determine the latest tag. The tag will be used to calculate the version string under the same constraints.

The release command provides a variety of skip flags:

- `--skip-git`: Skips deriving version information from Git. Must only be used in conjunction with `--version`.
- `--skip-submit`: Skips submitting for review
- `--skip-update-metadata`: Skips updating metadata (app info, localizations, assets, review details, etc.)
- `--skip-update-pricing`: Skips updating app pricing

Finally, you can provide a `--timeout` flag with a duration value to place a limit on the runtime of the release process. While Cider is intended to be fast, apps with a lot of localization fields or many apps will take a while longer to run through all of their metadata.

### Environment

In addition to flags, Cider requires a few environment variables to be set in order to operate. They each correspond to an element of authorization described by [Creating API Keys for App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi/creating_api_keys_for_app_store_connect_api) from the Apple Developer Documentation.

- `ASC_KEY_ID` – Your key ID
- `ASC_ISSUER_ID` – Your team's issuer ID
- `ASC_PRIVATE_KEY_PATH` – A path to your .p8 private key

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
