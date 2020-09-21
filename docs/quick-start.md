---
layout: default
nav_order: 3
---

# Quick Start

Once you've [installed](../install) Cider, you can get started setting it up for your project.

Run `cider init` to create a new `.cider.yml` file in the current directory:

```shell
cider init
```

This will run through a series of prompts where you will get to set some default values for your project. See [configuration.md](./configuration.md) for additional options and documentation on the entire project specification. This file should be checked in to source control.

Once this file is set up, you can either proceed to run `cider` [locally](#local), or set it up in [CI](#ci).

## Local

The most simple invocation of Cider to submit an app is as follows:

```
cider release --mode appstore
```

Cider contains a host of options enabling you to customize its runtime. Follow the guide on the [`release` command](../commands/cider_release).

## CI

### GitHub Actions

Cider can also be run autonomously using the official [Cider Action](https://github.com/marketplace/actions/cider-action) hosted on the GitHub Marketplace. The Action is versioned independently of Cider, and all of Cider's commands and internal capabilities are available.

#### Usage

```yaml
- uses: actions/checkout@v2
- uses: cidertool/cider-action@v0
  with:
    version: latest
    args: release --mode appstore --set-version ${{ env.VERSION }}
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    ASC_KEY_ID: ${{ secrets.ASC_KEY_ID }}
    ASC_ISSUER_ID: ${{ secrets.ASC_ISSUER_ID }}
    ASC_PRIVATE_KEY: ${{ secrets.ASC_PRIVATE_KEY }}
```

You can run this job in any context you see fit to use Cider to update app metadata or submit new versions of your apps to the App Store.

### Buildkite

If you're using Buildkite, you can use the [Cider Buildkite Plugin](https://github.com/cidertool/cider-buildkite-plugin). Similarly to the GitHub Action, the plugin is versioned independently of Cider and any function available in the Cider command line can be used. This plugin requires Docker.

#### Usage

```yaml
steps:
  - label: ":apple: Release with Cider"
    plugins:
      - cidertool/cider#v0.1.0:
          args: release --mode appstore
    env:
      ASC_KEY_ID: "..."
      ASC_ISSUER_ID: "..."
      ASC_PRIVATE_KEY: "..."
```
