package apiserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AtcApiServer struct {
	router *mux.Router
}

func Instance() *AtcApiServer {
	return &AtcApiServer{
		mux.NewRouter().StrictSlash(true),
	}
}

func (api *AtcApiServer) Start(host string) {
	api.router.HandleFunc("/api/webhook", api.webhook).Methods("POST")

	if host != "" {
		log.Println("Listening HTTP for", host)
		log.Fatal(http.ListenAndServe(host, api.router))
	}
	log.Print("ATC API Server didn't run!")
}
