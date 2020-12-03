package apiserver

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ActApiServer struct {
	router *mux.Router
}

func Instance() *ActApiServer {
	return &ActApiServer{
		mux.NewRouter().StrictSlash(true),
	}
}

func (api *ActApiServer) Start(host string) {
	api.router.HandleFunc("/api/webhook", api.webhook).Methods("POST")

	if host != "" {
		log.Println("Listening HTTP for", host)
		log.Fatal(http.ListenAndServe(host, api.router))
	}
	log.Print("ACT API Server didn't run!")
}