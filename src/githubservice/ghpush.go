package githubservice

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
)

var autoFetchers = map[string]VersionFetcher{
	"pom.xml": pomXmlFetcher{},
}

func detectFetchType(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

func PushAction(push *github.WebHookPayload, clientProvider ClientProvider) {
	id := *push.Installation.ID

	token, err := getAccessToken(id, clientProvider)
	if err != nil {
		log.Printf("getAccessToken Error: %v", err)
		return
	}
	owner := push.GetRepo().GetOwner().GetName()
	repo := push.GetRepo().GetName()
	fullname := push.GetRepo().GetFullName()
	ctx := context.Background()
	client := clientProvider.Get(token, ctx)

	ghOldContentProvider := ghContentProvider{
		owner:    owner,
		repo:     repo,
		sha1:     push.GetBefore(),
		ctx:      ctx,
		ghClient: client,
	}
	ghNewContentProvider := ghContentProvider{
		owner:    owner,
		repo:     repo,
		ctx:      ctx,
		ghClient: client,
	}

	settings, err := getAtcSetting(&ghNewContentProvider)
	if err != nil {
		settings = &AtcSettings{} //blank settings
	}

	newVersion := ""
	oldVersion := ""
	fetchType := detectFetchType(settings.Path)

	if fetchType != "" {
		var err error
		var reqError *RequestError
		fetcher := autoFetchers[fetchType]
		oldVersion, _, err = fetcher.GetVersion(&ghOldContentProvider, settings.Path) //ignore http api error
		if err != nil {
			log.Printf("get prev version error for %q: %v", fullname, err)
			return
		}
		newVersion, reqError, err = fetcher.GetVersion(&ghNewContentProvider, settings.Path)
		if err != nil {
			log.Printf("get version error for %q: %v", fullname, err)
			return
		}
		if reqError != nil {
			log.Printf("Wrong access status during getContent for installation %d for %q: %d", id, fullname, reqError.StatusCode)
			return
		}
	} else {
		fetched := false
		for defaultPath, fetcher := range autoFetchers {
			var err error
			var reqError *RequestError
			oldVersion, _, _ = fetcher.GetVersion(&ghOldContentProvider, defaultPath)

			newVersion, reqError, err = fetcher.GetVersion(&ghNewContentProvider, defaultPath)

			if reqError == nil && err == nil {
				fetched = true
				break
			} else {
				log.Printf("autofetcher error for %q: %v", defaultPath, err)
			}
		}
		if !fetched {
			log.Printf("Unable to fetch version using known methods!") //probably should be comment
			return
		}

	}

	if newVersion != oldVersion {
		log.Printf("There is a new version for %q! Old version: %q, new version: %q", fullname, oldVersion, newVersion)

		caption := settings.Prefix + newVersion
		sha := push.GetAfter()
		objType := "commit"
		timestamp := time.Now()

		tag := &github.Tag{
			Tag:     &caption,
			Message: &caption,
			Tagger: &github.CommitAuthor{
				Date:  &timestamp,
				Name:  push.GetPusher().Name,
				Email: push.GetPusher().Email,
				Login: push.GetPusher().Login,
			},
			Object: &github.GitObject{
				Type: &objType,
				SHA:  &sha,
			},
		}

		if err := addTagToCommit(client, owner, repo, tag); err != nil {
			log.Printf("addTagToCommit Error for %q: %v", fullname, err)
			return
		}

		cmnt := fmt.Sprintf("Added a new version for %q: %q", fullname, caption)
		addComment(client, owner, repo, sha, cmnt)
	}
}
