package githubservice

import "errors"

type testContentProvider struct {
	content string
	err     error
	reqErr  *RequestError
}

func (testContentProvider *testContentProvider) getContents(path string) (string, *RequestError, error) {
	return testContentProvider.content, testContentProvider.reqErr, testContentProvider.err
}

var (
	errUnmarshal = errors.New("unmarshal error")
	errGeneral   = errors.New("weird error")
)
