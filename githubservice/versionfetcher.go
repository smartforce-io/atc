package githubservice

import "errors"

var (
	errNoVers        = errors.New("empty number version")
	errNoGroupInConf = errors.New("regexStr don't have group")
	errParsRegex     = errors.New("pasre regexStr error")
)

type VersionFetcher interface {
	GetVersion(ghContentProvider contentProvider, settings AtcSettings) (string, error)
	GetVersionDefaultPath(ghContentProvider contentProvider) (string, error)
}
