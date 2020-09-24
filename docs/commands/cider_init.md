---
layout: page
parent: Commands
title: init
nav_order: 1
nav_exclude: false
---

## cider init

Generates a .cider.yml file

### Synopsis

Use to initialize a new Cider project. This will create a new configuration file
		in the current directory that should be checked into source control.

```
cider init [flags]
```

### Examples

```
cider init
```

### Options

```
  -f, --config string   Path of configuration file to create (default ".cider.yml")
  -h, --help            help for init
  -y, --skip-prompt     Skips onboarding prompts. This can result in an overwritten configuration file
```

### Options inherited from parent commands

```
      --debug   Enable debug mode
```

### SEE ALSO

* [cider](/commands/cider/)	 - Submit your builds to the Apple App Store in seconds

