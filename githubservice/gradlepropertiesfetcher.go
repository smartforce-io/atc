package githubservice

import (
	"errors"
	"regexp"
)

type GradleProperties struct {
	Version string `gradle:"version"`
}

type gradlePropertiesFetcher struct {
}

var unmarshalGradleProperties = func(content []byte, gradlePropertiesPtr *GradleProperties) error {
	regex, _ := regexp.Compile("version=([^\t\n\f\r]*)")
	res := regex.FindStringSubmatch(string(content))
	if len(res) == 0 {
		return errors.New("gradle.properties does not containt number version")
	}
	gradlePropertiesPtr.Version = res[1]
	return nil
}

func (gradleBuildFetcher *gradlePropertiesFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, error) {
	content, err := ghContentProvider.getContents(path)
	if err != nil {
		return "", err
	}
	gradle := &GradleProperties{}
	if err := unmarshalGradleProperties([]byte(content), gradle); err != nil {
		return "", err
	}
	return gradle.Version, nil
}

func (gradleBuildFetcher *gradlePropertiesFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return gradleBuildFetcher.GetVersion(ghContentProvider, "gradle.properties")
}
