package githubservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/smartforce-io/atc/envvars"

	"github.com/google/go-github/v39/github"
)

var testWebhookPayload = `
{
	"ref": "refs/tags/simple-tag",
	"before": "6113728f27ae82c7b1a177c8d03f9e96e0adf246",
	"after": "0000000000000000000000000000000000000000",
	"created": false,
	"deleted": true,
	"forced": false,
	"base_ref": null,
	"compare": "https://github.com/Codertocat/Hello-World/compare/6113728f27ae...000000000000",
	"commits": [],
	"head_commit": null,
	"repository": {
	  "id": 186853002,
	  "node_id": "MDEwOlJlcG9zaXRvcnkxODY4NTMwMDI=",
	  "name": "Hello-World",
	  "full_name": "Codertocat/Hello-World",
	  "private": false,
	  "owner": {
		"name": "Codertocat",
		"email": "21031067+Codertocat@users.noreply.github.com",
		"login": "Codertocat",
		"id": 21031067,
		"node_id": "MDQ6VXNlcjIxMDMxMDY3",
		"avatar_url": "https://avatars1.githubusercontent.com/u/21031067?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/Codertocat",
		"html_url": "https://github.com/Codertocat",
		"followers_url": "https://api.github.com/users/Codertocat/followers",
		"following_url": "https://api.github.com/users/Codertocat/following{/other_user}",
		"gists_url": "https://api.github.com/users/Codertocat/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/Codertocat/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/Codertocat/subscriptions",
		"organizations_url": "https://api.github.com/users/Codertocat/orgs",
		"repos_url": "https://api.github.com/users/Codertocat/repos",
		"events_url": "https://api.github.com/users/Codertocat/events{/privacy}",
		"received_events_url": "https://api.github.com/users/Codertocat/received_events",
		"type": "User",
		"site_admin": false
	  },
	  "html_url": "https://github.com/Codertocat/Hello-World",
	  "description": null,
	  "fork": false,
	  "url": "https://github.com/Codertocat/Hello-World",
	  "forks_url": "https://api.github.com/repos/Codertocat/Hello-World/forks",
	  "keys_url": "https://api.github.com/repos/Codertocat/Hello-World/keys{/key_id}",
	  "collaborators_url": "https://api.github.com/repos/Codertocat/Hello-World/collaborators{/collaborator}",
	  "teams_url": "https://api.github.com/repos/Codertocat/Hello-World/teams",
	  "hooks_url": "https://api.github.com/repos/Codertocat/Hello-World/hooks",
	  "issue_events_url": "https://api.github.com/repos/Codertocat/Hello-World/issues/events{/number}",
	  "events_url": "https://api.github.com/repos/Codertocat/Hello-World/events",
	  "assignees_url": "https://api.github.com/repos/Codertocat/Hello-World/assignees{/user}",
	  "branches_url": "https://api.github.com/repos/Codertocat/Hello-World/branches{/branch}",
	  "tags_url": "https://api.github.com/repos/Codertocat/Hello-World/tags",
	  "blobs_url": "https://api.github.com/repos/Codertocat/Hello-World/git/blobs{/sha}",
	  "git_tags_url": "https://api.github.com/repos/Codertocat/Hello-World/git/tags{/sha}",
	  "git_refs_url": "https://api.github.com/repos/Codertocat/Hello-World/git/refs{/sha}",
	  "trees_url": "https://api.github.com/repos/Codertocat/Hello-World/git/trees{/sha}",
	  "statuses_url": "https://api.github.com/repos/Codertocat/Hello-World/statuses/{sha}",
	  "languages_url": "https://api.github.com/repos/Codertocat/Hello-World/languages",
	  "stargazers_url": "https://api.github.com/repos/Codertocat/Hello-World/stargazers",
	  "contributors_url": "https://api.github.com/repos/Codertocat/Hello-World/contributors",
	  "subscribers_url": "https://api.github.com/repos/Codertocat/Hello-World/subscribers",
	  "subscription_url": "https://api.github.com/repos/Codertocat/Hello-World/subscription",
	  "commits_url": "https://api.github.com/repos/Codertocat/Hello-World/commits{/sha}",
	  "git_commits_url": "https://api.github.com/repos/Codertocat/Hello-World/git/commits{/sha}",
	  "comments_url": "https://api.github.com/repos/Codertocat/Hello-World/comments{/number}",
	  "issue_comment_url": "https://api.github.com/repos/Codertocat/Hello-World/issues/comments{/number}",
	  "contents_url": "https://api.github.com/repos/Codertocat/Hello-World/contents/{+path}",
	  "compare_url": "https://api.github.com/repos/Codertocat/Hello-World/compare/{base}...{head}",
	  "merges_url": "https://api.github.com/repos/Codertocat/Hello-World/merges",
	  "archive_url": "https://api.github.com/repos/Codertocat/Hello-World/{archive_format}{/ref}",
	  "downloads_url": "https://api.github.com/repos/Codertocat/Hello-World/downloads",
	  "issues_url": "https://api.github.com/repos/Codertocat/Hello-World/issues{/number}",
	  "pulls_url": "https://api.github.com/repos/Codertocat/Hello-World/pulls{/number}",
	  "milestones_url": "https://api.github.com/repos/Codertocat/Hello-World/milestones{/number}",
	  "notifications_url": "https://api.github.com/repos/Codertocat/Hello-World/notifications{?since,all,participating}",
	  "labels_url": "https://api.github.com/repos/Codertocat/Hello-World/labels{/name}",
	  "releases_url": "https://api.github.com/repos/Codertocat/Hello-World/releases{/id}",
	  "deployments_url": "https://api.github.com/repos/Codertocat/Hello-World/deployments",
	  "created_at": 1557933565,
	  "updated_at": "2019-05-15T15:20:41Z",
	  "pushed_at": 1557933657,
	  "git_url": "git://github.com/Codertocat/Hello-World.git",
	  "ssh_url": "git@github.com:Codertocat/Hello-World.git",
	  "clone_url": "https://github.com/Codertocat/Hello-World.git",
	  "svn_url": "https://github.com/Codertocat/Hello-World",
	  "homepage": null,
	  "size": 0,
	  "stargazers_count": 0,
	  "watchers_count": 0,
	  "language": "Ruby",
	  "has_issues": true,
	  "has_projects": true,
	  "has_downloads": true,
	  "has_wiki": true,
	  "has_pages": true,
	  "forks_count": 1,
	  "mirror_url": null,
	  "archived": false,
	  "disabled": false,
	  "open_issues_count": 2,
	  "license": null,
	  "forks": 1,
	  "open_issues": 2,
	  "watchers": 0,
	  "default_branch": "master",
	  "stargazers": 0,
	  "master_branch": "master"
	},
	"pusher": {
	  "name": "Codertocat",
	  "email": "21031067+Codertocat@users.noreply.github.com"
	},
	"sender": {
	  "login": "Codertocat",
	  "id": 21031067,
	  "node_id": "MDQ6VXNlcjIxMDMxMDY3",
	  "avatar_url": "https://avatars1.githubusercontent.com/u/21031067?v=4",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/Codertocat",
	  "html_url": "https://github.com/Codertocat",
	  "followers_url": "https://api.github.com/users/Codertocat/followers",
	  "following_url": "https://api.github.com/users/Codertocat/following{/other_user}",
	  "gists_url": "https://api.github.com/users/Codertocat/gists{/gist_id}",
	  "starred_url": "https://api.github.com/users/Codertocat/starred{/owner}{/repo}",
	  "subscriptions_url": "https://api.github.com/users/Codertocat/subscriptions",
	  "organizations_url": "https://api.github.com/users/Codertocat/orgs",
	  "repos_url": "https://api.github.com/users/Codertocat/repos",
	  "events_url": "https://api.github.com/users/Codertocat/events{/privacy}",
	  "received_events_url": "https://api.github.com/users/Codertocat/received_events",
	  "type": "User",
	  "site_admin": false
	},
	"installation": {
		"id": 8
	}
}
`

