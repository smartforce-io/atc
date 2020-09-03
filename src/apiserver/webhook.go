package apiserver

import (
	"io/ioutil"
	"log"
	"net/http"
)

func (api *ActApiServer) webhook(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	log.Println("Reveived webhook:", string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook was received!"))
}