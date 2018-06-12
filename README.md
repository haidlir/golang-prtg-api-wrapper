# golang-prtg-api-wrapper
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

	client := prtg.NewClient(server, username, password)
	prtgVersion, err := client.GetPrtgVersion()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("The version of PRTG on %v is %v.", server, prtgVersion)
	}
}
```
[More Example...](https://github.com/haidlir/golang-prtg-api-wrapper/tree/master/example)

## License
It is released under the MIT license. See
[LICENSE](https://github.com/haidlir/golang-prtg-api-wrapper/blob/master/LICENSE).
