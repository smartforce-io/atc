package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

func (api *ActApiServer) webhook(w http.ResponseWriter, r *http.Request) {
	switch event := r.Header.Get("X-GitHub-Event"); event {
	case "marketplace_purchase":
		w.WriteHeader(http.StatusOK)
		body, _ := ioutil.ReadAll(r.Body)
		log.Printf("markeplace purchase event: \n %s \n", body)

	case "create":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))

	case "delete":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))

	case "push":
		body, _ := ioutil.ReadAll(r.Body)

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
