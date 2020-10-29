package apiserver

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"githubservice"
	"io/ioutil"
	"log"
	"net/http"
)

type Webhook struct {
	Action string `json:"action"`
	Installation Installation `json:"installation"`
}

type Installation struct {
	Id int64 `json:"id"`
}

func (api *ActApiServer) webhook(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	wh := &Webhook{}
	if err := json.Unmarshal(body, wh); err != nil {
		log.Printf("webhook json.Unmarshal Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Can't unmarshal the webhook"))
		return
	}

	switch wh.Action {
	case "created":
		w.Write([]byte("Hello from Automated Tag Creator!"))
		w.WriteHeader(http.StatusOK)
		return

	case "deleted":
		w.Write([]byte("Bye-bye! Automated Tag Creator will be waiting you again!"))
		w.WriteHeader(http.StatusOK)
		return
	case "":
		push, err := emptyAction(body)
		if err != nil {
			log.Printf("emptyAction error: %v", err)
		} else {
			go githubservice.PushAction(push, wh.Installation.Id)
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("This webhook is undefined yet."))
}

func emptyAction(b []byte) (*github.WebHookPayload, error)  {
	push := &github.WebHookPayload{}
	if err := json.Unmarshal(b, push); err != nil {
		return nil, err
	}
	return push, nil
}