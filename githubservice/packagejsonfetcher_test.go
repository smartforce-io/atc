package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

var basicPackageJson = `
{"name": "my_package",
"description": "",
"version": "1.5.3"}
`

var failingPackageJsonFetcherUnmarshal = func(content []byte, packagejson *PackageJson) error {
	return errUnmarshal
}

func TestPackageJsonFetcherBasic(t *testing.T) {
	fetcher := packagejsonFetcher{}

	cp := mockContentProvider{basicPackageJson, nil}

	vers, err := fetcher.GetVersion(&cp, AtcSettings{Path: "package.json"})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if vers != "1.5.3" {
		t.Errorf("wrong settings File! Got %q, wanted %q", vers, 5)
	}
}

func TestUnmarshalPackageJson(t *testing.T) {
	var tests = []struct {
		content string
		version string
	}{
		{`{"version": "4.2.0"}`, `4.2.0`},
		{`{"version": "4.2.1",
		"name": "discord-musicbot"}`, `4.2.1`},
		{`{"name": "discord-musicbot",
		"version": "4.3.1"}`, `4.3.1`},
		{`{"version": "4.4.1-release"}`, "4.4.1-release"},
	}
	for _, test := range tests {
		packagejson := &PackageJson{}
		err := unmarshalPackageJson([]byte(test.content), packagejson)
		if err != nil {
			t.Errorf("Error unmarshal: %v", err)
			if packagejson.Version != test.version {
				t.Errorf("Unmarshal error for content: %s\n expected: %s, got: %s", test.content, test.version, packagejson.Version)
			}
		}
	}
}

func TestUnmarshalErrorPackageJson(t *testing.T) {
	var tests = []struct {
		content string
		err     string
	}{
		{`{"version": 4.1.1}`, `invalid character '.' after object key:value pair`},
		{`{version: "4.1.1"}`, `invalid character 'v' looking for beginning of object key string`},
		{`{v1.1}`, "invalid character 'v' looking for beginning of object key string"},
		{``, "unexpected end of JSON input"},
	}
	for _, test := range tests {
		packagejson := &PackageJson{}
		if err := unmarshalPackageJson([]byte(test.content), packagejson); fmt.Sprintf("%s", err) != test.err {
			t.Errorf("Error for content: %s\nexpected err: %v, got err: %v", test.content, test.err, err)
		}
	}
}

func TestErrorGetVersionPackageJson(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	pjf := &packagejsonFetcher{}
	//test error get contents
	_, err := pjf.GetVersion(&cp, AtcSettings{Path: "npm"})
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = pjf.GetVersionUsingDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error can't search verion
	cp.content = `{"name": "atc"}`
	cp.err = nil
	_, err = pjf.GetVersion(&cp, AtcSettings{Path: "NPM"})
	if err != errNoVers {
		t.Errorf("err:%s  !=  noVersErr:%s", err, errNoVers)
	}
}

func TestPackageJsonFetcherUnmarshalError(t *testing.T) {
	fetcher := packagejsonFetcher{}

	unmarshalPackageJsonCopy := unmarshalPackageJson

	unmarshalPackageJson = failingPackageJsonFetcherUnmarshal

	cp := mockContentProvider{"", nil}

	_, err := fetcher.GetVersion(&cp, AtcSettings{Path: "package.json"})

	if err != errUnmarshal {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshalPackageJson = unmarshalPackageJsonCopy
}
