package push

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/smartforce-io/atc/githubservice/provider"

	"github.com/smartforce-io/atc/envvars"

	"github.com/google/go-github/v39/github"
)

var (
	testWebhookPayload = `
{
	"ref": "refs/heads/main",
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
	  "default_branch": "main",
	  "stargazers": 0,
	  "master_branch": "main"
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
	testRsaKey = `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
)

func TestPushActionBasic(t *testing.T) {
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	commentCreated := false
	expectedMessage := `File .atc.yaml not found or path = "". Used default settings. Added a new version for "Codertocat/Hello-World": "v5"`
	var message string

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		commentCreated = true
		log.Println("req.Body: ", req.Body)
		j := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", j["body"])
		return defaultFn(req)
	})
	ActionPush(&p, mockClientProviderPtr)

	if !commentCreated {
		t.Errorf("Comment wasn't created\n")
	}
	if message != expectedMessage {
		t.Errorf("Wrong commit comment! expected: %s, got: %s\n", expectedMessage, message)
	}
}
func TestConfiguredPushAction(t *testing.T) {
	var testsConfigPath = []struct {
		confString      string
		expectedUrlPath string
		messageError    string
	}{
		{`path: projectA/pom.xml`, `projectA/pom.xml`, `Used default regexStr in file pom.xml. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: projectA/contents/pom.xml`, `projectA/contents/pom.xml`, `Used default regexStr in file pom.xml. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: build.gradle`, `build.gradle`, `Used default regexStr in file build.gradle. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: contents/build.gradle`, `contents/build.gradle`, `Used default regexStr in file build.gradle. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: package.json`, `package.json`, `Used default regexStr in file package.json. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: contents/package.json`, `contents/package.json`, `Used default regexStr in file package.json. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: pubspec.yaml`, `pubspec.yaml`, `Used default regexStr in file pubspec.yaml. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: contents/pubspec.yaml`, `contents/pubspec.yaml`, `Used default regexStr in file pubspec.yaml. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: /projectA/pom.xml`, ``, `error config file .atc.yaml; path has prefix "/"`},
		{`path: contents//build.gradle`, ``, `error config file .atc.yaml; path has "//"`},
		{`path: test.txt`, ``, `Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: `, ``, `File .atc.yaml not found or path = "". Used default settings. Added a new version for "Codertocat/Hello-World": "v5"`},
	}

	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	var receivedUrl string
	var config string
	var message string

	mockClientProviderPtr.OverrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(200, provider.MockContentResponse(config))
	})

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		json := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_MAVEN", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		receivedUrl = req.URL.String()
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_GRADLE", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		receivedUrl = req.URL.String()
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_NPM", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		receivedUrl = req.URL.String()
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_FLUTTER", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		receivedUrl = req.URL.String()
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_USERCONF", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		receivedUrl = req.URL.String()
		return defaultFn(req)
	})

	for _, test := range testsConfigPath {
		config = fmt.Sprintf(`
%s
behavior: before
template: v{{.Version}}
branch: main
regexstr: "vers: (.+)"`, test.confString)
		receivedUrl = ""
		message = ""

		ActionPush(&p, mockClientProviderPtr)

		matched, err := regexp.MatchString(test.expectedUrlPath, receivedUrl)
		if err != nil {
			t.Errorf("regexp error: %s", err)
		}
		if !matched {
			t.Errorf("Config:%s is not used:\nexpectedUrl: %s, receivedUrl: %s", test.confString, test.expectedUrlPath, receivedUrl)
		}
		if message != test.messageError {
			t.Errorf("Error config: %s\nexpectedErrorMessage: %s, got: %s", test.confString, test.messageError, message)
		}
	}
}

