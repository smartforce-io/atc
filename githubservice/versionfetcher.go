package githubservice

type VersionFetcher interface {
	GetVersion(ghContentProvider contentProvider, path string, tag *TagVersion) error
	GetVersionDefaultPath(ghContentProvider contentProvider, tagVerion *TagVersion) error
}
