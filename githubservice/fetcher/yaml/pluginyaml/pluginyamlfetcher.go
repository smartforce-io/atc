package pluginyaml

import (
	"github.com/smartforce-io/atc/githubservice/fetcher/yaml"
	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"
)

type Fetcher struct {
	*yaml.Fetcher
}

func (pluginYamlFetcher *Fetcher) GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error) {
	return pluginYamlFetcher.GetVersion(ghContentProvider, settings.AtcSettings{Path: "plugin.yaml"})
}
