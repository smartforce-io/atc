package yaml

import (
	"github.com/smartforce-io/atc/githubservice/fetcher"
	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"
	"gopkg.in/yaml.v2"
)

type Yaml struct {
	Version string `Yaml:"version"`
}

type Fetcher struct {
}

var UnmarshalYaml = func(content []byte, yamlPtr *Yaml) error {
	return yaml.Unmarshal(content, yamlPtr)
}

func (yamlFetcher *Fetcher) GetVersion(ghContentProvider provider.ContentProvider, settings settings.AtcSettings) (string, error) {
	content, err := ghContentProvider.GetContents(settings.Path)
	if err != nil {
		return "", err
	}
	y := &Yaml{}
	if err := UnmarshalYaml([]byte(content), y); err != nil {
		return "", err
	}
	if y.Version == "" {
		return "", fetcher.ErrNoVers
	}
	return y.Version, nil
}
