package githubservice

import (
	"fmt"
	"os"
	"testing"
)

func failingAtcSettingsUnmarshal(content []byte, atcSettingsPtr *AtcSettings) error {
	return errUnmarshal
}

var basicConfig = `
type: maven
path: /contests/pom.xml
behavior: before
template: v{{.version}}
//prefix: n
`

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
template: v{{.version}}
`, `contents/pom.xml`, `before`, `v{{.version}}`},
		{`
path: gradle.properties
behavior: after
template: vGR{{.version}}
`, `gradle.properties`, `after`, `vGR{{.version}}`},
		{`
path: 
behavior: 
template: 
`, ``, ``, ``},
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

func TestUnmarshalDefault(t *testing.T) {
	fileByte, err := os.ReadFile("../../.atc.yaml")
	if err != nil {
		t.Errorf("Err read file %v", err)
	}
	settings := &AtcSettings{}
	if err := unmarshal(fileByte, settings); err != nil {
		t.Errorf("err unmarshal: %v", err)
	}

}
func TestAtcSettingGeneralError(t *testing.T) {
	cp := mockContentProvider{content: "", err: errGeneral}
	set, err := getAtcSetting(&cp)

	fmt.Printf("t: %v\n", set)

	if err != errGeneral {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errGeneral)
	}

}

func TestAtcSettingUnmarshalError(t *testing.T) {
	unmarshalcp := unmarshal

	unmarshal = failingAtcSettingsUnmarshal

	cp := mockContentProvider{basicConfig, nil}
	_, err := getAtcSetting(&cp)

	if err != errUnmarshal {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshal = unmarshalcp

}
