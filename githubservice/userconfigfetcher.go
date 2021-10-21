package githubservice

import (
	"errors"
	"regexp"
)

type UserConfig struct {
	Version string `gradle:"version"`
}

type userConfigFetcher struct {
}

var unmarshalUserConfig = func(content []byte, regexStr string, userConfigPtr *UserConfig) error {
	regex, _ := regexp.Compile(regexStr)
	res := regex.FindStringSubmatch(string(content))
	if len(res) < 2 {
		return errNoVers
	}
	userConfigPtr.Version = res[1]
	return nil
}

func (userConfigFetcher *userConfigFetcher) GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error) {
	content, err := ghContentProvider.getContents(settings.Path)
	if err != nil {
		return "", err
	}
	userConfig := &UserConfig{}
	if err := unmarshalUserConfig([]byte(content), settings.RegexStr, userConfig); err != nil {
		return "", err
	}
	return userConfig.Version, nil
}

func (userConfigFetcher *userConfigFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return "", errors.New("UserConfig don't have default path")
}
