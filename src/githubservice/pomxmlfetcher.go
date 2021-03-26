package githubservice

import "encoding/xml"

type PoxXml struct {
	Version string `xml:"version"`
}

type pomXmlFetcher struct {
}

func (pomXmlFetcher pomXmlFetcher) GetVersion(ghContentProvider *ghContentProvider, path string) (string, *RequestError, error) {
	content, reqErr, err := ghContentProvider.getContents(path)
	if err != nil || reqErr != nil {
		return "", reqErr, err
	}
	pom := &PoxXml{}
	if err := xml.Unmarshal([]byte(content), pom); err != nil {
		return "", nil, err
	}
	return pom.Version, nil, nil

}
