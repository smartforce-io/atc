package githubservice

import "errors"

type mockContentProvider struct {
	content string
	err     error
	reqErr  *RequestError
}

func (mockContentProvider *mockContentProvider) getContents(path string) (string, *RequestError, error) {
	return mockContentProvider.content, mockContentProvider.reqErr, mockContentProvider.err
}

var (
	errUnmarshal = errors.New("unmarshal error")
	errGeneral   = errors.New("weird error")
)
