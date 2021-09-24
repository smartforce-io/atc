package githubservice

import (
	"context"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type ClientProvider interface {
	Get(token string, ctx context.Context) *github.Client
}

type GithubClientProvider struct {
}

func (githubClientProvider *GithubClientProvider) Get(token string, ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