func TestMissedOldNewVersionNoConfig(t *testing.T) {
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	commentCreated := false
	expectedMessage := `File .atc.yaml not found or path = "". Not found supported package manager.`
	var message string

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		commentCreated = true
		json := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("GET_OLD_VERSION_MAVEN", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	mockClientProviderPtr.OverrideResponseFn("GET_OLD_VERSION_GRADLE", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	mockClientProviderPtr.OverrideResponseFn("GET_OLD_VERSION_NPM", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	mockClientProviderPtr.OverrideResponseFn("GET_OLD_VERSION_FLUTTER", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	ActionPush(&p, mockClientProviderPtr)

	if !commentCreated {
		t.Errorf("Comment wasn't created\n")
	}

	if message != expectedMessage {
		t.Errorf("Wrong commit comment! expected: %s, got: %s\n", expectedMessage, message)
	}
}

func TestMissedOldVersionWithConfig(t *testing.T) {
	var testsMissVersion = []struct {
		confString                string
		defMockClientPrKeyVersion string
		expectedMessage           string
	}{
		{`
path: projectA/pom.xml
behavior: before
template: MavenV{{.Version}}
branch: main`, "GET_OLD_VERSION_MAVEN", "file pom.xml with old version not found"},
		{`
path: build.gradle
behavior: after
template: GradleV{{.Version}}
branch: main`, "GET_OLD_VERSION_GRADLE", "file build.gradle with old version not found"},
		{`
path: package.json
behavior: before
template: NPMv{{.Version}}
branch: main`, "GET_OLD_VERSION_NPM", "file package.json with old version not found"},
		{`
path: pubspec.yaml
behavior: after
template: FlutterV{{.Version}}
branch: main`, "GET_OLD_VERSION_FLUTTER", "file pubspec.yaml with old version not found"},
		{`
path: test.txt
behavior: after
template: FlutterV{{.Version}}
branch: main`, "GET_OLD_VERSION_USERCONF", ".atc.yaml don't have regexstr for not default package manager file test.txt."},
	}
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	commentCreated := false
	var message string

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		commentCreated = true
		json := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})

	for _, test := range testsMissVersion {
		mockClientProviderPtr.OverrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
			return provider.NewTestResponse(200, provider.MockContentResponse(test.confString))
		})

		mockClientProviderPtr.OverrideResponseFn(test.defMockClientPrKeyVersion, func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
			return provider.NewTestResponse(404, "not found")
		})
		message = ""

		ActionPush(&p, mockClientProviderPtr)

		if !commentCreated {
			t.Errorf("Comment should not be created\n")
		}

		if message != test.expectedMessage {
			t.Errorf("Wrong commit comment! expected: %s, got: %s\n", test.expectedMessage, message)
		}
	}
}

func TestMissedNewVersionNoConfig(t *testing.T) {
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	commentCreated := false
	expectedMessage := `File .atc.yaml not found or path = "". Not found supported package manager.`
	var message string

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		commentCreated = true
		json := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_MAVEN", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_GRADLE", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_NPM", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	mockClientProviderPtr.OverrideResponseFn("GET_NEW_VERSION_FLUTTER", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(404, "not found")
	})

	ActionPush(&p, mockClientProviderPtr)

	if !commentCreated {
		t.Errorf("Comment wasn't created\n")
	}

	if message != expectedMessage {
		t.Errorf("Wrong commit comment! expected: %s, got: %s\n", expectedMessage, message)
	}
}

func TestMissedNewVersionWithConfig(t *testing.T) {
	var testsMissVersion = []struct {
		confString                string
		defMockClientPrKeyVersion string
		expectedMessage           string
	}{
		{`
path: projectA/pom.xml
behavior: before
template: MavenV{{.Version}}
branch: main`, "GET_NEW_VERSION_MAVEN", "file pom.xml with new version not found"},
		{`
path: build.gradle
behavior: after
template: GradleV{{.Version}}
branch: main`, "GET_NEW_VERSION_GRADLE", "file build.gradle with new version not found"},
		{`
path: package.json
behavior: before
template: NPMv{{.Version}}
branch: main`, "GET_NEW_VERSION_NPM", "file package.json with new version not found"},
		{`
path: pubspec.yaml
behavior: after
template: FlutterV{{.Version}}
branch: main`, "GET_NEW_VERSION_FLUTTER", "file pubspec.yaml with new version not found"},
		{`
path: test.txt
behavior: after
template: FlutterV{{.Version}}
branch: main`, "GET_NEW_VERSION_USERCONF", ".atc.yaml don't have regexstr for not default package manager file test.txt."},
	}
	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	commentCreated := false
	var message string

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		commentCreated = true
		j := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", j["body"])
		return defaultFn(req)
	})

	for _, test := range testsMissVersion {
		mockClientProviderPtr.OverrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
			return provider.NewTestResponse(200, provider.MockContentResponse(test.confString))
		})

		mockClientProviderPtr.OverrideResponseFn(test.defMockClientPrKeyVersion, func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
			return provider.NewTestResponse(404, "not found")
		})
		message = ""

		ActionPush(&p, mockClientProviderPtr)

		if !commentCreated {
			t.Errorf("Comment should not be created\n")
		}

		if message != test.expectedMessage {
			t.Errorf("Wrong commit comment! expected: %s, got: %s\n", test.expectedMessage, message)
		}
	}
}

