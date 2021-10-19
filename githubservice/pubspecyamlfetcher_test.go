package githubservice

import (
	"errors"
	"testing"
)

var failingPubspecYamlFetcherUnmarshal = func(content []byte, pubspecyaml *PubspecYaml) error {
	return errUnmarshal
}

func TestUnmarshalPubspecYaml(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{`version: 1.2.0`, `1.2.0`},
		{`version: 1.2.1
name: newtify`, `1.2.1`},
		{`name: newtify
version: 1.2.2`, `1.2.2`},
		{`version: 1.2.3-release`, `1.2.3-release`},
		{`"version": 1.2.4`, `1.2.4`},
		{`version: "1.2.5"`, `1.2.5`},
		{`v1.1`, ""},
		{``, ""},
	}
	for _, test := range tests {
		pubspecyaml := &PubspecYaml{}
		if err := unmarshalPubspecYaml([]byte(test.input), pubspecyaml); pubspecyaml.Version != test.output {
			t.Errorf("err: %q, test.input= %s\nunmarshal vers= %s, expected vers= %s", err, test.input, pubspecyaml.Version, test.output)
		}
	}
}

func TestErrorGetVersionPubspecYaml(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	psf := &pubspecyamlFetcher{}
	//test error get contents
	_, err := psf.GetVersion(&cp, "npm")
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = psf.GetVersionDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
}

func TestPubspecYamlFetcherUnmarshalError(t *testing.T) {
	fetcher := pubspecyamlFetcher{}

	unmarshalPubspecYamlCopy := unmarshalPubspecYaml

	unmarshalPubspecYaml = failingPubspecYamlFetcherUnmarshal

	cp := mockContentProvider{"", nil}

	_, err := fetcher.GetVersion(&cp, "pubspec.yaml")

	if err != errUnmarshal {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshalPubspecYaml = unmarshalPubspecYamlCopy
}
