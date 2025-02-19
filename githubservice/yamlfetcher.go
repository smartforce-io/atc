package githubservice

import (
	"gopkg.in/yaml.v2"
)

type Yaml struct {
	Version string `Yaml:"version"`
}

type YamlFetcher struct {
}

var unmarshalYaml = func(content []byte, yamlPtr *Yaml) error {
	return yaml.Unmarshal(content, yamlPtr)
}

func (yamlFetcher *YamlFetcher) GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error) {
	content, err := ghContentProvider.getContents(settings.Path)
	if err != nil {
		return "", err
	}
	y := &Yaml{}
	if err := unmarshalYaml([]byte(content), y); err != nil {
		return "", err
	}
	if y.Version == "" {
		return "", errNoVers
	}
	return y.Version, nil
}
