package githubservice

import (
	"encoding/xml"
)

type PomXml struct {
	Version string `xml:"version"`
}

type pomXmlFetcher struct {
}

var unmarshalPomXml = func(content []byte, pomXmlPtr *PomXml) error {
	return xml.Unmarshal([]byte(content), pomXmlPtr)
}

func (pomXmlFetcher *pomXmlFetcher) GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error) {
	content, err := ghContentProvider.getContents(settings.Path)
	if err != nil {
		return "", err
	}
	pom := &PomXml{}
	if err := unmarshalPomXml([]byte(content), pom); err != nil {
		return "", err
	}
	if pom.Version == "" {
		return "", errNoVers
	}
	return pom.Version, nil
}

func (pomXmlFetcher *pomXmlFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return pomXmlFetcher.GetVersion(ghContentProvider, AtcSettings{Path: "contents/pom.xml"})
}
