package prtg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"os"
	"path/filepath"
)

func TestNewClient(t *testing.T) {
	server := "http://localhost"
	username := "user"
	password := "pass"

	// Trying to create new client
	client := NewClient(server, username, password)
	if client == nil {
		t.Error("A new connection object must have been made")
	}

	// Trying to change the server
	server = "http://127.0.0.1"
	_ =  NewClient(server, username, password)
	if client.server != "http://127.0.0.1" {
		t.Errorf("Server is %v instead of 127.0.0.1", client.server)
	}
}

func TestSetContextTimeout(t *testing.T) {
	server := "http://localhost"
	username := "user"
	password := "pass"
	client := NewClient(server, username, password)

	// Check whether client contains default context timeout or not.
	if client.timeout != 10000 {
		t.Errorf("client's context timeout is %vms instead of 10s", client.timeout)
	}

	// Trying to change the client context timeout more than or equals 30000
	client.SetContextTimeout(30001)
	if client.timeout != 10000 {
		t.Errorf("client's context timeout is %vms instead of 10s", client.timeout)
	}
	// Trying to change the client context timeout less than or equals to 30000
	client.SetContextTimeout(-1)
	if client.timeout != 10000 {
		t.Errorf("client's context timeout is %vms instead of 10s", client.timeout)
	}
	// Trying to change the client context timeout more than 30000
	client.SetContextTimeout(30000)
	if client.timeout != 30000 {
		t.Errorf("client's context timeout is %vms instead of 30000ms", client.timeout)
	}
}

func TestGetCompleteUrl(t *testing.T) {
	server := " http://localhost"
	client := NewClient(server, "", "")

	_, err := client.GetPrtgVersion()
	if err == nil {
		t.Errorf("It Should be error when server %v", client.server)
	}
}

func TestIncompleteUrl(t *testing.T) {
	server := "localhost"
	client := NewClient(server, "", "")

	_, err := client.GetPrtgVersion()
	if err == nil {
		t.Errorf("It Should be error when server %v", client.server)
		return
	}
}

// Inspired by go-octokit
// setup sets up a test HTTP server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(mux *http.ServeMux) *httptest.Server {
	// test server
	server := httptest.NewServer(mux)
	return server
}

func loadfixture(f string) string {
	pwd, _ := os.Getwd()
	p := filepath.Join(pwd, "..", "fixtures", f)
	c, _ := ioutil.ReadFile(p)
	return string(c)
}

func TestGetPrtgVersion(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc(GetSensorDetailsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, loadfixture("/prtg_version.json"))
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	server := fmt.Sprintf("%v", serverURL)
	username := "user"
	password := "pass"
	client := NewClient(server, username, password)
	prtgVersion, err := client.GetPrtgVersion()
	if err != nil {
		t.Errorf("Unable to get PRTG Version: %v", err)
		return
	}
	if prtgVersion != "18.2.41.1636" {
		t.Errorf("PRTG Version is %v instead of 18.2.41.1636", client.server)
	}
}