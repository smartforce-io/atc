package githubservice

import (
	"errors"
	"testing"
)

var failingPackageJsonFetcherUnmarshal = func(content []byte, packagejson *PackageJson) error {
	return errUnmarshal
}

func TestUnmarshalPackageJson(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{`{"version": "4.2.0"}`, `4.2.0`},
		{`{"version": "4.2.1",
		"name": "discord-musicbot"}`, `4.2.1`},
		{`{"name": "discord-musicbot",
		"version": "4.3.1"}`, `4.3.1`},
		{`{"version": "4.4.1-release"}`, "4.4.1-release"},
		{`{"version": 4.1.1}`, ``},
		{`{version: "4.1.1"}`, ``},
		{`{v1.1}`, ""},
		{``, ""},
	}
	for _, test := range tests {
		packagejson := &PackageJson{}
		if err := unmarshalPackageJson([]byte(test.input), packagejson); packagejson.Version != test.output {
			t.Errorf("err: %q, test.input= %s\nunmarshal vers= %s, expected vers= %s", err, test.input, packagejson.Version, test.output)
		}
	}
}

func TestErrorGetVersionPackageJson(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	pjf := &packagejsonFetcher{}
	//test error get contents
	_, err := pjf.GetVersion(&cp, "npm")
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = pjf.GetVersionDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
}

func TestPackageJsonFetcherUnmarshalError(t *testing.T) {
	fetcher := packagejsonFetcher{}

	unmarshalPackageJsonCopy := unmarshalPackageJson

	unmarshalPackageJson = failingPackageJsonFetcherUnmarshal

	cp := mockContentProvider{"", nil}

	_, err := fetcher.GetVersion(&cp, "package.json")

	if err != errUnmarshal {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshalPackageJson = unmarshalPackageJsonCopy
}
