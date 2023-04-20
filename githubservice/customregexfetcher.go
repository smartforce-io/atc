package githubservice

import (
	"errors"
	"regexp"
)

type CustomRegexConfig struct {
	Version string `customregexfetcher:"version"`
}

type customRegexFetcher struct {
}

var unmarshalCustomRegexConfig = func(content []byte, regexStr string, customRegexConfigPtr *CustomRegexConfig) error {
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		errParsRegex = err
		return errParsRegex
	}
	res := regex.FindStringSubmatch(string(content))
	if len(res) == 0 {
		return errNoVers
	}
	if len(res) == 1 {
		return errNoGroupInConf
	}
	customRegexConfigPtr.Version = res[1]
	return nil
}

func (customRegexFetcher *customRegexFetcher) GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error) {
	content, err := ghContentProvider.getContents(settings.Path)
	if err != nil {
		return "", err
	}
	customRegexConfig := &CustomRegexConfig{}
	if err := unmarshalCustomRegexConfig([]byte(content), settings.RegexStr, customRegexConfig); err != nil {
		return "", err
	}
	return customRegexConfig.Version, nil
}

func (customRegexFetcher *customRegexFetcher) GetVersionUsingDefaultPath(ghContentProvider contentProvider) (string, error) {
	return "", errors.New("CustomRegexConfig doesn't have a default path")
}
