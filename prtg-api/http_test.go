package prtg

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestGetHttpBody(t *testing.T) {
	var url string
	var timeout int64
	var err error

	// Wrong written url
	url = " http://localhost"
	timeout = 10000
	_, err = getHTTPBody(url, timeout)
	if err == nil {
		t.Errorf("It Should be error (at NewRequest()) if url %v", url)
	}

	// When server not found or inactive
	url = "http://localhost"
	_, err = getHTTPBody(url, timeout)
	if err == nil {
		t.Errorf("It Should be error (at Send Request) if server down: %v", err)
	}
}

func TestRespStatusCode(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc(GetSensorDetailsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, loadfixture("/prtg_version.json"))
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	path := "wrong/path"
	u := fmt.Sprintf("%v/%v", serverURL, path)
	var timeout int64 = 10000
	_, err := getHTTPBody(u, timeout)
	if err == nil {
		t.Errorf("%v", err)
	}
}

func TestJsonUnmarshalGetSensorDetail(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, loadfixture("/prtg_version_error.json"))
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	u := fmt.Sprintf("%v/", serverURL)
	var timeout int64 = 10000
	result, err := getSensorDetail(u, timeout)
	if err == nil {
		t.Errorf("%v", result)
	}
}