<p align="center">
  <!-- <img alt="applereleaser logo" src="assets/go.png" height="150" /> -->
  <h3 align="center">applereleaser</h3>
  <p align="center">Submit to the App Store in seconds!</p>
</p>

---

`applereleaser` is a tool managing the entire release process of an iOS, macOS or tvOS application, supported by official Apple APIs. It takes the builds you've uploaded to App Store Connect, updates their metadata, and submits them for review automatically using an expressive YAML configuration. Unlike Xcode or altool, `applereleaser` is designed to be useful on Linux and Windows, in addition to macOS. 

## Install `applereleaser`

- [On my machine](https://aaronsky.github.io/applereleaser/usage/install/#local);
- [On CI/CD systems](https://aaronsky.github.io/applereleaser/usage/install/#ci).

## Documentation

Documentation is hosted at <https://aaronsky.github.io/applereleaser>.

## Badges

![build](https://github.com/aaronsky/applereleaser/workflows/build/badge.svg)
[![License](https://img.shields.io/github/license/aaronsky/applereleaser)](/LICENSE)
[![Release](https://img.shields.io/github/release/aaronsky/applereleaser.svg)](https://github.com/aaronsky/applereleaser/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/aaronsky/applereleaser)](https://hub.docker.com/r/aaronsky/applereleaser)
[![Github Releases Stats of applereleaser](https://img.shields.io/github/downloads/aaronsky/applereleaser/total.svg?logo=github)](https://somsubhra.com/github-release-stats/?username=aaronsky&repository=applereleaser)

## Contributing

This project's primary goal is to simplify the process to release on the App Store, and enable the entire build + test + release process to be executable in the command line. Until the package's version stabilizes with v1, there isn't a strong roadmap beyond those stated goals. However, contributions are always welcome. If you want to get involved or you just want to offer feedback, please see [`CONTRIBUTING.md`](./.github/CONTRIBUTING.md) for details.

## License

This library is licensed under the [MIT License](./LICENSE)
