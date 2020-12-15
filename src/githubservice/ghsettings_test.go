package githubservice

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestUnmarshallSample(t *testing.T) {
	file := "pom.xml"
	prefix := "n"
	field := "version"
	sample := `
file: pom.xml
prefix: n
field: version
`
	settings := &AtcSettings{}

	if err := yaml.Unmarshal([]byte(sample), settings); err != nil {
		panic(err)
	}
	if settings.File != file { t.Errorf("wrong settings File! Got %q, wanted %q", settings.File, file) }
	if settings.Prefix != prefix { t.Errorf("wrong settings Prefix! Got %q, wanted %q", settings.Prefix, prefix) }
	if settings.Field != field { t.Errorf("wrong settings Field! Got %q, wanted %q", settings.Field, field) }
}