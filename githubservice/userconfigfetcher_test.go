package githubservice

import (
	"errors"
	"fmt"
	"testing"
)

var basicUserConfig = `
name: test
project: testt
vers: 1.0.1
end
}
`

func TestUserConfigFetcherBasic(t *testing.T) {
	fetcher := userConfigFetcher{}

	cp := mockContentProvider{basicUserConfig, nil}

	vers, err := fetcher.GetVersion(&cp, AtcSettings{Path: "test.txt", RegexStr: "vers: (.+)"})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if vers != "1.0.1" {
		t.Errorf("wrong settings File! Got %q, wanted %q", vers, 5)
	}
}

func TestUnmarshalUserConfig(t *testing.T) {
	var tests = []struct {
		content  string
		regexstr string
		version  string
	}{
		{`vers: "1.1"`, "vers: (.+)", "1.1"},
		{`vers: "1.2"
		versionCode 1`, "vers: (.+)", "1.2"},
		{`vers: 1
		versionName "1.3"`, "vers: (.+)", "1.3"},
		{`vers: "1.4-release"`, "vers: (.+)", "1.4-release"},
	}
	for _, test := range tests {
		userConf := &UserConfig{}
		err := unmarshalUserConfig([]byte(test.content), test.regexstr, userConf)
		if err != nil {
			t.Errorf("Error unmarshal: %v", err)
			if userConf.Version != test.version {
				t.Errorf("Unmarshal error for content: %s\n expected: %s, got: %s", test.content, test.version, userConf.Version)
			}
		}
	}
}

func TestUnmarshalErrorUserConfig(t *testing.T) {
	var tests = []struct {
		content  string
		regexstr string
		err      string
	}{
		{`versionName = 1.5`, "vers: (.{}", "error parsing regexp: missing closing ): `vers: (.{}`"},
		{`vers: 11"`, "versious: 1", "empty number version"},
		{``, "", "regexStr don't have group"},
	}
	for _, test := range tests {
		userConf := &UserConfig{}
		if err := unmarshalUserConfig([]byte(test.content), test.regexstr, userConf); fmt.Sprintf("%s", err) != test.err {
			t.Errorf("Error for content: %s\nexpected err: %v, got err: %v", test.content, test.err, err)
		}
	}
}

func TestErrorGetVersionUserConfig(t *testing.T) {
	noContentErr := errors.New("can't get content")
	defaultPathErr := "UserConfig don't have default path"
	cp := mockContentProvider{"", noContentErr}
	usf := &userConfigFetcher{}
	//test error get contents
	_, err := usf.GetVersion(&cp, AtcSettings{Path: "test"})
	if err != noContentErr {
		t.Errorf("err:%s  !=  noContentErr:%s", err, noContentErr)
	}
	//test error get contents when use DefaultPath
	_, err = usf.GetVersionUsingDefaultPath(&cp)
	if fmt.Sprintf("%s", err) != defaultPathErr {
		t.Errorf("err:%s  !=  defaultPathErr:%s", err, defaultPathErr)
	}
}
