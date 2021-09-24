package main

import (
	"log"
	"os"

	"github.com/smartforce-io/atc/apiserver"
)

func main() {
	log.Println("Automated Tag Creator")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	apiserver.Instance().Start(":" + port)
}
