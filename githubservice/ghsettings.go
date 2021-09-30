package githubservice

import (
	"gopkg.in/yaml.v2"
)

var unmarshal = func(content []byte, atcSettingsPtr *AtcSettings) error {
	return yaml.Unmarshal([]byte(content), atcSettingsPtr)
}

type AtcSettings struct {
	//Type     string `json:"type"`
	Path     string `json:"path"`
	Behavior string `json:"behavior"`
	Template string `json:"template"`
	//Prefix string `json:"prefix"`
}

func getAtcSetting(ghcp contentProvider) (*AtcSettings, error) {
	content, err := ghcp.getContents(".atc.yaml")
	if err != nil {
		return nil, err
	}
	//GOTO check content != ""
	//.atc.yaml not found

	settings := &AtcSettings{}
	if err := unmarshal([]byte(content), settings); err != nil {
		return nil, err
	}
	//GOTO check template for contains ".version"
	return settings, nil
}
