---
layout: page
title: FAQ
nav_order: 6
---

# Frequently Asked Questions

## Can I use `applereleaser` to upload builds?

`applereleaser` is not designed to replace `altool`, the recommended Apple-provided tool for validating and uploading archives to Xcode. That tool has an exceptionally well-documented manpage and comes with stability guarantees from the platform-holder. Please use applereleaser to submit your builds once they've been uploaded and completed processing.

## Can I use an Apple ID to authenticate?

No. `applereleaser` is built on the backbone of Apple's official App Store Connect API, and a valid [JSON Web Token](https://tools.ietf.org/html/rfc7519) generated from an issuer ID, key ID and private key are essential components in generating a compliant token. Create a key for a user on your team [following Apple's documentation](https://developer.apple.com/documentation/appstoreconnectapi/creating_api_keys_for_app_store_connect_api) and use those credentials in your CI environment. If you ever need to revoke these credentials, you can do so by following [these instructions](https://developer.apple.com/documentation/appstoreconnectapi/revoking_api_keys).

## How is this different from Fastlane/Spaceship

Spaceship, and by extension Fastlane, are designed to be customizable for a variety of features and functions. You can do largely do anything, but that comes with the inherent overhead that "anything" entails. Spaceship has served Fastlane and the broader Apple development community well for years, but the investment cost can't be denied. Additionally, Spaceship was originally designed around Apple's private iTunes Connect API, and its migration to the official App Store Connect API has been slow. `applereleaser` has been designed with simplicity and portability in mind, which has required limiting its scope from "anything". In addition, `applereleaser` has been built around the App Store Connect API from the very beginning. What you get is a tool that is useful out-of-the-box, with simple configuration options, that runs quickly anywhere. 
