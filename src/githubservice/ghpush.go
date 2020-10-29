package githubservice

import (
	"context"
	"github.com/google/go-github/github"
	"log"
	"net/http"
)

func PushAction(push *github.WebHookPayload, id int64) {
	const versionSource = "pom.xml"

	jwt, err := getJwt()
	if err != nil {
		log.Printf("getJWT Error: %v", err)
		return
	}

	ctx := context.Background()
	client := getGithubClient(jwt, ctx)
	inst, resp, err := client.Apps.CreateInstallationToken(ctx, id, nil)
	if err != nil {
		log.Printf("create access token error: %v", err)
		return
	}
	if resp.StatusCode !=  http.StatusCreated {
		log.Printf("Wrong access status during create access token for installation %d: %s", id, resp.Status)
		return
	}

	client = getGithubClient(inst.GetToken(), ctx)
	owner := push.GetRepo().GetOwner().GetName()
	repo := push.GetRepo().GetName()
	f, _, resp, err := client.Repositories.GetContents( ctx, owner, repo, versionSource, nil)
	if err != nil {
		log.Printf("get contents error: %v", err)
		return
	}
	if resp.StatusCode !=  http.StatusOK {
		log.Printf("Wrong access status during getContent for installation %d: %s", id, resp.Status)
		return
	}

	log.Print(f.GetContent())

	for _, commit := range push.Commits {
		for _, modified := range commit.Modified {
			if modified == versionSource {
				log.Printf("There is a new version!!!!!")
			}
		}
	}
}