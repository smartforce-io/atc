package githubservice

type VersionFetcher interface {
	GetVersion(ghContentProvider *ghContentProvider, path string) (string, *RequestError, error)
}
