package main

import (
	"apiserver"
	"log"
	"os"
)

func main() {
	log.Println("Automated Tag Creator")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	apiserver.Instance().Start(":"+port)
}