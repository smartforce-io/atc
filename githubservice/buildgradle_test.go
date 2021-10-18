package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

func TestUnmarshalBuildGradle(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{`versionName "1.1"`, "1.1"},
		{`versionName "1.2"
		versionCode 1`, "1.2"},
		{`versionCode 1
		versionName "1.3"`, "1.3"},
		{`versionName "1.4-release"`, "1.4-release"},
		{`versionName = 1.5`, ""},
		{`version "1.6"`, ""},
		{`v1.1`, ""},
		{``, ""},
	}
	for _, test := range tests {
		gradle := &BuildGradle{}
		if err := unmarshalBuildGradle([]byte(test.input), gradle); gradle.Version != test.output {
			t.Errorf("err: %q\n, gradle.Version = %s\n test.input= %s", err, gradle.Version, test.input)
		}
	}
}

func TestErrorGetVersionGradle(t *testing.T) {
	noContentErr := errors.New("can't get content")
	cp := mockContentProvider{"", noContentErr}
	bgf := &buildGradleFetcher{}
	//test error get contents
	_, err := bgf.GetVersion(&cp, "gradle")
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = bgf.GetVersionDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error can't search verion
	noVersErr := "build.gradle does not containt number version"
	cp.err = nil
	_, err = bgf.GetVersion(&cp, "gradle")
	if fmt.Sprint(err) != noVersErr {
		t.Errorf("err:%s  !=  noVersErr:%s", err, noVersErr)
	}
}
