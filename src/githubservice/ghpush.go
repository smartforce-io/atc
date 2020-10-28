package githubservice

import (
	"github.com/google/go-github/github"
	"log"
)

func PushAction(push *github.WebHookPayload) {
	const versionSource = "pom.xml"

	for _, commit := range push.Commits {
		for _, modified := range commit.Modified {
			if modified == versionSource && isNewVersion(versionSource) {
				log.Printf("There is a new version!!!!!")
			}
		}
	}
}

func isNewVersion(s string) bool {
	return true
}