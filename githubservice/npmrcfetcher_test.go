package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

func TestUnmarshalNpmrc(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{`npm config set init.version "0.0.1"`, `0.0.1`},
		{`npm config set init.version "0.0.1"
		npm config set init.license "MIT"`, `0.0.1`},
		{`npm config set init.author.url "http://hiro.snowcrash.io"
		npm config set init.version "0.0.1"`, `0.0.1`},
		{`npm config set init.author.url "http://hiro.snowcrash.io" npm config set init.version "0.0.1"`, `0.0.1`},
		{`npm config set init.version "0.0.1-relise"`, "0.0.1-relise"},
		{`init.version = 1.0.32`, ""},
		{`v1.1`, ""},
		{``, ""},
	}
	for _, test := range tests {
		npm := &Npmrc{}
		if err := unmarshalNpmrc([]byte(test.input), npm); npm.Version != test.output {
			t.Errorf("err: %q\n, npm.Version = %s\n test.input= %s", err, npm.Version, test.input)
		}
	}
}

func TestErrorGetVersionNpm(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	nf := &npmrcFetcher{}
	//test error get contents
	_, err := nf.GetVersion(&cp, "npm")
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = nf.GetVersionDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error can't search verion
	noVersErr := ".npmrc does not containt number version"
	cp.err = nil
	_, err = nf.GetVersion(&cp, "npm")
	if fmt.Sprint(err) != noVersErr {
		t.Errorf("err:%s  !=  noVersErr:%s", err, noVersErr)
	}
}
