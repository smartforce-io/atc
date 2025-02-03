package githubservice

type pubspecyamlFetcher struct {
	*YamlFetcher
}

func (pubspecyamlFetcher *pubspecyamlFetcher) GetVersionUsingDefaultPath(ghContentProvider contentProvider) (string, error) {
	return pubspecyamlFetcher.GetVersion(ghContentProvider, AtcSettings{Path: "pubspec.yaml"})
}
