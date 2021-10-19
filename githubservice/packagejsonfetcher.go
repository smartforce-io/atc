package githubservice

import (
	"encoding/json"
)

type PackageJson struct {
	Version string `PackageJson:"version"`
}

type packagejsonFetcher struct {
}

var unmarshalPackageJson = func(content []byte, packagejsonPtr *PackageJson) error {
	return json.Unmarshal(content, packagejsonPtr)
}

func (packagejsonFetcher *packagejsonFetcher) GetVersion(ghContentProvider contentProvider, path string) (string, error) {
	content, err := ghContentProvider.getContents(path)
	if err != nil {
		return "", err
	}
	packagejson := &PackageJson{}
	if err := unmarshalPackageJson([]byte(content), packagejson); err != nil {
		return "", err
	}
	if packagejson.Version == "" {
		return "", errNoVers
	}
	return packagejson.Version, nil
}

func (packagejsonFetcher *packagejsonFetcher) GetVersionDefaultPath(ghContentProvider contentProvider) (string, error) {
	return packagejsonFetcher.GetVersion(ghContentProvider, "package.json")
}
