package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

var basicBuildGradle = `
android {
    defaultConfig {
        versionCode 1
        versionName "1.7.1"

        testInstrumentationRunner "androidx.test.runner.AndroidJUnitRunner"
        }
    }
}
`

func TestBuildGradleFetcherBasic(t *testing.T) {
	fetcher := buildGradleFetcher{}

	cp := mockContentProvider{basicBuildGradle, nil}

	vers, err := fetcher.GetVersion(&cp, "build.gradle")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if vers != "1.7.1" {
		t.Errorf("wrong settings File! Got %q, wanted %q", vers, 5)
	}
}

func TestUnmarshalBuildGradle(t *testing.T) {
	var tests = []struct {
		content string
		version string
	}{
		{`versionName "1.1"`, "1.1"},
		{`versionName "1.2"
		versionCode 1`, "1.2"},
		{`versionCode 1
		versionName "1.3"`, "1.3"},
		{`versionName "1.4-release"`, "1.4-release"},
	}
	for _, test := range tests {
		gradle := &BuildGradle{}
		err := unmarshalBuildGradle([]byte(test.content), gradle)
		if err != nil {
			t.Errorf("Error unmarshal: %v", err)
			if gradle.Version != test.version {
				t.Errorf("Unmarshal error for content: %s\n expected: %s, got: %s", test.content, test.version, gradle.Version)
			}
		}
	}
}

func TestUnmarshalErrorBuildGradle(t *testing.T) {
	var tests = []struct {
		content string
		err     string
	}{
		{`versionName = 1.5`, "empty number version"},
		{`version "1.6"`, "empty number version"},
		{`v1.1`, "empty number version"},
		{``, "empty number version"},
	}
	for _, test := range tests {
		gradle := &BuildGradle{}
		if err := unmarshalBuildGradle([]byte(test.content), gradle); fmt.Sprintf("%s", err) != test.err {
			t.Errorf("Error for content: %s\nexpected err: %v, got err: %v", test.content, test.err, err)
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
}
