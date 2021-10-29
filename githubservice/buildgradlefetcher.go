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
	if len(res) < 3 {
		return errNoVers
	}
	buildGradlePtr.Version = res[2]
	return nil
}

func (buildGradleFetcher *buildGradleFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, error) {
	content, err := ghContentProvider.getContents(path)
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
	return buildGradleFetcher.GetVersion(ghContentProvider, "app/build.gradle")
}
