package githubservice

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v39/github"
)

type TagContent struct {
	Version string
}

const (
	behaviorBefore = "before"
)

var autoFetchers = map[string]VersionFetcher{
	"pom.xml":      &pomXmlFetcher{},
	"build.gradle": &buildGradleFetcher{},
	"package.json": &packagejsonFetcher{},
	"pubspec.yaml": &pubspecyamlFetcher{},
}

func detectFetchType(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

func madeСaptionToTemplate(templateString, version string) (string, error) {
	buf := new(bytes.Buffer)
	tagContent := TagContent{version}
	tmplFuncMap := template.FuncMap{
		"Time": func() time.Time { return time.Now() },
	}
	tmpl, err := template.New("template tagContent").Funcs(tmplFuncMap).Parse(templateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(buf, tagContent)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func madeShaToBehavior(push *github.WebHookPayload, behavior string) *string {
	if strings.ToLower(behavior) == behaviorBefore {
		return push.Before
	}
	return push.After
}

func createBranchToClientProvider(settings *AtcSettings, push *github.WebHookPayload) string {
	if settings.Branch != "" {
		return settings.Branch
	}
	return push.GetRepo().GetDefaultBranch()
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
		ref:      push.GetBefore(),
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
		log.Println("err. send user: ", err)
		addComment(client, owner, repo, push.GetAfter(), fmt.Sprint(err))
		return
	}

	ghNewContentProviderPtr.ref = createBranchToClientProvider(settings, push)
	if push.GetRef() != "refs/heads/"+ghNewContentProviderPtr.ref { // checking which branch is in work
		return
	}

	commitComment := ""
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
			if err == errNoVers {
				addComment(client, owner, repo, push.GetAfter(), fmt.Sprintf("file %s with old version err: %v", fetchType, err))
			} else {
				addComment(client, owner, repo, push.GetAfter(), fmt.Sprintf("file %s with old version not found", fetchType))
			}
			return
		}
		newVersion, err = fetcher.GetVersion(ghNewContentProviderPtr, settings.Path)
		if err != nil {
			if err == errHttpStatusCode {
				log.Printf("Wrong access status during getContent for installation %d for %q: %d", id, fullname, reqError.StatusCode)
			} else if err == errNoVers {
				log.Printf("get version error for %q: %v", fullname, err)
				addComment(client, owner, repo, push.GetAfter(), fmt.Sprintf("file %s with new version err: %v", fetchType, err))
			} else {
				log.Printf("get version error for %q: %v", fullname, err)
				addComment(client, owner, repo, push.GetAfter(), fmt.Sprintf("file %s with new version not found", fetchType))
			}
			return
		}
	} else {
		commitComment = `File .atc.yaml not found or path = "". `
		fetched := false
		for defaultPath, fetcher := range autoFetchers {
			var err error
			oldVersion, err = fetcher.GetVersionDefaultPath(ghOldContentProviderPtr)
			if err != nil && err != errHttpStatusCode { //ignore http api error
				log.Printf("get prev version error for %q, default path: %s, err: %v", fullname, defaultPath, err)
				continue
			}

			newVersion, err = fetcher.GetVersionDefaultPath(ghNewContentProviderPtr)

			if err == nil {
				fetched = true
				commitComment += "Used default settings. "
				break
			} else {
				log.Printf("autofetcher error for %q: %v", defaultPath, err)
			}
		}
		if !fetched {
			commitComment += "Not found supported package manager."
			addComment(client, owner, repo, push.GetAfter(), commitComment)
			log.Printf("Unable to fetch version using known methods!") //probably should be comment
			return
		}
	}

	if newVersion != oldVersion {
		log.Printf("There is a new version for %q! Old version: %q, new version: %q", fullname, oldVersion, newVersion)
		caption, err := madeСaptionToTemplate(settings.Template, newVersion)
		if err != nil {
			log.Printf("error in go templates: %v", err)
			return
		}
		sha := *madeShaToBehavior(push, settings.Behavior)
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
			addComment(client, owner, repo, sha, fmt.Sprintf("can't add tag to commit, error : %v", err))
			return
		}

		commitComment += fmt.Sprintf("Added a new version for %q: %q", fullname, caption)
		addComment(client, owner, repo, sha, commitComment)
	}
}
