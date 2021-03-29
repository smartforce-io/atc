package githubservice

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
)

type contentProvider interface {
	getContents(path string) (string, *RequestError, error)
}

type RequestError struct {
	StatusCode int
}

type ghContentProvider struct {
	owner    string
	repo     string
	sha1     string
	ctx      context.Context
	ghClient *github.Client
}

func (ghcp *ghContentProvider) getContents(path string) (string, *RequestError, error) {

	fileContent, _, response, err := ghcp.ghClient.Repositories.GetContents(ghcp.ctx,
		ghcp.owner, ghcp.repo, path,
		&github.RepositoryContentGetOptions{Ref: ghcp.sha1})

	if err != nil {
		return "", nil, err
	}
	content, _ := fileContent.GetContent()

	if response.StatusCode != http.StatusOK {
		return content, &RequestError{StatusCode: response.StatusCode}, nil
	} else {
		return content, nil, nil
	}
}
