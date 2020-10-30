<p align="center">
  <img 
    alt="Cider logo" 
    src="docs/assets/images/header.png" 
    height="100px"
  />
  <p align="center">Submit to the App Store in seconds!</p>
</p>

---

Cider is a tool managing the entire release process of an iOS, macOS or tvOS application, supported by official Apple APIs. It takes the builds you've uploaded to App Store Connect, updates their metadata, and submits them for review automatically using an expressive YAML configuration. Unlike Xcode or altool, Cider is designed to be useful on Linux and Windows, in addition to macOS.

## Documentation

Documentation is hosted at <https://cidertool.github.io/cider>. Check out our installation and quick start documentation!

## Integrations

-   [GitHub Action](https://github.com/marketplace/actions/cider-action)
-   [Buildkite Plugin](https://github.com/cidertool/cider-buildkite-plugin)

## Badges

![build](https://github.com/cidertool/cider/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/cidertool/cider/branch/main/graph/badge.svg)](https://codecov.io/gh/cidertool/cider)
[![License](https://img.shields.io/github/license/cidertool/cider)](/COPYING)
[![Release](https://img.shields.io/github/release/cidertool/cider.svg)](https://github.com/cidertool/cider/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/cidertool/cider)](https://hub.docker.com/r/cidertool/cider)
[![Github Releases Stats of Cider](https://img.shields.io/github/downloads/cidertool/cider/total.svg?logo=github)](https://somsubhra.com/github-release-stats/?username=cidertool&repository=cider)

## Contributing

This project's primary goal is to simplify the process to release on the App Store, and enable the entire build + test + release process to be executable in the command line. Until the package's version stabilizes with v1, there isn't a strong roadmap beyond those stated goals. However, contributions are always welcome. If you want to get involved or you just want to offer feedback, please see [`CONTRIBUTING.md`](./.github/CONTRIBUTING.md) for details.

## Credits

Special thanks to:

-   [GoReleaser](https://goreleaser.com/) for inspiring the architecture and open sourcing several components used in Cider

## License

This library is licensed under the GNU General Public License v3.0 or later

See [COPYING](./COPYING) to see the full text.
