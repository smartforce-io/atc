package githubservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/github"
)

var (
	errCreateTagWrongStatus = errors.New("wrong status for create a tag")
	errCreateRefWrongStatus = errors.New("wrong status for create a ref")
)

func addTagToCommit(client *github.Client, owner, repo string, tag *github.Tag) error {
	t, resp, err := client.Git.CreateTag(context.Background(), owner, repo, tag)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errCreateTagWrongStatus
	}

	refs := "refs/tags/" + t.GetTag()
	_, resp, err = client.Git.CreateRef(context.Background(), owner, repo, &github.Reference{
		Ref: &refs,
		Object: &github.GitObject{
			SHA: t.SHA,
		},
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errCreateRefWrongStatus
	}
	return nil
}
