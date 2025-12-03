package buildgradle

import (
	"errors"
	"fmt"
	"testing"

	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"
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
	fetcher := Fetcher{}

	cp := provider.MockContentProvider{Content: basicBuildGradle}

	vers, err := fetcher.GetVersion(&cp, settings.AtcSettings{Path: "build.gradle"})

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
		{`versionName "1.5"
		test{
			test2
		}`, "1.5"},
		{`test{
			test2
		}
		versionName "1.6"`, "1.6"},
		{`test{
			test2
		}
		versionName "1.7"
		test{
			test2
		}`, "1.7"},
		{`//versionName "1.81"
		versionName "1.8"`, "1.8"},
		{`versionName "1.9"
		//versionName "1.91"`, "1.9"},
		{`//versionName "2.01"
		versionName "2.0"
		//versionName "2.02"`, "2.0"},
	}
	for _, test := range tests {
		content := fmt.Sprintf(`
android {
    defaultConfig {
        versionCode 1
        %s

        testInstrumentationRunner "androidx.test.runner.AndroidJUnitRunner"
        }
    }
}
`, test.content)
		gradle := &BuildGradle{}
		err := unmarshalBuildGradle([]byte(content), gradle)
		if err != nil {
			t.Errorf("Error unmarshal: %v", err)
		}
		if gradle.Version != test.version {
			t.Errorf("Unmarshal error for content: %s\n expected: %s, got: %s", test.content, test.version, gradle.Version)
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
		{`defaultConfig {
			versionCode 1
			}`, "empty number version"},
		{`defaultConfig {
			versionCode 1
			versionName 111
			}`, "empty number version"},
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
	cp := provider.MockContentProvider{"", noContentErr}
	bgf := &Fetcher{}
	//test error get contents
	_, err := bgf.GetVersion(&cp, settings.AtcSettings{Path: "gradle"})
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = bgf.GetVersionUsingDefaultPath(&cp)
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
}
