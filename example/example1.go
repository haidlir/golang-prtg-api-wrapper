package main

import (
	"log"

	"github.com/haidlir/golang-prtg-api-wrapper/prtg-api"
)

func main() {
	// Configuration
	server := "https://prtg.paessler.com"
	username := "demo"
	password := "demodemo"

	client := prtg.NewClient(server, username, password)
	prtgVersion, err := client.GetPrtgVersion()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(prtgVersion)
	}
}