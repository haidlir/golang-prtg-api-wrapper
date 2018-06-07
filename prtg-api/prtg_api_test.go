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
	"time"
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

func composeDummyHistAPIParam() (sensorId, average int64, sDate, eDate time.Time) {
	sensorId = 14254
	average = 0
	sDate = time.Date(2018, time.May, 1, 0, 0, 0, 0, time.UTC)
	eDate = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.UTC)
	return
}

func TestGetCompleteUrl(t *testing.T) {
	servers := []string{" http://localhost", "localhost"}

	for _, server := range(servers) {
		client := NewClient(server, "", "")

		_, err := client.GetPrtgVersion()
		if err == nil {
			t.Errorf("It Should be error when server %v", client.server)
		}

		sensorId, average, sDate, eDate := composeDummyHistAPIParam()
		_, err = client.GetHistoricData(sensorId, average, sDate, eDate)
		if err == nil {
			t.Errorf("It Should be error when server %v", client.server)
			return
		}
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
		t.Errorf("PRTG Version is %v instead of 18.2.41.1636", prtgVersion)
	}
}

func TestGetSensorDetail(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc(GetSensorDetailsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		sensorId := r.FormValue("id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if sensorId == "9182" {
			fmt.Fprint(w, loadfixture("/prtg_sensor_9182.json"))
		} else if sensorId == "9321" {
			fmt.Fprint(w, loadfixture("/prtg_sensor_9321.json"))
		} else if sensorId == "1337" {
			time.Sleep(2 * time.Millisecond)
			fmt.Fprint(w, "")
		}
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	server := fmt.Sprintf("%v", serverURL)
	username := "user"
	password := "pass"
	var sensorId int64
	client := NewClient(server, username, password)

	// for sensor id 9182
	sensorId = 9182
	sensorDetail, err := client.GetSensorDetail(sensorId)
	if err != nil {
		t.Errorf("Unable to get PRTG's Sensor Detail: %v", err)
		return
	}
	if sensorDetail.Name != "NetFlow V5 1" {
		t.Errorf("Sensor's name %v instead of NetFlow V5 1", sensorDetail.Name)
	}

	// for sensor id 9321
	sensorId = 9321
	sensorDetail, err = client.GetSensorDetail(sensorId)
	if err != nil {
		t.Errorf("Unable to get PRTG's Sensor Detail: %v", err)
		return
	}
	if sensorDetail.Name != "Ping" {
		t.Errorf("Sensor's name %v instead of Ping", sensorDetail.Name)
	}

	// for sensor id 1337
	sensorId = 1337
	client.SetContextTimeout(1)
	sensorDetail, err = client.GetSensorDetail(sensorId)
	if err == nil {
		t.Errorf("Since context's timeout reached, error should occur")
		return
	}
}

func TestHistData(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc(GetHistoricDatasEndpoint, func(w http.ResponseWriter, r *http.Request) {
		sensorId := r.FormValue("id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if sensorId == "14254" {
			fmt.Fprint(w, loadfixture("/prtg_histdata_14254.json"))
		} else if sensorId == "9321" {
			fmt.Fprint(w, loadfixture("/prtg_histdata_9321.xml"))
		} else if sensorId == "9000" {
			fmt.Fprint(w, loadfixture("/prtg_histdata_9000_empty.json"))
		}
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	server := fmt.Sprintf("%v", serverURL)
	username := "user"
	password := "pass"
	var sensorId int64
	var average int64 = 0
	sDate := time.Date(2018, time.May, 1, 0, 0, 0, 0, time.UTC)
	eDate := time.Date(2018, time.June, 1, 0, 0, 0, 0, time.UTC)
	client := NewClient(server, username, password)

	// for sensor id 14254
	sensorId = 14254
	histData, err := client.GetHistoricData(sensorId, average, sDate, eDate)
	if err != nil {
		t.Errorf("Unable to get PRTG's Historic Data: %v", err)
		return
	}
	if len(histData) <= 0 {
		t.Errorf("No data within historic data")
	}

	// for sensor id 9000
	sensorId = 9000
	histData, err = client.GetHistoricData(sensorId, average, sDate, eDate)
	if err == nil {
		t.Errorf("Since no historic data found, an error should occur.")
	}

	// for sensor id 9321
	sensorId = 9321
	histData, err = client.GetHistoricData(sensorId, average, sDate, eDate)
	if err == nil {
		t.Errorf("Since the response's body is XML, an error should occur.")
	}


	// Should return error, if data range is more than 31 days
	sDate = time.Date(2018, time.May, 1, 0, 0, 0, 0, time.UTC)
	eDate = time.Date(2018, time.June, 1, 0, 0, 1, 0, time.UTC)
	histData, err = client.GetHistoricData(sensorId, average, sDate, eDate)
	if err == nil {
		t.Errorf("Since the date range is more than 31 days, an error should occur.")
	}

	// id should be more than or equals to zero
	sensorId = -1
	histData, err = client.GetHistoricData(sensorId, average, sDate, eDate)
	if err == nil {
		t.Errorf("Since the daverage is less than zero, an error should occur.")
	}

	// Average should be more than or equals to zero
	sensorId = 0
	average = -1
	histData, err = client.GetHistoricData(sensorId, average, sDate, eDate)
	if err == nil {
		t.Errorf("Since the daverage is less than zero, an error should occur.")
	}
}