package apiserver

import (
	"encoding/json"
	"githubservice"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
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
	case "created":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Automated Tag Creator!"))
		return

	case "deleted":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bye-bye! Automated Tag Creator will be waiting you again!"))
		return

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
			go githubservice.PushAction(push)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("This webhook is undefined yet."))
}
