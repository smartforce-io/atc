package pomxml

import (
	"encoding/xml"

	"github.com/smartforce-io/atc/githubservice/provider"

	"github.com/smartforce-io/atc/githubservice/fetcher"
	"github.com/smartforce-io/atc/githubservice/settings"
)

type PomXml struct {
	Version string `xml:"version"`
}

type Fetcher struct {
}

var unmarshalPomXml = func(content []byte, pomXmlPtr *PomXml) error {
	return xml.Unmarshal([]byte(content), pomXmlPtr)
}

func (pomXmlFetcher *Fetcher) GetVersion(ghContentProvider provider.ContentProvider, settings settings.AtcSettings) (string, error) {
	content, err := ghContentProvider.GetContents(settings.Path)
	if err != nil {
		return "", err
	}
	pom := &PomXml{}
	if err := unmarshalPomXml([]byte(content), pom); err != nil {
		return "", err
	}
	if pom.Version == "" {
		return "", fetcher.ErrNoVers
	}
	return pom.Version, nil
}

func (pomXmlFetcher *Fetcher) GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error) {
	return pomXmlFetcher.GetVersion(ghContentProvider, settings.AtcSettings{Path: "pom.xml"})
}
