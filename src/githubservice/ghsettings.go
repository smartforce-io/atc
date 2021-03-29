package githubservice

import (
	"errors"

	"gopkg.in/yaml.v2"
)

var unmarshal = func(content []byte, atcSettingsPtr *AtcSettings) error {
	return yaml.Unmarshal([]byte(content), atcSettingsPtr)
}

type AtcSettings struct {
	Path   string `json:"path"`
	Prefix string `json:"prefix"`
}

var (
	errFailedResponse = errors.New("failed response")
)

func getAtcSetting(ghcp contentProvider) (*AtcSettings, error) {
	content, reqErr, err := ghcp.getContents(".atc.yaml")

	if err != nil {
		return nil, err
	}

	if reqErr != nil {
		return nil, errFailedResponse
	}

	settings := &AtcSettings{}
	if err := unmarshal([]byte(content), settings); err != nil {
		return nil, err
	}
	return settings, nil
}
