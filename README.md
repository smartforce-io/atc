# Automated Tag Creator
ATC scans specified file that contains version (e.g., `plugin.yaml`, `package.json`, etc.) then if a new version is detected, it generates a corresponding Git tag and pushes it to the repository

By default ATC detects files such as `build.gradle`, `package.json`, `pom.xml`, `pubspec.yaml`, `plugin.yaml` but if you want to use another file for tag creation you can use `regex` key in your step.with configuration
## Usage

```
name: create tag
on:
  push:

permissions:
  contents: write

jobs:
  tests:
    name: Create Tag
    runs-on: ubuntu-24.04
    steps:
      - name: Check out code into the Go module directory
        uses:  uses: actions/checkout@v5

      - name: Create tag
        uses: smartforce-io/atc@master
        with:
          type: 'plugin.yaml'
          secrets: ${{ secrets.GITHUB_TOKEN }}
```

## Customizing

### inputs

Can be used as step.with keys:


| Name       | Type   | Default         | Description                                                                                                                                               |
|------------|--------|-----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `behavior` | String | `after`         | Determines which commit to read the version from.Use `after` to read version from the current commit,or `before` to read version from the previous commit |
| `template` | String | `v{{.Version}}` | Template for tag                                                                                                                                          |
| `regex`    | String | `version: (.+)` | Create regex string if you are not using the default ATC package manager. The regexstr must contain one group with version number.                        |
