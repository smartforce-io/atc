package githubservice

import (
	"gopkg.in/yaml.v2"
)

var unmarshal = func(content []byte, atcSettingsPtr *AtcSettings) error {
	return yaml.Unmarshal([]byte(content), atcSettingsPtr)
}

type AtcSettings struct {
	Path   string `json:"path"`
	Prefix string `json:"prefix"`
}

func getAtcSetting(ghcp contentProvider) (*AtcSettings, error) {
	content, err := ghcp.getContents(".atc.yaml")

	if err != nil {
		return nil, err
	}

	settings := &AtcSettings{}
	if err := unmarshal([]byte(content), settings); err != nil {
		return nil, err
	}
	return settings, nil
}
