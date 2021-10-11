package githubservice

import (
	"context"
	"log"

	"github.com/google/go-github/v39/github"
)

func addComment(client *github.Client, owner, repo, sha, text string) {
	if _, _, err := client.Repositories.CreateComment(context.Background(), owner, repo, sha, &github.RepositoryComment{
		Body: &text,
	}); err != nil {
		log.Printf("add comment error for %s/%s: %v", owner, repo, err)
	}
	log.Printf("add comment: %s", text)
}
