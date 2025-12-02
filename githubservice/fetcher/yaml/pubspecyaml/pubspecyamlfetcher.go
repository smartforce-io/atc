package pubspecyaml

import (
	"github.com/smartforce-io/atc/githubservice/fetcher/yaml"
	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"
)

type Fetcher struct {
	*yaml.Fetcher
}

func (pubspecyamlFetcher *Fetcher) GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error) {
	return pubspecyamlFetcher.GetVersion(ghContentProvider, settings.AtcSettings{Path: "pubspec.yaml"})
}
