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
	if len(res) == 0 {
		return errNoVers
	}
	buildGradlePtr.Version = res[1]
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

func (buildGradleFetcher *buildGradleFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return buildGradleFetcher.GetVersion(ghContentProvider, "app/build.gradle")
}
