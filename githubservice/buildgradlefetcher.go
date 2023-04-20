package githubservice

import (
	"log"
	"regexp"
)

type BuildGradle struct {
	Version string `gradle:"version"`
}

type buildGradleFetcher struct {
}

var unmarshalBuildGradle = func(content []byte, buildGradlePtr *BuildGradle) error {
	regex, err := regexp.Compile(`defaultConfig {[^{}]*([^{}]*{[\s\S]*}[^{}]*)*[^{}]*\n[\t ]*versionName "(.+)"`)
	if err != nil {
		log.Printf("regexp compile err: %v", err)
		return err
	}
	res := regex.FindStringSubmatch(string(content))
	if len(res) < 2 {
		return errNoVers
	}
	buildGradlePtr.Version = res[2]
	return nil
}

func (buildGradleFetcher *buildGradleFetcher) GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error) {
	content, err := ghContentProvider.getContents(settings.Path)
	if err != nil {
		return "", err
	}
	gradle := &BuildGradle{}
	if err := unmarshalBuildGradle([]byte(content), gradle); err != nil {
		return "", err
	}
	return gradle.Version, nil
}

func (buildGradleFetcher *buildGradleFetcher) GetVersionUsingDefaultPath(ghContentProvider contentProvider) (string, error) {
	return buildGradleFetcher.GetVersion(ghContentProvider, AtcSettings{Path: "app/build.gradle"})
}
