package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

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

	vers, err := fetcher.GetVersion(&cp, "pom.xml")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if vers != "5" {
		t.Errorf("wrong settings File! Got %q, wanted %q", vers, 5)
	}
}

func TestUnmarshalPomXml(t *testing.T) {
	var tests = []struct {
		content string
		version string
	}{
		{`<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
<version>4.2.0</version>
</project>`, `4.2.0`},
		{`<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
<version>4.2.0-release</version>
</project>`, `4.2.0-release`},
	}
	for _, test := range tests {
		pomxml := &PomXml{}
		err := unmarshalPomXml([]byte(test.content), pomxml)
		if err != nil {
			t.Errorf("Error unmarshal: %v", err)
			if pomxml.Version != test.version {
				t.Errorf("Unmarshal error for content: %s\n expected: %s, got: %s", test.content, test.version, pomxml.Version)
			}
		}
	}
}

func TestUnmarshalErrorPomXml(t *testing.T) {
	var tests = []struct {
		content string
		err     string
	}{
		{`<?xml version="1.0" encoding="UTF-8"?>`, `EOF`},
		{``, "EOF"},
	}
	for _, test := range tests {
		pomxml := &PomXml{}
		if err := unmarshalPomXml([]byte(test.content), pomxml); fmt.Sprintf("%s", err) != test.err {
			t.Errorf("Error for content: %s\nexpected err: %v, got err: %v", test.content, test.err, err)
		}
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

func TestErrorGetVersionPomXml(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	pxf := &pomXmlFetcher{}
	//test error get contents
	_, err := pxf.GetVersion(&cp, "Maven")
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = pxf.GetVersionDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error can't search verion
	cp.content = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
	</project>`
	cp.err = nil
	_, err = pxf.GetVersion(&cp, "Maven")
	if err != errNoVers {
		t.Errorf("err:%s  !=  noVersErr:%s", err, errNoVers)
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
