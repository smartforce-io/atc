package githubservice

type PluginYamlFetcher struct {
	*YamlFetcher
}

func (pluginYamlFetcher *PluginYamlFetcher) GetVersionUsingDefaultPath(ghContentProvider contentProvider) (string, error) {
	return pluginYamlFetcher.GetVersion(ghContentProvider, AtcSettings{Path: "plugin.yaml"})
}
