package main

import (
	"log"

	"github.com/haidlir/golang-prtg-api-wrapper/prtg-api"
)

func main() {
	// Configuration
	server := "https://prtg.paessler.com"
	username := "demo"
	passwordHash := "000000000"

	client := prtg.NewClientWithHashedPass(server, username, passwordHash)
	prtgVersion, err := client.GetPrtgVersion()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("The version of PRTG on %v is %v.", server, prtgVersion)
	}
}
