package main

import (
	"log"
	"os"

	"github.com/smartforce-io/atc/apiserver"
	"github.com/smartforce-io/atc/githubservice/push"
)

func main() {
	log.Println("Automated Tag Creator")
	mode := os.Getenv("CI_MODE")
	switch {
	case mode != "":
		err := push.CIActionPush()
		if err != nil {
			log.Fatalf("error creating tag %v", err)
		}
	default:
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		apiserver.Instance().Start(":" + port)
	}
}
