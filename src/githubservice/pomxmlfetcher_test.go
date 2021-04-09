package githubservice

import "testing"

var basicPomXml = `
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		 xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
		 <version>5</version>
</project>
`
var failingPomXmlFetcherUnmarshal = func(content []byte, pomXmlPtr *PomXml) error {
	return errUnmarshal
}

func TestPomXmlFetcherBasic(t *testing.T) {
	fetcher := pomXmlFetcher{}

	cp := mockContentProvider{basicPomXml, nil}

	v, err := fetcher.GetVersion(&cp, "pom.xml")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if v != "5" {
		t.Errorf("wrong settings File! Got %q, wanted %q", v, 5)
	}
}
func TestPomXmlFetcherGeneralError(t *testing.T) {
	fetcher := pomXmlFetcher{}

	cp := mockContentProvider{content: "", err: errGeneral}

	_, err := fetcher.GetVersion(&cp, "pom.xml")

	if err != errGeneral {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errGeneral)
	}
}
func TestPomXmlFetcherUnmarshalError(t *testing.T) {
	fetcher := pomXmlFetcher{}

	unmarshalPomXmlCopy := unmarshalPomXml

	unmarshalPomXml = failingPomXmlFetcherUnmarshal

	cp := mockContentProvider{basicPomXml, nil}

	_, err := fetcher.GetVersion(&cp, "pom.xml")

	if err != errUnmarshal {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshalPomXml = unmarshalPomXmlCopy
}
