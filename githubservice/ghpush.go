package githubservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type TagContent struct {
	Version string
}

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

func renderTagNameTemplate(templateString, version string) (string, error) {
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

func getShaByBehavior(push *github.WebHookPayload, behavior string) *string {
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

	var commitComment, newVersion, oldVersion string

	fetchType := detectFetchType(settings.Path)

	if fetchType != "" {
		var err error
		var reqError *RequestError
		fetcher := autoFetchers[fetchType]
		if fetcher == nil { //not default file
			if settings.RegexStr == "" {
				addComment(client, owner, repo, push.GetAfter(), fmt.Sprintf(".atc.yaml don't have regexstr for not default package manager file %s.", fetchType))
				return
			}
			fetcher = &customRegexFetcher{}
		} else {
			if settings.RegexStr != "" {
				commitComment += fmt.Sprintf("Used default regexStr in file %s. ", fetchType)
			}
		}
		oldVersion, err = fetcher.GetVersion(ghOldContentProviderPtr, *settings)
		if err != nil && !errors.Is(err, errHttpStatusCode) { //ignore http api error
			log.Printf("get prev version error for %q: %v", fullname, err)
			if errors.Is(err, errNoVers) || errors.Is(err, errNoGroupInConf) {
				addComment(client, owner, repo, push.GetAfter(), fmt.Sprintf("file %s with old version err: %v", fetchType, err))
			} else {
				addComment(client, owner, repo, push.GetAfter(), fmt.Sprintf("file %s with old version not found", fetchType))
			}
			return
		}
		newVersion, err = fetcher.GetVersion(ghNewContentProviderPtr, *settings)
		if err != nil {
			if errors.Is(err, errHttpStatusCode) {
				log.Printf("Wrong access status during getContent for installation %d for %q: %d", id, fullname, reqError.StatusCode)
			} else if errors.Is(err, errNoVers) || errors.Is(err, errNoGroupInConf) {
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
			oldVersion, err = fetcher.GetVersionUsingDefaultPath(ghOldContentProviderPtr)
			if err != nil && !errors.Is(err, errHttpStatusCode) { //ignore http api error
				log.Printf("get prev version error for %q, default path: %s, err: %v", fullname, defaultPath, err)
				continue
			}

			newVersion, err = fetcher.GetVersionUsingDefaultPath(ghNewContentProviderPtr)
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
		caption, err := renderTagNameTemplate(settings.Template, newVersion)
		if err != nil {
			log.Printf("error in go templates: %v", err)
			return
		}
		sha := *getShaByBehavior(push, settings.Behavior)
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

func CIPushAction() error {
	githubToken := os.Getenv("GITHUB_TOKEN")
	fullname := os.Getenv("GITHUB_REPOSITORY")
	commitSHA := os.Getenv("COMMIT_SHA")

	settings := &AtcSettings{
		Path:     os.Getenv("FILE_TYPE"),
		Behavior: os.Getenv("BEHAVIOR"),
		Template: os.Getenv("TEMPLATE"),
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	s := strings.Split(fullname, "/")
	owner := s[0]
	repo := s[1]

	commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, commitSHA, nil)
	if err != nil {
		return fmt.Errorf("error getting commit %s %v", commitSHA, err)
	}

	parents := commit.Parents
	if len(parents) == 0 {
		log.Printf("this branch has no older commits")
		return nil
	}

	ghOldContentProviderPtr := &ghContentProvider{
		owner:    owner,
		repo:     repo,
		ref:      parents[0].GetSHA(),
		ctx:      ctx,
		ghClient: client,
	}
	ghNewContentProviderPtr := &ghContentProvider{
		owner:    owner,
		repo:     repo,
		ref:      commitSHA,
		ctx:      ctx,
		ghClient: client,
	}

	var sha string

	caption, err := fetch(settings, ghOldContentProviderPtr, ghNewContentProviderPtr, fullname)
	if err != nil {
		return fmt.Errorf("fetch version error: %v", err)
	}

	if settings.Behavior == behaviorAfter {
		sha = commit.GetSHA()
	} else {
		sha = parents[0].GetSHA()
	}

	objType := "commit"
	timestamp := time.Now()

	tag := &github.Tag{
		Tag:     &caption,
		Message: &caption,
		Tagger: &github.CommitAuthor{
			Date:  &timestamp,
			Name:  commit.Commit.Author.Name,
			Email: commit.Commit.Author.Email,
			Login: commit.Commit.Author.Login,
		},
		Object: &github.GitObject{
			Type: &objType,
			SHA:  &sha,
		},
	}

	if err = addTagToCommit(client, owner, repo, tag); err != nil {
		return fmt.Errorf("error when adding tag to commit %q: %v", fullname, err)
	}

	log.Printf("Added a new version for %q: %q", fullname, caption)
	return nil
}

func fetch(settings *AtcSettings, ghOldContentProviderPtr,
	ghNewContentProviderPtr contentProvider, fullname string) (string, error) {
	fetchType := detectFetchType(settings.Path)
	var newVersion string
	var oldVersion string
	if fetchType != "" {
		var err error
		fetcher := autoFetchers[fetchType]
		if fetcher == nil {
			log.Printf("using custom fetcher")
			fetcher = &customRegexFetcher{}
		}

		oldVersion, err = fetcher.GetVersion(ghOldContentProviderPtr, *settings)
		if err != nil && !errors.Is(err, errHttpStatusCode) { //ignore http api error
			return "", fmt.Errorf("get prev version error for %q: %w", fullname, err)
		}

		log.Printf("old version %s", oldVersion)
		newVersion, err = fetcher.GetVersion(ghNewContentProviderPtr, *settings)
		if err != nil {
			return "", fmt.Errorf("get new version error for %q: %w", fullname, err)
		}
	} else {
		fetched := false
		for defaultPath, fetcher := range autoFetchers {
			var err error
			oldVersion, err = fetcher.GetVersionUsingDefaultPath(ghOldContentProviderPtr)
			if err != nil && !errors.Is(err, errHttpStatusCode) { //ignore http api error
				log.Printf("get prev version error for %q, default path: %s, err: %v", fullname, defaultPath, err)
				continue
			}

			newVersion, err = fetcher.GetVersionUsingDefaultPath(ghNewContentProviderPtr)
			if err == nil {
				fetched = true
				break
			} else {
				return "", fmt.Errorf("autofetcher error for %q: %w", defaultPath, err)
			}
		}

		if !fetched {
			return "", fmt.Errorf("unable to fetch version using known methods")
		}
	}

	if newVersion != oldVersion {
		log.Printf("There is a new version for %q! Old version: %q, new version: %q", fullname, oldVersion, newVersion)
		caption, err := renderTagNameTemplate(settings.Template, newVersion)
		if err != nil {
			return "", fmt.Errorf("error in go templates: %v", err)
		}
		return caption, nil
	}

	return "", nil
}
