package packagejson

import (
	"encoding/json"

	"github.com/smartforce-io/atc/githubservice/provider"

	"github.com/smartforce-io/atc/githubservice/fetcher"
	"github.com/smartforce-io/atc/githubservice/settings"
)

type PackageJson struct {
	Version string `PackageJson:"version"`
}

type Fetcher struct {
}

var unmarshalPackageJson = func(content []byte, packagejsonPtr *PackageJson) error {
	return json.Unmarshal(content, packagejsonPtr)
}

func (packagejsonFetcher *Fetcher) GetVersion(ghContentProvider provider.ContentProvider, settings settings.AtcSettings) (string, error) {
	content, err := ghContentProvider.GetContents(settings.Path)
	if err != nil {
		return "", err
	}
	packagejson := &PackageJson{}
	if err := unmarshalPackageJson([]byte(content), packagejson); err != nil {
		return "", err
	}
	if packagejson.Version == "" {
		return "", fetcher.ErrNoVers
	}
	return packagejson.Version, nil
}

func (packagejsonFetcher *Fetcher) GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error) {
	return packagejsonFetcher.GetVersion(ghContentProvider, settings.AtcSettings{Path: "package.json"})
}
