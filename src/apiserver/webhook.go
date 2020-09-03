package apiserver

import "net/http"

func (api *ActApiServer) webhook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook was received!"))
}