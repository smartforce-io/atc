package githubservice

import (
	"fmt"
	"os"
	"path/filepath"
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
	// typel := "maven"
	file := "pom.xml"
	behavior := "before"
	template := "v{{.version}}"
	//prefix := "n"

	cp := mockContentProvider{basicConfig, nil}

	settings, err := getAtcSetting(&cp)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	// if settings.Type != typel {
	// 	t.Errorf("wrong settings File! Got %q, wanted %q", settings.Type, typel)
	// }
	if settings.Path != file {
		if filepath.Base(settings.Path) != file {
			t.Errorf("wrong settings File! Got %q, wanted %q", settings.Path, file)
		}
	}

	if settings.Behavior != behavior {
		t.Errorf("wrong settings Behavior! Got %q, wanted %q", settings.Behavior, behavior)
	}
	if settings.Template != template {
		t.Errorf("wrong settings Template! Got %q, wanted %q", settings.Template, template)
	}
	// if settings.Prefix != prefix {
	// 	t.Errorf("wrong settings Prefix! Got %q, wanted %q", settings.Prefix, prefix)
	// }
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
