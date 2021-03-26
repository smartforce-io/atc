package githubservice

import (
	"errors"
	"log"

	"gopkg.in/yaml.v2"
)

type AtcSettings struct {
	Path   string `json:"path"`
	Prefix string `json:"prefix"`
}

var (
	errFailedResponse = errors.New("failed response")
)

func getAtcSetting(ghcp *ghContentProvider) (*AtcSettings, error) {
	content, reqErr, err := ghcp.getContents(".atc.yaml")

	if err != nil {
		return nil, err
	}

	if reqErr != nil {
		log.Printf("getAtcSetting received a failed status for %s/%s: %q", ghcp.owner, ghcp.repo, reqErr.StatusCode)
		return nil, errFailedResponse
	}

	settings := &AtcSettings{}
	if err := yaml.Unmarshal([]byte(content), settings); err != nil {
		return nil, err
	}
	return settings, nil
}
