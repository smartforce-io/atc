package githubservice

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	behaviorBefore = "before"
	behaviorAfter  = "after"
	pathPrefix     = "/"
)

var unmarshal = func(content []byte, atcSettingsPtr *AtcSettings) error {
	return yaml.Unmarshal([]byte(content), atcSettingsPtr)
}

type AtcSettings struct {
	Path     string `json:"path"`
	Behavior string `json:"behavior"`
	Template string `json:"template"`
	Branch   string `json:"branch"`
}

func validateSettings(settings *AtcSettings) error {
	//check settins to "" and use default value:
	if settings.Behavior == "" {
		settings.Behavior = behaviorAfter
	}
	if settings.Template == "" {
		settings.Template = "v{{.Version}}"
	}

	//check Behavior:
	if strings.ToLower(settings.Behavior) != behaviorAfter && strings.ToLower(settings.Behavior) != behaviorBefore {
		return errors.New(`error config file .atc.yaml: behavior doesn't contain "before" or "after"`)
	}
	//check Template:
	if !strings.Contains(settings.Template, `{{.Version}}`) {
		return errors.New(`error config file .atc.yaml: template doesn't contain "{{.Version}}"`)
	}
	//check Path:
	if settings.Path == "" {
		return nil
	}
	if strings.HasPrefix(settings.Path, pathPrefix) {
		return errors.New(`error config file .atc.yaml; path has prefix "/"`)
	}
	if strings.Contains(settings.Path, "//") {
		return errors.New(`error config file .atc.yaml; path has "//"`)
	}
	sufOk := false
	for fether := range autoFetchers {
		if filepath.Base(settings.Path) == fether {
			sufOk = true
		}
	}
	if !sufOk {
		return errors.New(`error config file .atc.yaml: path no has suffix "pom.xml", "build.gradle", "package.json" or "pubspec.yaml"`)
	}
	return nil
}

func getAtcSetting(ghcp contentProvider) (*AtcSettings, error) {
	settings := &AtcSettings{}

	content, err := ghcp.getContents(".atc.yaml")
	if err != nil {
		log.Printf("get .atc.yaml error: %s. Used default settings", err)
		return &AtcSettings{Behavior: "after", Template: "v{{.Version}}"}, nil
	}

	if err := unmarshal([]byte(content), settings); err != nil {
		return nil, errors.New(`error config file .atc.yaml; can't unmarshal file`)
	}

	if err := validateSettings(settings); err != nil {
		return nil, err
	}
	return settings, nil
}
