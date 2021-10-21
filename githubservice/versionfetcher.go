package githubservice

import "errors"

var (
	errNoVers        = errors.New("empty number version")
	errNoGroupInConf = errors.New("regexStr don't have group")
)

type VersionFetcher interface {
	GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error)
	GetVersionDefaultPath(ghContentProvider contentProvider) (string, error)
}
