package githubservice

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestUnmarshallSample(t *testing.T) {
	file := "pom.xml"
	prefix := "n"
	sample := `
path: pom.xml
prefix: n
`
	settings := &AtcSettings{}

	if err := yaml.Unmarshal([]byte(sample), settings); err != nil {
		panic(err)
	}
	if settings.Path != file {
		t.Errorf("wrong settings File! Got %q, wanted %q", settings.Path, file)
	}
	if settings.Prefix != prefix {
		t.Errorf("wrong settings Prefix! Got %q, wanted %q", settings.Prefix, prefix)
	}
}
