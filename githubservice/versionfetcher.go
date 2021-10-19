package githubservice

import "errors"

var (
	errNoVers = errors.New("empty number version")
)

type VersionFetcher interface {
	GetVersion(ghContentProvider contentProvider, path string) (string, error)
	GetVersionDefaultPath(ghContentProvider contentProvider) (string, error)
}
