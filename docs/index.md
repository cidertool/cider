---
layout: home
title: Home
nav_order: 1
description: "Cider is a command-line application that makes it easy to submit your Apple App Store apps for review."
permalink: /
---

# Cider

Cider is a tool managing the entire release process of an iOS, macOS or tvOS application, supported by official Apple APIs. It takes the builds you've uploaded to App Store Connect, updates their metadata, and submits them for review automatically using an expressive YAML configuration. Unlike Xcode or altool, Cider is designed to be useful on Linux and Windows, in addition to macOS. 

Cider is not a replacement for `altool`, the official command-line interface for uploading, validating, and notarizing archives to Apple. It's instead designed to complement `altool`, and by extension `xcodebuild`. With Cider, your pipeline can build, test, upload, and now release your app without any required manual action.
