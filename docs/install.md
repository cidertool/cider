---
layout: page
nav_order: 2
---

# Installation

Cider can be installed from a variety of sources.

### Homebrew (macOS/Linux)

```shell
brew install cidertool/tap/cider
```

### Scoop (Windows)

```shell
scoop bucket add cider https://github.com/cidertool/scoop-bucket.git
scoop install cider
```

### Docker

```shell
docker run --rm \
    --volume $PWD:/app \
    --workdir /app \
    --env ASC_KEY_ID \
    --env ASC_ISSUER_ID \
    --env ASC_PRIVATE_KEY \
    cidertool/cider release
```

### Pre-built Binary

Download the specific version for your platform on Cider's [releases page](https://github.com/cidertool/cider/releases).

### Compile from Source

```shell
git clone git@github.com:cidertool/cider.git
cd cider
go build -o cider .
```

From there, you can use your locally-built Cider binary for whatever purposes you need.

## Missing your favorite?

Please [file an issue](https://github.com/cidertool/cider/issues/new)!
