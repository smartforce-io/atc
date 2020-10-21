package main

import (
	"apiserver"
	"log"
)

func main() {
	log.Println("Automated Tag Creator")
	apiserver.Instance().Start(":8080")
}