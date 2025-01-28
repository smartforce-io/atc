package customregex

import (
	"errors"
	"regexp"

	"github.com/smartforce-io/atc/githubservice/provider"

	"github.com/smartforce-io/atc/githubservice/fetcher"
	"github.com/smartforce-io/atc/githubservice/settings"
)

type Config struct {
	Version string `customregexfetcher:"version"`
}

type Fetcher struct {
}

var unmarshalCustomRegexConfig = func(content []byte, regexStr string, customRegexConfigPtr *Config) error {
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		fetcher.ErrParsRegex = err
		return fetcher.ErrParsRegex
	}
	res := regex.FindStringSubmatch(string(content))
	if len(res) == 0 {
		return fetcher.ErrNoVers
	}
	if len(res) == 1 {
		return fetcher.ErrNoGroupInConf
	}
	customRegexConfigPtr.Version = res[1]
	return nil
}

func (customRegexFetcher *Fetcher) GetVersion(ghContentProvider provider.ContentProvider, settings settings.AtcSettings) (string, error) {
	content, err := ghContentProvider.GetContents(settings.Path)
	if err != nil {
		return "", err
	}
	customRegexConfig := &Config{}
	if err := unmarshalCustomRegexConfig([]byte(content), settings.RegexStr, customRegexConfig); err != nil {
		return "", err
	}
	return customRegexConfig.Version, nil
}

func (customRegexFetcher *Fetcher) GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error) {
	return "", errors.New("CustomRegexConfig doesn't have a default path")
}
