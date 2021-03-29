package githubservice

import "encoding/xml"

type PomXml struct {
	Version string `xml:"version"`
}

type pomXmlFetcher struct {
}

var unmarshalPomXml = func(content []byte, pomXmlPtr *PomXml) error {
	return xml.Unmarshal([]byte(content), pomXmlPtr)
}

func (pomXmlFetcher pomXmlFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, *RequestError, error) {
	content, reqErr, err := ghContentProvider.getContents(path)
	if err != nil || reqErr != nil {
		return "", reqErr, err
	}
	pom := &PomXml{}
	if err := unmarshalPomXml([]byte(content), pom); err != nil {
		return "", nil, err
	}
	return pom.Version, nil, nil

}
