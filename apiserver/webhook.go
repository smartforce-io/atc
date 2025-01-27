package apiserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/smartforce-io/atc/githubservice"

	"github.com/google/go-github/v39/github"
)

type Webhook struct {
	Action       string       `json:"action"`
	Installation Installation `json:"installation"`
}

type Installation struct {
	Id int64 `json:"id"`
}

func (api *AtcApiServer) webhook(w http.ResponseWriter, r *http.Request) {
	switch event := r.Header.Get("X-GitHub-Event"); event {
	case "marketplace_purchase":
		w.WriteHeader(http.StatusOK)
		body, _ := io.ReadAll(r.Body)
		log.Printf("markeplace purchase event: \n %s \n", body)

	case "create":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))

	case "delete":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))

	case "push":
		body, _ := io.ReadAll(r.Body)
		body = removeOrgFromWebhookRequest(body)

		push := &github.WebHookPayload{}
		if err := json.Unmarshal(body, push); err != nil {
			log.Printf("webhook json.Unmarshal Error: %v", err)
			http.Error(w, "can't parse a webhook payload", http.StatusInternalServerError)
			return
		}
		if push.Installation == nil || push.Installation.ID == nil {
			log.Printf("push webhook doesn't contain installation info: %v", push)
			http.Error(w, "push webhook doesn't contain installation info", http.StatusBadRequest)
			return
		}
		if strings.HasPrefix(push.GetRef(), "refs/heads/") {
			go githubservice.PushAction(push, &githubservice.GithubClientProvider{}) //it's not clear who is resposible for DI
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This webhook is undefined yet."))
	}
}

func removeOrgFromWebhookRequest(body []byte) []byte {
	reg, err := regexp.Compile(`,"organization":"[^\t\n\f\r\"]+"`)
	if err != nil {
		log.Printf("err compile regexp: %v", err)
	}
	return []byte(reg.ReplaceAllString(string(body), ""))
}
