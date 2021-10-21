package githubservice

import (
	"regexp"
)

type BuildGradle struct {
	Version string `gradle:"version"`
}

type buildGradleFetcher struct {
}

var unmarshalBuildGradle = func(content []byte, buildGradlePtr *BuildGradle) error {
	regex, _ := regexp.Compile(`versionName "([^\t\n\f\r]+)"`)
	res := regex.FindStringSubmatch(string(content))
	if len(res) < 2 {
		return errNoVers
	}
	buildGradlePtr.Version = res[1]
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

func (buildGradleFetcher *buildGradleFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return buildGradleFetcher.GetVersion(ghContentProvider, AtcSettings{Path: "app/build.gradle"})
}
