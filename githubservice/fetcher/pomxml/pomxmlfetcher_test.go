package pomxml

import (
	"errors"
	"fmt"

	"github.com/smartforce-io/atc/githubservice/provider"

	"testing"

	"github.com/smartforce-io/atc/githubservice/fetcher"
	"github.com/smartforce-io/atc/githubservice/settings"
)

var basicPomXml = `
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		 xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
		 <version>5</version>
</project>
`
var failingPomXmlFetcherUnmarshal = func(content []byte, pomXmlPtr *PomXml) error {
	return provider.ErrUnmarshal
}

func TestPomXmlFetcherBasic(t *testing.T) {
	f := Fetcher{}

	cp := provider.MockContentProvider{Content: basicPomXml}

	vers, err := f.GetVersion(&cp, settings.AtcSettings{Path: "pom.xml"})

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
	f := Fetcher{}

	cp := provider.MockContentProvider{Content: "", Err: provider.ErrGeneral}

	_, err := f.GetVersion(&cp, settings.AtcSettings{Path: "pom.xml"})

	if !errors.Is(err, provider.ErrGeneral) {
		t.Errorf("Invalid error, Got %v, wanted %v", err, provider.ErrGeneral)
	}
}

func TestErrorGetVersionPomXml(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := provider.MockContentProvider{Err: noContentErr}
	pxf := &Fetcher{}
	//test error get contents
	_, err := pxf.GetVersion(&cp, settings.AtcSettings{Path: "Maven"})
	if !errors.Is(err, noContentErr) {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = pxf.GetVersionUsingDefaultPath(&cp)
	if !errors.Is(err, noContentErr) {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error can't search verion
	cp.Content = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
	</project>`
	cp.Err = nil
	_, err = pxf.GetVersion(&cp, settings.AtcSettings{Path: "Maven"})
	if !errors.Is(err, fetcher.ErrNoVers) {
		t.Errorf("err:%s  !=  noVersErr:%s", err, fetcher.ErrNoVers)
	}
}

func TestPomXmlFetcherUnmarshalError(t *testing.T) {
	f := Fetcher{}

	unmarshalPomXmlCopy := unmarshalPomXml

	unmarshalPomXml = failingPomXmlFetcherUnmarshal

	cp := provider.MockContentProvider{Content: basicPomXml}

	_, err := f.GetVersion(&cp, settings.AtcSettings{Path: "pom.xml"})

	if !errors.Is(err, provider.ErrUnmarshal) {
		t.Errorf("Invalid error, Got %v, wanted %v", err, provider.ErrUnmarshal)
	}

	unmarshalPomXml = unmarshalPomXmlCopy
}
