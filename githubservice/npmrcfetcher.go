package githubservice

import (
	"errors"
	"regexp"
)

type Npmrc struct {
	Version string `Npmrc:"version"`
}

type npmrcFetcher struct {
}

var unmarshalNpmrc = func(content []byte, npmrcPtr *Npmrc) error {
	regex, _ := regexp.Compile("npm config set init.version \"([^\t\n\f\r]*?)\"")
	res := regex.FindStringSubmatch(string(content))
	if len(res) == 0 {
		return errors.New(".npmrc does not containt number version")
	}
	npmrcPtr.Version = res[1]
	return nil
}

func (npmrcFetcher *npmrcFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, error) {
	content, err := ghContentProvider.getContents(path)
	if err != nil {
		return "", err
	}
	npmrc := &Npmrc{}
	if err := unmarshalNpmrc([]byte(content), npmrc); err != nil {
		return "", err
	}
	return npmrc.Version, nil
}

func (npmrcFetcher *npmrcFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return npmrcFetcher.GetVersion(ghContentProvider, ".npmrc")
}
