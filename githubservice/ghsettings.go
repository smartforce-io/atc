package githubservice

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var unmarshal = func(content []byte, atcSettingsPtr *AtcSettings) error {
	return yaml.Unmarshal([]byte(content), atcSettingsPtr)
}

type AtcSettings struct {
	Path     string `json:"path"`
	Behavior string `json:"behavior"`
	Template string `json:"template"`
}

func checkSettingsForErrors(settings *AtcSettings) error {
	//check Path:
	pathPrefix := "/"
	pathSuffix := [3]string{"pom.xml", "gradle.properties", ".npmrc"}

	if settings.Path == "" {
		return errors.New(`error config file .atc.yaml; path = ""`)
	}
	if strings.HasPrefix(settings.Path, pathPrefix) {
		return errors.New(`error config file .atc.yaml; path has prefix "/"`)
	}
	if strings.Contains(settings.Path, "//") {
		return errors.New(`error config file .atc.yaml; path has "//"`)
	}
	sufOk := false
	for _, suf := range pathSuffix {
		if filepath.Base(settings.Path) == suf {
			sufOk = true
		}
	}
	if !sufOk {
		return errors.New(`error config file .atc.yaml: path no has suffix "pom.xml" or "gradle.properties" or ".npmrc"`)
	}
	//check Behavior:
	if settings.Behavior == "" {
		return errors.New(`error config file .atc.yaml; behavior = ""`)
	}
	if strings.ToLower(settings.Behavior) != "after" && strings.ToLower(settings.Behavior) != "before" {
		return errors.New(`error config file .atc.yaml: behavior no contains "before" or "after"`)
	}
	//check Template:
	if settings.Template == "" {
		return errors.New(`error config file .atc.yaml; template = ""`)
	}
	if !strings.Contains(settings.Template, `{{.Version}}`) {
		return errors.New(`error config file .atc.yaml: template no contains "{{.Version}}"`)
	}
	return nil
}

func getAtcSetting(ghcp contentProvider) (*AtcSettings, error) {
	settings := &AtcSettings{}

	content, err := ghcp.getContents(".atc.yaml")
	if err != nil {
		log.Printf("get .atc.yaml error: %s. Used default settings", err)
		return settings, nil
	}

	if err := unmarshal([]byte(content), settings); err != nil {
		return nil, errors.New(`error config file .atc.yaml; can't unmarshal file`)
	}

	if err := checkSettingsForErrors(settings); err != nil {
		return nil, err
	}

	return settings, nil
}
