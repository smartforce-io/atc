package githubservice

import (
	"testing"
)

func failingAtcSettingsUnmarshal(content []byte, atcSettingsPtr *AtcSettings) error {
	return errUnmarshal
}

var basicConfig = `
path: pom.xml
prefix: n
`

func TestBasicAtcSetting(t *testing.T) {
	file := "pom.xml"
	prefix := "n"

	cp := testContentProvider{basicConfig, nil, nil}

	settings, err := getAtcSetting(&cp)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if settings.Path != file {
		t.Errorf("wrong settings File! Got %q, wanted %q", settings.Path, file)
	}
	if settings.Prefix != prefix {
		t.Errorf("wrong settings Prefix! Got %q, wanted %q", settings.Prefix, prefix)
	}
}
func TestWrongAtcSettingFile(t *testing.T) {
	reqErr := RequestError{StatusCode: 404}
	cp := testContentProvider{content: "", reqErr: &reqErr, err: nil}
	_, err := getAtcSetting(&cp)

	if err != errFailedResponse {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errFailedResponse)
	}

}
func TestAtcSettingGeneralError(t *testing.T) {
	cp := testContentProvider{content: "", reqErr: nil, err: errGeneral}
	_, err := getAtcSetting(&cp)

	if err != errGeneral {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errGeneral)
	}

}

func TestAtcSettingUnmarshalError(t *testing.T) {
	unmarshalcp := unmarshal

	unmarshal = failingAtcSettingsUnmarshal

	cp := testContentProvider{basicConfig, nil, nil}
	_, err := getAtcSetting(&cp)

	if err != errUnmarshal {
		t.Errorf("Invalid error, Got %v, wanted %v", err, errUnmarshal)
	}

	unmarshal = unmarshalcp

}
