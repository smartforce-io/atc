# Config file ATC .atc.yaml

Config file uses syntax Yaml.\
.Atc.yaml settings are always taken from ** default ** branch! 

## Action Inputs
- [**Path**](#path): Path to package manager configuration file.
- [**Behavior**](#behavior): Commit to be used to create tag.
- [**Template**](#template): Tag template.
- [**Branch**](#branch): Branch for track version changes. 
- [**RegexStr**](#regexstr): Regex String to get version from custom configuration file.

## Examples:
* **Only path configuration.** When the configuration file *contents/pom.xml* changes the version project in *default* branch from 1.0.0 to 1.0.1, this example will create the tag "v1.0.1" in current commit.
```yaml
path: "contents/pom.xml"
```
* **Full configuration with default package manager.** When the configurate file *package.json* changes the version project in *release* branch from 1.1.0 to 1.1.1, this example will create a tag "v1.1.1-alfaNPM" in previous commit.
```yaml
path: "package.json"
behavior: "before"
template: "{{.Version}}-alfaNPM"
branch: "release"
```
* **Full configuration with custom package manager.** When the configurate file *test.txt* changes the version project in *release* branch from 1.2.0 to 1.2.1, this example will create a tag "v1.2.1-custom" in previous commit.
```yaml
path: "package.json"
behavior: "before"
template: "{{.Version}}-custom"
branch: "release"
regexStr: "vers: (.+)"
```
###### File test.txt:
```
name: test
vers: 1.2.1
project: test
```

## Сustomization ATC config file(.atc.yaml):
### Path
ATC supports: Gradle(build.gradle), NPM(package.json), Maven(pom.xml), Flutter(pubspec.yaml) or other config file if [RegexStr](#regexstr) is used. 
Path uses relative link to the package manager configuration file. Don't use "/" preffix.
```yaml
path: "contents/pom.xml"
path: "package.json"
path: "test.txt"
```
### Behavior
ATC can create tag for current commit, use **after** for this, or previous commit, use **before** for this. The default behavior is **after**.
###### Behavior examples:
```yaml
behavior: "after" # for use current commit
behavior: "before" # for use previous commit
```
### Template
ATC Template works with [GO Template](https://pkg.go.dev/text/template). Use "{{.Version}}" to write the number to the tag.
ATC supports the function time.Now(), to use it use {{Time}}. The default template is "v{{.Version}}".
###### Template examples:
```yaml
template: "v{{.Version}}" # for version = 2.0.0, tag = "v2.0.0"
template: "v{{.Version}}-alfa{{.Version}}" # for version = 2.0.1, tag = "v2.0.1-alfa2.0.1"
template: "{{.Version}}-{{Time.Hour}}" # for version = 2.0.2, tag = "v2.0.2-`Hours now`"
```
### Branch
ATC can track non-default branch. 
###### Branch examples:
```yaml
branch: "main"
branch: "issue-404"
branch: "testbranch"
```
### RegexStr
Write [Path](#path) to configuration file and create regex string if you are not using the default ATC package manager. 
The regexstr must contain one group with version number.
###### Regexstr examples:
```yaml
regexstr: "version: (.+)" # for `version: 2.0.0`
regexstr: "\"version\": \"(.+)\"" # for `"version": "2.0.1""`
```
