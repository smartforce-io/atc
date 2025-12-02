package fetcher

import (
	"errors"

	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"
)

var (
	ErrNoVers        = errors.New("empty number version")
	ErrNoGroupInConf = errors.New("regexStr don't have group")
	ErrParsRegex     = errors.New("pasre regexStr error")
)

type VersionFetcher interface {
	GetVersion(ghContentProvider provider.ContentProvider, settings settings.AtcSettings) (string, error)
	GetVersionUsingDefaultPath(ghContentProvider provider.ContentProvider) (string, error)
}
