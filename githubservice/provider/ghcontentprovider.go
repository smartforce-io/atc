package provider

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v39/github"
)

type ContentProvider interface {
	GetContents(path string) (string, error)
}

type RequestError struct {
	StatusCode int
}

var ErrHttpStatusCode = errors.New("http status code error")

type GhContentProvider struct {
	Owner    string
	Repo     string
	Ref      string
	Ctx      context.Context
	GhClient *github.Client
}

func (ghcp *GhContentProvider) GetContents(path string) (string, error) {

	fileContent, _, response, err := ghcp.GhClient.Repositories.GetContents(ghcp.Ctx,
		ghcp.Owner, ghcp.Repo, path,
		&github.RepositoryContentGetOptions{Ref: ghcp.Ref})

	if err != nil {
		return "", err
	}
	content, _ := fileContent.GetContent()

	if response.StatusCode != http.StatusOK {
		return content, ErrHttpStatusCode
	} else {
		return content, nil
	}
}
