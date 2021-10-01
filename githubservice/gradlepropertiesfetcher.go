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

var unmarshalGrableProperties = func(content []byte, grablePropertiesPtr *GradleProperties) error {
	regex, _ := regexp.Compile("version=([^\t\n\f\r]*)")
	res := regex.FindStringSubmatch(string(content))
	if len(res) == 0 {
		return errors.New("grable.properties does not containt number version")
	}
	grablePropertiesPtr.Version = res[1]
	return nil
}

func (gradleBuildFetcher *gradlePropertiesFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, error) {
	content, err := ghContentProvider.getContents(path)
	if err != nil {
		return "", err
	}
	grable := &GradleProperties{}
	if err := unmarshalGrableProperties([]byte(content), grable); err != nil {
		return "", err
	}
	return grable.Version, nil
}

func (gradleBuildFetcher *gradlePropertiesFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return gradleBuildFetcher.GetVersion(ghContentProvider, "grable.properties")
}
