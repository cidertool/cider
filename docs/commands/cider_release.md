---
layout: page
parent: Commands
title: release
nav_order: 2
nav_exclude: false
---

## cider release

Release the selected apps in the current project

### Synopsis

Release the selected apps in the current project.
		
You can provide a path to a project directory as an argument to be the root directory
of all relative path expansions in the program, such as the Git repository, preview sets,
and screenshot resources. An exception to this is if you set a custom configuration file
path with the `--config` flag.

Additionally, Cider requires a few environment variables to be set in order to operate.
They each correspond to an element of authorization described by the Apple Developer Documentation.

- `ASC_KEY_ID`: The key's ID.
- `ASC_ISSUER_ID`: Your team's issuer ID.
- `ASC_PRIVATE_KEY` or `ASC_PRIVATE_KEY_PATH`: The .p8 private key issued by Apple.

These three values each have varying degrees of sensetivity and should be treated as secrets. Store
them securely in your environment so Cider can leverage them safely.

More info: https://developer.apple.com/documentation/appstoreconnectapi/creating_api_keys_for_app_store_connect_api

```
cider release [path] [flags]
```

### Examples

```
cider release --mode=appstore --set-version="1.0"
```

### Options

```
  -A, --all-apps --app                Process all apps in the configuration file. Supercedes any usage of the --app flag.
  -a, --app stringArray               Process the given app, providing the app key name used in your configuration file.
                                      
                                      This flag can be provided repeatedly for each app you want to process. You can omit
                                      this flag if your configuration file has only one app defined.
  -f, --config string                 Load configuration from file
  -h, --help                          help for release
  -p, --max-processes int             Run certain metadata syncing and asset uploading logic in parallel with
                                      the maximum allowable concurrency. (default 1)
      --mode {appstore,testflight}    Mode used to declare the publishing target for submission.
                                      		
                                      The default is "testflight" for submitting to Testflight, and the other alternative
                                      option is "appstore" for submitting to the App Store.
      --set-beta-group stringArray    Provide names of beta groups to release to instead of using
                                      the configuration file.
      --set-beta-tester stringArray   Provide email addresses of beta testers to release to instead of
                                      using the configuration file.
  -B, --set-build string              Build override to use instead of "latest". Corresponds to the CFBundleVersion
                                      of your build.
                                      		
                                      The default behavior without this flag is to select the latest build. In both cases,
                                      if the selected build has an invalid processing state, Cider will abort with an error
                                      to ensure your release is handled safely.
  -V, --set-version string            Version string override to use instead of parsing Git tags. Corresponds to the
                                      CFBundleShortVersionString of your build.
                                      
                                      Cider expects this string to follow the Major.Minor.Patch semantics outlined in Apple documentation
                                      and Semantic Versioning (semver). If this flag is omitted, Git will be leveraged to determine the
                                      latest tag. The tag will be used to calculate the version string under the same constraints.
      --skip-git --set-version        Skips deriving version information from Git. Must only be used in conjunction with the --set-version flag.
      --skip-submit                   Skips submitting for review
      --skip-update-metadata          Skips updating metadata (app info, localizations, assets, review details, etc.)
      --skip-update-pricing           Skips updating app pricing
      --timeout duration              Timeout for the entire release process.
                                      		
                                      If the command takes longer than this amount of time to run, Cider will abort. (default 30m0s)
```

### Options inherited from parent commands

```
      --debug   Enable debug mode
```

### SEE ALSO

* [cider](/commands/cider/)	 - Submit your builds to the Apple App Store in seconds

