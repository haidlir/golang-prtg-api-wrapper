package prtg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func isContentXML(header http.Header) bool {
	contentDisposition := header.Get("Content-Disposition")
	return contentDisposition == "attachment; filename=table.xml"
}

func getHTTPBody(url string, timeout int64) ([]byte, *http.Header, error) {
	// Skipping TLS Verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to create GET method: %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(timeout)*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to create HTTP request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		return nil, nil, fmt.Errorf("Wrong Username and/or Password | HTTP Response status NOK: %v", res.StatusCode)
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("HTTP Response status NOK: %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to read response body: %v", err)
	}
	return body, &res.Header, nil
}

func getPrtgResponse(url string, timeout int64, v interface{}) error {
	body, header, err := getHTTPBody(url, timeout)
	if err != nil {
		return err
	}

	// Unmarshal XML
	if isContentXML(*header) {
		err = xml.Unmarshal(body, &v)
		if err != nil {
			return fmt.Errorf("Unable to unmarshal xml response: %v", err)
		}
		return nil
	}
	// Unmarshal JSON for default
	err = json.Unmarshal(body, &v)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal json response: %v", err)
	}
	return nil
}
