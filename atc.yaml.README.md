# Config file ATC .atc.yaml

Config file use syntax Yaml. Settings .atc.yaml always takes from your **main** branch!

## Action Inputs
- [**Path**](#path): Path to your package manager config file.
- [**Behavior**](#behavior): What is commits will be used for create tag.
- [**Template**](#template): Template your tag.
- [**Branch**](#branch): Branch for traking changes versions.
- [**RegexStr**](#regexstr): Regex String for get version from custom config file.

## Examples:
* **Only path configurate.** When in config file *contents/pom.xml* version project will change in *main(master)* branch from 1.0.0 to 1.0.1, this example will create a tag "v1.0.1" in current commit.
```yaml
path: "contents/pom.xml"
```
* **full configure with default package manager.** When in config file *package.json* version project will change in *release* branch from 1.1.0 to 1.1.1, this example will create a tag "v1.1.1-alfaNPM" in previos commit.
```yaml
path: "package.json"
behavior: "before"
template: "{{.Version}}-alfaNPM"
branch: "release"
```
* **full configure with custom package manager.** When in config file *test.txt* version project will change in *release* branch from 1.2.0 to 1.2.1, this example will create a tag "v1.2.1-custom" in previos commit.
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
ATC supports: Gradle(build.gradle), NPM(package.json), Maven(pom.xml), Flutter(pubspec.yaml) or other config file with configurate [RegexStr](#regexstr). 
In field Path use relative link to your package manager config file. Don't use preffix "/".
```yaml
path: "contents/pom.xml"
path: "package.json"
path: "test.txt"
```
### Behavior
ATC can create tag for current commit, for this use **after**, or previos commit, then use **before**. Default behavior: **after**.
###### Behavior examples:
```yaml
behavior: "after" # for use current commit
behavior: "before" # for use previos commit
```
### Template
ATC template work with [GO Template](https://pkg.go.dev/text/template). For write number in tag use {{.Version}}.
ATC support function time.Now(), for use it {{Time}}. Default template: "v{{.Version}}".
###### Template examples:
```yaml
template: "v{{.Version}}" # for version = 2.0.0, tag = "v2.0.0"
template: "v{{.Version}}-alfa{{.Version}}" # for version = 2.0.1, tag = "v2.0.1-alfa2.0.1"
template: "{{.Version}}-{{Time.Hour}}" # for version = 2.0.2, tag = "v2.0.2-`Hours now`"
```
### Branch
Use this field, if u need track not main(master) branch in your project.
###### Branch examples:
```yaml
branch: "main"
branch: "issue-404"
branch: "testbranch"
```
### RegexStr
If you used no default ATC package manager, write config file to [Path](#path) and create regex string for this field.
Regex must contain one group with number version.
###### Regexstr examples:
```yaml
regexstr: "version: (.+)" # for `version: 2.0.0`
regexstr: "\"version\": \"(.+)\"" # for `"version": "2.0.1""`
```