func TestPushActionBasic(t *testing.T) {
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := DefaultMockClientProvider()

	commentCreated := false
	expectedMessage := `Added a new version for "Codertocat/Hello-World": "v5"`
	var message string

	mockClientProviderPtr.overrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		commentCreated = true
		log.Println(req.Body)
		json := getBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})

	PushAction(&p, mockClientProviderPtr)

	if !commentCreated {
		t.Errorf("Comment wasn't created\n")
	}
	if message != expectedMessage {
		t.Errorf("Wrong commit comment! expected: %s, got: %s\n", expectedMessage, message)
	}

}
func TestConfiguredPushAction(t *testing.T) {
	var withConfiguredPath = "path: projectA/pom.xml"

	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := DefaultMockClientProvider()

	mockClientProviderPtr.overrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		return newTestResponse(200, mockContentResponse(withConfiguredPath))
	})

	var receivedUrl string

	mockClientProviderPtr.overrideResponseFn("GET_NEW_VERSION", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		receivedUrl = req.URL.String()
		return defaultFn(req)
	})

	PushAction(&p, mockClientProviderPtr)

	if !strings.Contains(receivedUrl, "projectA/pom.xml") {
		t.Errorf("Config is not used\n")
	}
}
func TestMissedOldVersion(t *testing.T) {
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := DefaultMockClientProvider()

	commentCreated := false

	mockClientProviderPtr.overrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		//func not go there
		commentCreated = true
		return defaultFn(req)
	})

	mockClientProviderPtr.overrideResponseFn("GET_OLD_VERSION", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		return newTestResponse(404, "not found")
	})

	PushAction(&p, mockClientProviderPtr)

	if commentCreated {
		t.Errorf("Comment wasn't created\n")
	}
}
func TestMissedNewVersion(t *testing.T) {
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := DefaultMockClientProvider()

	commentCreated := false

	mockClientProviderPtr.overrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		commentCreated = true
		return defaultFn(req)
	})

	mockClientProviderPtr.overrideResponseFn("GET_NEW_VERSION", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		return newTestResponse(404, "not found")
	})

	PushAction(&p, mockClientProviderPtr)

	if commentCreated {
		t.Errorf("Comment should not be created\n")
	}
}
func TestConfiguredTagPrefix(t *testing.T) {
	var withConfiguredPath = "prefix: n"

	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := DefaultMockClientProvider()

	mockClientProviderPtr.overrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		return newTestResponse(200, mockContentResponse(withConfiguredPath))
	})
	expectedMessage := `Added a new version for "Codertocat/Hello-World": "v5"`
	expectedTag := "v5"
	var message string
	var tag string

	mockClientProviderPtr.overrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		json := getBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})
	mockClientProviderPtr.overrideResponseFn("ADD_TAG", func(req *http.Request, defaultFn RoundTripFunc) *http.Response {
		json := getBodyJson(req)
		tag = fmt.Sprintf("%v", json["tag"])
		return defaultFn(req)
	})

	PushAction(&p, mockClientProviderPtr)

	if tag != expectedTag {
		t.Errorf("Wrong tag! expected: %s, got: %s\n", expectedTag, tag)
	}

	if message != expectedMessage {
		t.Errorf("Wrong commit comment! expected: %s, got: %s\n", expectedMessage, message)
	}

}

func TestMadeСaptionToTemplate(t *testing.T) {
	var tests = []struct {
		template string
		version  string
		result   string
	}{
		{`v{{.version}}`, `1.0`, `v1.0`},
		{`vNN{{.version}}`, `1.0`, `vNN1.0`},
		{`v_{{.version}}`, `1.0`, `v_1.0`},
		{`v{{.version}}`, `1.0-relise`, `v1.0-relise`},
		{`v{{.versio}}`, `1.0`, `v1.0`},
		{`{{.version}}`, `1.0`, `1.0`},
		{``, `1.0`, `v1.0`},
	}
	for _, test := range tests {
		result := madeСaptionToTemplate(test.template, test.version)
		if result != test.result {
			t.Errorf("template: %q, version: %q\nwant: %q, got: %q", test.template, test.version, result, test.result)
		}
	}
}
