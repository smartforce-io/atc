package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-github/v39/github"
)

const (
	tokenResponse = `
{
	"token" : "aaa",
	"expires_at": "2016-07-11T22:14:10Z"
}
`
	oldPomXml = `
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
	<version>4</version>
</project>
`
	newPomXml = `
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
	<version>5</version>
</project>
`
	oldGradle = `
android {
	defaultConfig {
		versionCode 1
		versionName "4"
		}
	}
}
`
	newGradle = `
android {
	defaultConfig {
		versionCode 1
		versionName "5"
		}
	}
}
`
	oldNpm = `
{"version": "4",
"name": "atc"}
`
	newNpm = `
{"version": "5",
"name": "atc"}
`
	oldFlutter = `
{version: 4,
name: atc}
`
	newFlutter = `
{version: 5,
name: atc}
`
	oldUserConf = `
name: testUserConf
vers: 4
`
	newUserConf = `
name: testUserConf
vers: 5
`
)

var (
	ErrUnmarshal = errors.New("unmarshal error")
	ErrGeneral   = errors.New("weird error")
)

func MockContentResponse(content string) string {
	response := fmt.Sprintf(`{"Content" : "%s", "size": %d, "encoding":"base64"}`, base64.StdEncoding.EncodeToString([]byte(content)), len(content))
	return response
}

type MockContentProvider struct {
	Content string
	Err     error
}

func (mockContentProvider *MockContentProvider) GetContents(path string) (string, error) {
	return mockContentProvider.Content, mockContentProvider.Err
}

// RoundTripFunc
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { //this is kind of wrapper where original function is used in interface implementation
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
func NewTestResponse(status int, response string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBufferString(response)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
}
func GetBodyJson(req *http.Request) map[string]interface{} {
	result := make(map[string]interface{})

	reqBytes, _ := io.ReadAll(req.Body)

	json.Unmarshal(reqBytes, &result)

	return result
}

type evaluation struct {
	conditionFn func(req *http.Request) bool
	responseFn  RoundTripFunc
}

type MockClientProvider struct {
	evaluations map[string]evaluation
}

func (mockClientProvider *MockClientProvider) OverrideResponseFn(action string, ovveride func(req *http.Request, defaultFn RoundTripFunc) *http.Response) {
	eval := mockClientProvider.evaluations[action]
	saved := eval.responseFn
	eval.responseFn = func(req *http.Request) *http.Response {
		return ovveride(req, saved)
	}
	mockClientProvider.evaluations[action] = eval
}

func DefaultMockClientProvider() *MockClientProvider {
	defaultEvaluations := map[string]evaluation{
		"GET_TOKEN": {
			func(req *http.Request) bool {
				return strings.HasPrefix(req.URL.String(), "https://api.github.com/app/installations/")
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(201, tokenResponse)
			},
		},
		"GET_ATC_CONFIG": {
			func(req *http.Request) bool {
				return strings.HasSuffix(req.URL.String(), "atc.yaml")
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(404, "not found")
			},
		},
		"GET_OLD_VERSION_MAVEN": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("pom\\.xml\\?ref=[^main]+", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(oldPomXml))
			},
		},
		"GET_NEW_VERSION_MAVEN": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("pom\\.xml\\?ref=main", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(newPomXml))
			},
		},
		"GET_OLD_VERSION_GRADLE": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("build\\.gradle\\?ref=[^main]+", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(oldGradle))
			},
		},
		"GET_NEW_VERSION_GRADLE": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("build\\.gradle\\?ref=main", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(newGradle))
			},
		},
		"GET_OLD_VERSION_NPM": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("package\\.json\\?ref=[^main]+", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(oldNpm))
			},
		},
		"GET_NEW_VERSION_NPM": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("package\\.json\\?ref=main", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(newNpm))
			},
		},
		"GET_OLD_VERSION_FLUTTER": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("pubspec\\.yaml\\?ref=[^main]+", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(oldFlutter))
			},
		},
		"GET_NEW_VERSION_FLUTTER": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("pubspec\\.yaml\\?ref=main", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(newFlutter))
			},
		},
		"GET_OLD_VERSION_USERCONF": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("test\\.txt\\?ref=[^main]+", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(oldUserConf))
			},
		},
		"GET_NEW_VERSION_USERCONF": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString("test\\.txt\\?ref=main", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(200, MockContentResponse(newUserConf))
			},
		},
		"ADD_TAG": {
			func(req *http.Request) bool {
				return strings.Contains(req.URL.String(), "/git/tags")
			},
			func(req *http.Request) *http.Response {
				log.Println("ADD_TAG req.Body: ", req.Body)
				jsonMap := GetBodyJson(req)

				return NewTestResponse(201, fmt.Sprintf(`{"tag":"%s", "sha":"940bd336248efae0f9ee5bc7b2d5c985887b16ac"}`, jsonMap["tag"]))
			},
		},
		"ADD_REF": {
			func(req *http.Request) bool {
				return strings.Contains(req.URL.String(), "/git/refs")
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(201, `{}`)
			},
		},
		"ADD_COMMENT": {
			func(req *http.Request) bool {
				matched, err := regexp.MatchString(".*/commits/(.{40})/comments", req.URL.String())
				if err != nil {
					return false
				}
				return matched
			},
			func(req *http.Request) *http.Response {
				return NewTestResponse(201, `{}`)
			},
		},
	}
	return &MockClientProvider{defaultEvaluations}
}

func (mockClientProvider *MockClientProvider) Get(token string, ctx context.Context) *github.Client {
	client := NewTestClient(func(req *http.Request) *http.Response {
		for _, eval := range mockClientProvider.evaluations {
			if eval.conditionFn(req) {
				log.Println("req.URL: ", req.URL.String())
				return eval.responseFn(req)
			}
		}
		return NewTestResponse(404, "not found")
	})

	return github.NewClient(client)
}