func TestConfiguredTagTemplate(t *testing.T) {
	var testsConfigTemplate = []struct {
		confString      string
		expectedMessage string
		expectedTag     string
	}{
		{`template: v{{.Version}}`, `Added a new version for "Codertocat/Hello-World": "v5"`, `v5`},
		{`template: v{{.Version}}-{{.Version}}`, `Added a new version for "Codertocat/Hello-World": "v5-5"`, `v5-5`},
		{`template: vTest{{.Version}}`, `Added a new version for "Codertocat/Hello-World": "vTest5"`, `vTest5`},
		{`template: "{{.Version}}Vte"`, `Added a new version for "Codertocat/Hello-World": "5Vte"`, `5Vte`},
		{`template: vVv{.Version}`, `error config file .atc.yaml: template doesn't contain "{{.Version}}"`, ``},
		{`template: `, `Added a new version for "Codertocat/Hello-World": "v5"`, `v5`},
	}

	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	var message string
	var tag string
	var config string

	mockClientProviderPtr.OverrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(200, provider.MockContentResponse(config))
	})

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		j := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", j["body"])
		return defaultFn(req)
	})
	mockClientProviderPtr.OverrideResponseFn("ADD_TAG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		j := provider.GetBodyJson(req)
		tag = fmt.Sprintf("%v", j["tag"])
		return defaultFn(req)
	})

	for _, test := range testsConfigTemplate {
		config = fmt.Sprintf(`
path: contents/pom.xml
behavior: before
%s
branch: main`, test.confString)

		message = ""
		tag = ""

		ActionPush(&p, mockClientProviderPtr)

		if tag != test.expectedTag {
			t.Errorf("Wrong tag! expected: %s, got: %s\n", test.expectedTag, tag)
		}

		if message != test.expectedMessage {
			t.Errorf("Wrong commit comment! confString: %s\nexpected: %s, got: %s\n", test.confString, test.expectedMessage, message)
		}
	}
}

func TestConfiguredTagBehavior(t *testing.T) {
	var testsConfigBehavior = []struct {
		confString  string
		expectedSha string
	}{
		{`behavior: after`, `0000000000000000000000000000000000000000`},
		{`behavior: before`, `6113728f27ae82c7b1a177c8d03f9e96e0adf246`},
		{`behavior: bef`, `0000000000000000000000000000000000000000`},
		{`behavior: `, `0000000000000000000000000000000000000000`},
	}

	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	var config string
	var sha string
	var message string
	errorMessage := `error config file .atc.yaml: behavior doesn't contain "before" or "after"`

	mockClientProviderPtr.OverrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(200, provider.MockContentResponse(config))
	})

	mockClientProviderPtr.OverrideResponseFn("ADD_TAG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		j := provider.GetBodyJson(req)
		sha = fmt.Sprintf("%v", j["object"])
		return defaultFn(req)
	})

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		j := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", j["body"])
		return defaultFn(req)
	})

	for _, test := range testsConfigBehavior {
		config = fmt.Sprintf(`
path: contents/pom.xml
%s
template: v{{.Version}}
branch: main`, test.confString)

		sha = ""

		ActionPush(&p, mockClientProviderPtr)

		if sha != test.expectedSha {
			if message != errorMessage {
				t.Errorf("Wrong sha! confString: %s\nexpected sha: %s, got sha: %s\nexpected err: %s, got err: %s", test.confString, test.expectedSha, sha, errorMessage, message)
			}
		}
	}
}

