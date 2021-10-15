package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

func TestUnmarshalGradleProperties(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{"version=1.0.32", "1.0.32"},
		{`version=1.0.32
		org.gradle.caching=true`, "1.0.32"},
		{`# ....
		version=1.0.32`, "1.0.32"},
		{`# ....version=1.0.32`, "1.0.32"},
		{`version=1.4-relise`, "1.4-relise"},
		{`version = 1.0.32`, ""},
		{`version="1.1"`, "\"1.1\""},
		{`v1.1`, ""},
		{``, ""},
	}
	for _, test := range tests {
		gradle := &GradleProperties{}
		if err := unmarshalGradleProperties([]byte(test.input), gradle); gradle.Version != test.output {
			t.Errorf("err: %q\n, gradle.Version = %s\n test.input= %s", err, gradle.Version, test.input)
		}
	}
}

func TestErrorGetVersionGradle(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	gpf := &gradlePropertiesFetcher{}
	//test error get contents
	err := gpf.GetVersion(&cp, "gradle", &TagVersion{})
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	err = gpf.GetVersionDefaultPath(&cp, &TagVersion{})
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error can't search verion
	noVersErr := "gradle.properties does not containt number version"
	cp.err = nil
	err = gpf.GetVersion(&cp, "gradle", &TagVersion{})
	if fmt.Sprint(err) != noVersErr {
		t.Errorf("err:%s  !=  noVersErr:%s", err, noVersErr)
	}
}
