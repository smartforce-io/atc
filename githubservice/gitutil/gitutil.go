package gitutil

import (
	"context"
	"errors"
	"github.com/google/go-github/v39/github"
	"log"
	"net/http"
)

var (
	errCreateTagWrongStatus = errors.New("wrong status for create a tag")
	errCreateRefWrongStatus = errors.New("wrong status for create a ref")
)

func AddComment(client *github.Client, owner, repo, sha, text string) {
	if _, _, err := client.Repositories.CreateComment(context.Background(), owner, repo, sha, &github.RepositoryComment{
		Body: &text,
	}); err != nil {
		log.Printf("add comment error for %s/%s: %v", owner, repo, err)
	}
}

func AddTagToCommit(client *github.Client, owner, repo string, tag *github.Tag) error {
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
