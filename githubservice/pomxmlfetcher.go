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

func (pomXmlFetcher *pomXmlFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, error) {
	content, err := ghContentProvider.getContents(path)
	if err != nil {
		return "", err
	}
	pom := &PomXml{}
	if err := unmarshalPomXml([]byte(content), pom); err != nil {
		return "", err
	}
	return pom.Version, nil

}

func (pomXmlFetcher *pomXmlFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return pomXmlFetcher.GetVersion(ghContentProvider, "pom.xml")
}
