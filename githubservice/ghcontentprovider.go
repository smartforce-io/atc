package githubservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v39/github"
)

type contentProvider interface {
	getContents(path string) (string, error)
}

type RequestError struct {
	StatusCode int
}

var errHttpStatusCode = errors.New("http status code error")

type ghContentProvider struct {
	owner    string
	repo     string
	ref      string
	ctx      context.Context
	ghClient *github.Client
}

func (ghcp *ghContentProvider) getContents(path string) (string, error) {

	fileContent, _, response, err := ghcp.ghClient.Repositories.GetContents(ghcp.ctx,
		ghcp.owner, ghcp.repo, path,
		&github.RepositoryContentGetOptions{Ref: ghcp.ref})

	if err != nil {
		return "", err
	}
	content, _ := fileContent.GetContent()

	if response.StatusCode != http.StatusOK {
		return content, errHttpStatusCode
	} else {
		return content, nil
	}
}