func TestConfiguredBranch(t *testing.T) {
	var testsConfigBehavior = []struct {
		confString      string
		expectedMessage string
	}{
		{`branch: test`, ``},
		{`branch: main`, `Used default regexStr in file pom.xml. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`branch: `, `Used default regexStr in file pom.xml. Added a new version for "Codertocat/Hello-World": "v5"`},
		{``, `Used default regexStr in file pom.xml. Added a new version for "Codertocat/Hello-World": "v5"`},
	}

	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	var config string
	var message string
	errorMessageEmtry := `error config file .atc.yaml; behavior = ""`
	errorMessage := `error config file .atc.yaml: behavior no contains "before" or "after"`

	mockClientProviderPtr.OverrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(200, provider.MockContentResponse(config))
	})

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		json := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})

	for _, test := range testsConfigBehavior {
		config = fmt.Sprintf(`
path: contents/pom.xml
behavior: after
template: v{{.Version}}
%s
regexstr: "vers: (.+)"`, test.confString)

		message = ""

		ActionPush(&p, mockClientProviderPtr)

		if message != test.expectedMessage {
			if message != errorMessage && message != errorMessageEmtry {
				t.Errorf("Wrong branch! confString: %s\nexpected: %s, got: %s\n", test.confString, test.expectedMessage, message)
			}
		}
	}
}

func TestConfiguredRegexStr(t *testing.T) {
	var testsConfigBehavior = []struct {
		confPath        string
		confRegexStr    string
		expectedMessage string
	}{
		{`path: pom.xml`, `regexstr: "vers: (.+)"`, `Used default regexStr in file pom.xml. Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: pom.xml`, `regexstr: `, `Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: pom.xml`, ``, `Added a new version for "Codertocat/Hello-World": "v5"`},
		{`path: test.txt`, ``, `.atc.yaml don't have regexstr for not default package manager file test.txt.`},
		{`path: test.txt`, `regexstr: "vers: (.+)"`, `Added a new version for "Codertocat/Hello-World": "v5"`},
	}

	p := github.WebHookPayload{}
	json.Unmarshal([]byte(testWebhookPayload), &p)

	os.Setenv(envvars.PemData, testRsaKey)

	mockClientProviderPtr := provider.DefaultMockClientProvider()

	var config string
	var message string
	errorMessageEmtry := `error config file .atc.yaml; behavior = ""`
	errorMessage := `error config file .atc.yaml: behavior no contains "before" or "after"`

	mockClientProviderPtr.OverrideResponseFn("GET_ATC_CONFIG", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		return provider.NewTestResponse(200, provider.MockContentResponse(config))
	})

	mockClientProviderPtr.OverrideResponseFn("ADD_COMMENT", func(req *http.Request, defaultFn provider.RoundTripFunc) *http.Response {
		json := provider.GetBodyJson(req)
		message = fmt.Sprintf("%v", json["body"])
		return defaultFn(req)
	})

	for _, test := range testsConfigBehavior {
		config = fmt.Sprintf(`
%s
behavior: after
template: v{{.Version}}
branch: main
%s`, test.confPath, test.confRegexStr)

		message = ""

		ActionPush(&p, mockClientProviderPtr)

		if message != test.expectedMessage {
			if message != errorMessage && message != errorMessageEmtry {
				t.Errorf("Wrong branch! confString: %s\nexpected: %s, got: %s\n", test.confPath, test.expectedMessage, message)
			}
		}
	}
}

func TestMadeСaptionToTemplate(t *testing.T) {
	var tests = []struct {
		template string
		version  string
		result   string
	}{
		{`v{{.Version}}`, `1.0`, `v1.0`},
		{`vNN{{.Version}}`, `1.0`, `vNN1.0`},
		{`v_{{.Version}}`, `1.0`, `v_1.0`},
		{`v{{.Version}}`, `1.0-relise`, `v1.0-relise`},
		{`{{.Version}}`, `1.0`, `1.0`},
		{`Time hour now: {{Time.Hour}}, {{.Version}}`, `1.0`, "Time hour now: " + strconv.Itoa(time.Now().Hour()) + ", 1.0"},
		{``, `1.0`, ``},
	}
	for _, test := range tests {
		result, _ := renderTagNameTemplate(test.template, test.version)
		if result != test.result {
			t.Errorf("template: %q, version: %q\nwant: %q, got: %q", test.template, test.version, test.result, result)
		}
	}
}

func TestMadeСaptionToTemplateError(t *testing.T) {
	var tests = []struct {
		template  string
		version   string
		errString string
	}{
		{`v{{.Versio}}`, `1.0`, `template: template tagContent:1:3: executing "template tagContent" at <.Versio>: can't evaluate field Versio in type push.TagContent`},
	}
	for _, test := range tests {
		_, err := renderTagNameTemplate(test.template, test.version)
		if fmt.Sprint(err) != test.errString {
			t.Errorf("template: %q, version: %q\nerr want: %v, err got: %v", test.template, test.version, test.errString, err)
		}
	}
}
