package apiserver

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"githubservice"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Webhook struct {
	Action string `json:"action"`
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
		if strings.HasPrefix(push.GetRef(), "refs/heads/") {
			wh := &Webhook{}
			if err := json.Unmarshal(body, wh); err != nil {
				log.Printf("webhook json.Unmarshal Error: %v", err)
				http.Error(w, "can't parse a webhook payload for getting of installation id", http.StatusInternalServerError)
				return
			}
			go githubservice.PushAction(push, wh.Installation.Id)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("This webhook is undefined yet."))
}