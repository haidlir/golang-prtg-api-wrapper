package prtg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getHTTPBody(url string, timeout int64) ([]byte, error) {
	// Skipping TLS Verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create GET method: %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(timeout)*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to create HTTP request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Response status NOK: %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read response body: %v", err)
	}
	return body, nil
}

func getSensorDetail(url string, timeout int64) (*PrtgSensorDetailsResponse, error) {
	body, err := getHTTPBody(url, timeout)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var msg PrtgSensorDetailsResponse
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal json response: %v", err)
	}

	return &msg, nil

}

func getHistoricData(url string, timeout int64) (*PrtgHistoricDataResponse, error) {
	body, err := getHTTPBody(url, timeout)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var msg PrtgHistoricDataResponse
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal json response: %v", err)
	}

	return &msg, nil

}