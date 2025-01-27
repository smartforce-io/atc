package pubspecyaml

import (
	"gopkg.in/yaml.v2"

	"github.com/smartforce-io/atc/githubservice/fetcher"
	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"
)

type PubspecYaml struct {
	Version string `PubspecYaml:"version"`
}

type Fetcher struct {
}

var unmarshalPubspecYaml = func(content []byte, pubspecyamlPtr *PubspecYaml) error {
	return yaml.Unmarshal(content, pubspecyamlPtr)
}

func (pubspecyamlFetcher *Fetcher) GetVersion(ghContentProvider provider.ContentProvider, settings settings.AtcSettings) (string, error) {
	content, err := ghContentProvider.GetContents(settings.Path)
	if err != nil {
		return "", err
	}
	pubspecyaml := &PubspecYaml{}
	if err := unmarshalPubspecYaml([]byte(content), pubspecyaml); err != nil {
		return "", err
	}
	if pubspecyaml.Version == "" {
		return "", fetcher.ErrNoVers
	}
	return pubspecyaml.Version, nil
}

func (pubspecyamlFetcher *Fetcher) GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error) {
	return pubspecyamlFetcher.GetVersion(ghContentProvider, settings.AtcSettings{Path: "pubspec.yaml"})
}
