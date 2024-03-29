package githubservice

import "gopkg.in/yaml.v2"

type PubspecYaml struct {
	Version string `PubspecYaml:"version"`
}

type pubspecyamlFetcher struct {
}

var unmarshalPubspecYaml = func(content []byte, pubspecyamlPtr *PubspecYaml) error {
	return yaml.Unmarshal(content, pubspecyamlPtr)
}

func (pubspecyamlFetcher *pubspecyamlFetcher) GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error) {
	content, err := ghContentProvider.getContents(settings.Path)
	if err != nil {
		return "", err
	}
	pubspecyaml := &PubspecYaml{}
	if err := unmarshalPubspecYaml([]byte(content), pubspecyaml); err != nil {
		return "", err
	}
	if pubspecyaml.Version == "" {
		return "", errNoVers
	}
	return pubspecyaml.Version, nil
}

func (pubspecyamlFetcher *pubspecyamlFetcher) GetVersionUsingDefaultPath(ghContentProvider contentProvider) (string, error) {
	return pubspecyamlFetcher.GetVersion(ghContentProvider, AtcSettings{Path: "pubspec.yaml"})
}
