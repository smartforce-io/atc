package githubservice

type VersionFetcher interface {
	GetVersion(ghContentProvider contentProvider, path string) (string, error)
}
