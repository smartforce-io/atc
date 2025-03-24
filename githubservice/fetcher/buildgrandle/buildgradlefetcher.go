package buildgrandle

import (
	"log"
	"regexp"

	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"

	"github.com/smartforce-io/atc/githubservice/fetcher"
)

type BuildGradle struct {
	Version string `gradle:"version"`
}

type Fetcher struct {
}

var unmarshalBuildGradle = func(content []byte, buildGradlePtr *BuildGradle) error {
	regex, err := regexp.Compile(`defaultConfig {[^{}]*([^{}]*{[\s\S]*}[^{}]*)*[^{}]*\n[\t ]*versionName "(.+)"`)
	if err != nil {
		log.Printf("regexp compile err: %v", err)
		return err
	}
	res := regex.FindStringSubmatch(string(content))
	if len(res) < 2 {
		return fetcher.ErrNoVers
	}
	buildGradlePtr.Version = res[2]
	return nil
}

func (buildGradleFetcher *Fetcher) GetVersion(ghContentProvider provider.ContentProvider, settings settings.AtcSettings) (string, error) {
	content, err := ghContentProvider.GetContents(settings.Path)
	if err != nil {
		return "", err
	}
	gradle := &BuildGradle{}
	if err := unmarshalBuildGradle([]byte(content), gradle); err != nil {
		return "", err
	}
	return gradle.Version, nil
}

func (buildGradleFetcher *Fetcher) GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error) {
	return buildGradleFetcher.GetVersion(ghContentProvider, settings.AtcSettings{Path: "app/build.gradle"})
}
