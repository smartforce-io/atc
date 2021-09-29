package githubservice

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
)

var autoFetchers = map[string]VersionFetcher{
	"pom.xml":           &pomXmlFetcher{},
	"grable.properties": &gradlePropertiesFetcher{},
	".npmrc":            &npmrcFetcher{},
}

func detectFetchType(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

func madeСaptionToTemplate(template, version string) string {
	if template == "" {
		return "" // need ckeck on ""????
	}
	return strings.Replace(template, ".version", version, -1)
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

	ghOldContentProviderPtr := &ghContentProvider{
		owner:    owner,
		repo:     repo,
		sha1:     push.GetBefore(),
		ctx:      ctx,
		ghClient: client,
	}
	ghNewContentProviderPtr := &ghContentProvider{
		owner:    owner,
		repo:     repo,
		ctx:      ctx,
		ghClient: client,
	}

	settings, err := getAtcSetting(ghNewContentProviderPtr)
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
		oldVersion, err = fetcher.GetVersion(ghOldContentProviderPtr, settings.Path)
		if err != nil && err != errHttpStatusCode { //ignore http api error
			log.Printf("get prev version error for %q: %v", fullname, err)
			return
		}
		newVersion, err = fetcher.GetVersion(ghNewContentProviderPtr, settings.Path)
		if err != nil {
			if err == errHttpStatusCode {
				log.Printf("Wrong access status during getContent for installation %d for %q: %d", id, fullname, reqError.StatusCode)
			} else {
				log.Printf("get version error for %q: %v", fullname, err)
			}
			return
		}
	} else {
		fetched := false
		for defaultPath, fetcher := range autoFetchers {
			var err error
			oldVersion, _ = fetcher.GetVersion(ghOldContentProviderPtr, defaultPath)
			if err != nil && err != errHttpStatusCode { //ignore http api error
				log.Printf("get prev version error for %q: %v", fullname, err)
				return
			}

			newVersion, err = fetcher.GetVersion(ghNewContentProviderPtr, defaultPath)

			if err == nil {
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

		caption := madeСaptionToTemplate(settings.Template, newVersion)
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
