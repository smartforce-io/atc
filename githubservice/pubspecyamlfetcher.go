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

func (pubspecyamlFetcher *pubspecyamlFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, error) {
	content, err := ghContentProvider.getContents(path)
	if err != nil {
		return "", err
	}
	pubspecyaml := &PubspecYaml{}
	if err := unmarshalPubspecYaml([]byte(content), pubspecyaml); err != nil {
		return "", err
	}
	return pubspecyaml.Version, nil
}

func (pubspecyamlFetcher *pubspecyamlFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return pubspecyamlFetcher.GetVersion(ghContentProvider, "pubspec.yaml")
}
