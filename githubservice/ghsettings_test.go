package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

func failingAtcSettingsUnmarshal(content []byte, atcSettingsPtr *AtcSettings) error {
	return errUnmarshal
}

var basicConfig = `
path: contests/pom.xml
behavior: before
template: v{{.Version}}`

func TestBasicAtcSetting(t *testing.T) {

	var testsConfig = []struct {
		config   string
		path     string
		behavior string
		template string
	}{
		{`
path: contents/pom.xml
behavior: before
template: v{{.Version}}`, `contents/pom.xml`, `before`, `v{{.Version}}`},
		{`
path: build.gradle
behavior: after
template: vGR{{.Version}}`, `build.gradle`, `after`, `vGR{{.Version}}`},
	}

	cp := mockContentProvider{basicConfig, nil}

	_, err := getAtcSetting(&cp)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	for _, test := range testsConfig {
		cp = mockContentProvider{test.config, nil}
		settings, _ := getAtcSetting(&cp)
		if settings.Path != test.path {
			t.Errorf("wrong settings Path! Got %q, wanted %q", settings.Path, test.path)
		}
		if settings.Behavior != test.behavior {
			t.Errorf("wrong settings Begavior! Got %q, wanted %q", settings.Behavior, test.behavior)
		}
		if settings.Template != test.template {
			t.Errorf("wrong settings Template! Got %q, wanted %q", settings.Template, test.template)
		}
	}
}

func TestCheckSettingsForErrors(t *testing.T) {
	var tests = []struct {
		path             string
		behavior         string
		template         string
		expectedErrorStr string
	}{
		{"", "", "", `error config file .atc.yaml; path = ""`},
		{"/contents/pom.xml", "", "", `error config file .atc.yaml; path has prefix "/"`},
		{"contents//asd.txt", "", "", `error config file .atc.yaml; path has "//"`},
		{"contents/asd.txt", "", "", `error config file .atc.yaml: path no has suffix "pom.xml" or "build.gradle" or "package.json"`},
		{"contents/pom.xml", "", "", `error config file .atc.yaml; behavior = ""`},
		{"contents/pom.xml/", "bef", "", `error config file .atc.yaml: behavior no contains "before" or "after"`},
		{"contents/pom.xml", "after", "", `error config file .atc.yaml; template = ""`},
		{"contents/pom.xml", "after", "{.version}", `error config file .atc.yaml: template no contains "{{.Version}}"`},
		{"contents/pom.xml", "before", ".vers", `error config file .atc.yaml: template no contains "{{.Version}}"`},
		{"contents/pom.xml", "before", "v{{.Version}}V", fmt.Sprint(nil)},
	}

	for _, test := range tests {
		settings := &AtcSettings{test.path, test.behavior, test.template}
		err := checkSettingsForErrors(settings)
		if fmt.Sprint(err) != test.expectedErrorStr {
			t.Errorf("no takes error settings:%s\nexpected: %s, got: %s", settings, test.expectedErrorStr, err)
		}
	}
}

func TestUnmarshalDefault(t *testing.T) {
	var tests = []struct {
		atcYamlFile     string
		unexpectedError error
	}{
		{`
path: build.gradle
behavior: before
template: v{{.version}}`, errors.New(``)},
		{`
		path: build.gradle
		behavior: before
		template: v{{.version}}`, nil},
		{`
path: build.gradle
behavior: before
template: {{.version}}`, nil},
		{``, errors.New(``)},
	}

	for _, test := range tests {
		settings := &AtcSettings{}
		if err := unmarshal([]byte(test.atcYamlFile), settings); err == test.unexpectedError {
			t.Errorf("err unmarshal file:%s\n: %v", test.atcYamlFile, err)
		}
	}
}

func TestAtcSettingGetContentsError(t *testing.T) {
	confFilStr := `
path: contents/pom.xml
behavior: before
template: v{{.version}}`
	emptySettings := &AtcSettings{}

	cp := mockContentProvider{content: confFilStr, err: errGeneral}
	set, err := getAtcSetting(&cp)

	if set != emptySettings && err != nil {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errGeneral)
	}
}

func TestAtcSettingUnmarshalError(t *testing.T) {
	unmarshalcp := unmarshal

	unmarshal = failingAtcSettingsUnmarshal

	cp := mockContentProvider{basicConfig, nil}
	_, err := getAtcSetting(&cp)

	if fmt.Sprint(err) != `error config file .atc.yaml; can't unmarshal file` {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshal = unmarshalcp
}
