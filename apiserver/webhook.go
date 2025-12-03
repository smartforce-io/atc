package apiserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-github/v39/github"

	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/push"
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

		p := &github.WebHookPayload{}
		if err := json.Unmarshal(body, p); err != nil {
			log.Printf("webhook json.Unmarshal Error: %v", err)
			http.Error(w, "can't parse a webhook payload", http.StatusInternalServerError)
			return
		}
		if p.Installation == nil || p.Installation.ID == nil {
			log.Printf("p webhook doesn't contain installation info: %v", p)
			http.Error(w, "p webhook doesn't contain installation info", http.StatusBadRequest)
			return
		}
		if strings.HasPrefix(p.GetRef(), "refs/heads/") {
			go push.ActionPush(p, &provider.GithubClientProvider{}) //it's not clear who is resposible for DI
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
