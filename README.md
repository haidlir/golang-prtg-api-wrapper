# golang-prtg-api-wrapper
[![Build Status](https://travis-ci.org/haidlir/golang-prtg-api-wrapper.svg?branch=master)](https://travis-ci.org/haidlir/golang-prtg-api-wrapper) [![Coverage Status](https://coveralls.io/repos/github/haidlir/golang-prtg-api-wrapper/badge.svg?branch=master)](https://coveralls.io/github/haidlir/golang-prtg-api-wrapper?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/haidlir/golang-prtg-api-wrapper)](https://goreportcard.com/report/github.com/haidlir/golang-prtg-api-wrapper) [![GoDoc](https://sonarcloud.io/api/project_badges/measure?project=golang-prtg-api-wrapper&metric=alert_status)](https://sonarcloud.io/dashboard?id=golang-prtg-api-wrapper) [![GoDoc](https://godoc.org/github.com/haidlir/golang-prtg-api-wrapper/prtg-api?status.svg)](https://godoc.org/github.com/haidlir/golang-prtg-api-wrapper/prtg-api)<br />
PRTG API WRAPPER for Golang Developer

## Status
Experimental

## Motivation
To be used by developers in my team, to fetch data from PRTG.
Freely used by others according to the [LICENSE](https://github.com/haidlir/golang-prtg-api-wrapper/blob/master/LICENSE).

## How to Start
```bash
$ go get github.com/haidlir/golang-prtg-api-wrapper
```

## Example
```go
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
	// or
	passwordHash := "passhash"

	client := prtg.NewClient(server, username, password, passwordHash)
	prtgVersion, err := client.GetPrtgVersion()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("The version of PRTG on %v is %v.", server, prtgVersion)
	}
}
```
[More Example...](https://github.com/haidlir/golang-prtg-api-wrapper/tree/master/_example)

## License
It is released under the MIT license. See
[LICENSE](https://github.com/haidlir/golang-prtg-api-wrapper/blob/master/LICENSE).
