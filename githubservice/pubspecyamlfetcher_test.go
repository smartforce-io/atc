package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

var basicPubspecYaml = `
name: newtify
"version": 1.2.5
description: >-
`

var failingPubspecYamlFetcherUnmarshal = func(content []byte, pubspecyaml *PubspecYaml) error {
	return errUnmarshal
}

func TestPubspecYamlFetcherBasic(t *testing.T) {
	fetcher := pubspecyamlFetcher{}

	cp := mockContentProvider{basicPubspecYaml, nil}

	vers, err := fetcher.GetVersion(&cp, AtcSettings{Path: "pubspec.yaml"})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if vers != "1.2.5" {
		t.Errorf("wrong settings File! Got %q, wanted %q", vers, 5)
	}
}

func TestUnmarshalPubspecYaml(t *testing.T) {
	var tests = []struct {
		content string
		version string
	}{
		{`version: 1.2.0`, `1.2.0`},
		{`version: 1.2.1
name: newtify`, `1.2.1`},
		{`name: newtify
version: 1.2.2`, `1.2.2`},
		{`version: 1.2.3-release`, `1.2.3-release`},
		{`"version": 1.2.4`, `1.2.4`},
		{`version: "1.2.5"`, `1.2.5`},
		{``, ""},
	}
	for _, test := range tests {
		pubspecyaml := &PubspecYaml{}
		err := unmarshalPubspecYaml([]byte(test.content), pubspecyaml)
		if err != nil {
			t.Errorf("Error unmarshal: %v", err)
			if pubspecyaml.Version != test.version {
				t.Errorf("Unmarshal error for content: %s\n expected: %s, got: %s", test.content, test.version, pubspecyaml.Version)
			}
		}
	}
}

func TestUnmarshalErrorPubspecYaml(t *testing.T) {
	var tests = []struct {
		content string
		err     string
	}{
		{`v1.1`, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `v1.1` into githubservice.PubspecYaml"},
	}
	for _, test := range tests {
		pubspecyaml := &PubspecYaml{}
		if err := unmarshalPubspecYaml([]byte(test.content), pubspecyaml); fmt.Sprintf("%s", err) != test.err {
			t.Errorf("Error for content: %s\nexpected err: %v, got err: %v", test.content, test.err, err)
		}
	}
}

func TestErrorGetVersionPubspecYaml(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	psf := &pubspecyamlFetcher{}
	//test error get contents
	_, err := psf.GetVersion(&cp, AtcSettings{Path: "Flutter"})
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = psf.GetVersionDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error can't search verion
	cp.content = `name: atc`
	cp.err = nil
	_, err = psf.GetVersion(&cp, AtcSettings{Path: "Flutter"})
	if err != errNoVers {
		t.Errorf("err:%s  !=  noVersErr:%s", err, errNoVers)
	}

}

func TestPubspecYamlFetcherUnmarshalError(t *testing.T) {
	fetcher := pubspecyamlFetcher{}

	unmarshalPubspecYamlCopy := unmarshalPubspecYaml

	unmarshalPubspecYaml = failingPubspecYamlFetcherUnmarshal

	cp := mockContentProvider{"", nil}

	_, err := fetcher.GetVersion(&cp, AtcSettings{Path: "pubspec.yaml"})

	if err != errUnmarshal {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshalPubspecYaml = unmarshalPubspecYamlCopy
}
