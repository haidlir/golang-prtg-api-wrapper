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
	client =  NewClient(server, username, password)
	if client.Server != "http://127.0.0.1" {
		t.Errorf("Server is %v instead of http://127.0.0.1", client.Server)
	}
}

func TestSetContextTimeout(t *testing.T) {
	server := "http://localhost"
	username := "user"
	password := "pass"
	client := NewClient(server, username, password)

	// Check whether client contains default context timeout or not.
	if client.Timeout != 10000 {
		t.Errorf("client's context timeout is %vms instead of 10s", client.Timeout)
	}

	// Trying to change the client context timeout more than or equals 30000
	client.SetContextTimeout(30001)
	if client.Timeout != 10000 {
		t.Errorf("client's context timeout is %vms instead of 10s", client.Timeout)
	}
	// Trying to change the client context timeout less than or equals to 30000
	client.SetContextTimeout(-1)
	if client.Timeout != 10000 {
		t.Errorf("client's context timeout is %vms instead of 10s", client.Timeout)
	}
	// Trying to change the client context timeout more than 30000
	client.SetContextTimeout(30000)
	if client.Timeout != 30000 {
		t.Errorf("client's context timeout is %vms instead of 30000ms", client.Timeout)
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
			t.Errorf("It Should be error when server %v", client.Server)
		}

		sensorId, average, sDate, eDate := composeDummyHistAPIParam()
		_, err = client.GetHistoricData(sensorId, average, sDate, eDate)
		if err == nil {
			t.Errorf("It Should be error when server %v", client.Server)
		}

		_, err = client.GetSensorList(sensorId, nil)
		if err == nil {
			t.Errorf("It Should be error when server %v", client.Server)
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
		t.Errorf("Since the id is less than zero, an error should occur.")
	}

	// Average should be more than or equals to zero
	sensorId = 0
	average = -1
	histData, err = client.GetHistoricData(sensorId, average, sDate, eDate)
	if err == nil {
		t.Errorf("Since the average is less than zero, an error should occur.")
	}
}

func TestGetSensorList(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc(GetTableListsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		sensorId := r.FormValue("id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if sensorId == "9301" {
			fmt.Fprint(w, loadfixture("/prtg_sensor-list_9301.json"))
		} else if sensorId == "9000" {
			fmt.Fprint(w, loadfixture("/prtg_sensor-list_9000_empty.json"))
		} else if sensorId == "9321" {
			fmt.Fprint(w, loadfixture("/prtg_histdata_9321.xml"))
		}
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	server := fmt.Sprintf("%v", serverURL)
	username := "user"
	password := "pass"
	var sensorId int64
	var columns []string
	client := NewClient(server, username, password)

	// Check sensor list within id 9301
	sensorId = 9301
	columns = []string{"objid","probe","group","device","sensor","status","message",
												"lastvalue","priority","favorite"}
	sensorList, err := client.GetSensorList(sensorId, columns)
	if err != nil {
		t.Errorf("It should be success but error: %v", err)
	}
	if len(sensorList) <= 0 {
		t.Errorf("It should be not empty.")
	}

	// Check sensor list within id 9000 (empty)
	sensorId = 9000
	sensorList, err = client.GetSensorList(sensorId, columns)
	if len(sensorList) > 0 {
		t.Errorf("It should be empty.")
	}

	// Check sensor id less than zero
	sensorId = -1
	sensorList, err = client.GetSensorList(sensorId, columns)
	if err == nil {
		t.Errorf("Since the id is less than zero, an error should occur.")
	}

	// Check columns is nil
	sensorId = 9301
	columns = nil
	sensorList, err = client.GetSensorList(sensorId, columns)
	if err != nil {
		t.Errorf("Columns should turn to default column if the columns are nil, but error: %v", err)
	}
	// Check columns contain so many random string
	columns = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	sensorList, err = client.GetSensorList(sensorId, columns)
	if err != nil {
		t.Errorf("Columns should turn to default column if the column's values are too much, but error: %v", err)
	}
	// for sensor id 9321
	sensorId = 9321
	sensorList, err = client.GetSensorList(sensorId, columns)
	if err == nil {
		t.Errorf("Since the response's body is XML, an error should occur.")
	}
}

func TestGetDeviceList(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc(GetTableListsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		sensorId := r.FormValue("id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if sensorId == "9217" {
			fmt.Fprint(w, loadfixture("/prtg_device-list_9217.json"))
		} else if sensorId == "9000" {
			fmt.Fprint(w, loadfixture("/prtg_device-list_9000_empty.json"))
		} else if sensorId == "9321" {
			fmt.Fprint(w, loadfixture("/prtg_histdata_9321.xml"))
		}
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	server := fmt.Sprintf("%v", serverURL)
	username := "user"
	password := "pass"
	var sensorId int64
	var columns []string
	client := NewClient(server, username, password)

	// Check sensor list within id 9301
	sensorId = 9217
	columns = []string{"objid","probe","group","device","host","downsens","partialdownsens",
						"downacksens","upsens","warnsens","pausedsens","unusualsens",
						"undefinedsens"}
	deviceList, err := client.GetDeviceList(sensorId, columns)
	if err != nil {
		t.Errorf("It should be success but error: %v", err)
	}
	if len(deviceList) <= 0 {
		t.Errorf("It should be not empty.")
	}

	// Check sensor list within id 9000 (empty)
	sensorId = 9000
	deviceList, err = client.GetDeviceList(sensorId, columns)
	if len(deviceList) > 0 {
		t.Errorf("It should be empty.")
	}

	// Check sensor id less than zero
	sensorId = -1
	deviceList, err = client.GetDeviceList(sensorId, columns)
	if err == nil {
		t.Errorf("Since the id is less than zero, an error should occur.")
	}

	// Check columns is nil
	sensorId = 9217
	columns = nil
	deviceList, err = client.GetDeviceList(sensorId, columns)
	if err != nil {
		t.Errorf("Columns should turn to default column if the columns are nil, but error: %v", err)
	}
	// Check columns contain so many random string
	columns = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s",
						"t", "u", "v", "w", "x", "y", "z"}
	deviceList, err = client.GetDeviceList(sensorId, columns)
	if err != nil {
		t.Errorf("Columns should turn to default column if the column's values are too much, but error: %v", err)
	}
	// for sensor id 9321
	sensorId = 9321
	deviceList, err = client.GetDeviceList(sensorId, columns)
	if err == nil {
		t.Errorf("Since the response's body is XML, an error should occur.")
	}
}

func TestGetGroupList(t *testing.T) {
	mux := new(http.ServeMux)
	mux.HandleFunc(GetTableListsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		sensorId := r.FormValue("id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if sensorId == "0" {
			fmt.Fprint(w, loadfixture("/prtg_group-list_0.json"))
		} else if sensorId == "9000" {
			fmt.Fprint(w, loadfixture("/prtg_group-list_9000_empty.json"))
		} else if sensorId == "9321" {
			fmt.Fprint(w, loadfixture("/prtg_histdata_9321.xml"))
		}
	})
	httpServer := setup(mux)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	server := fmt.Sprintf("%v", serverURL)
	username := "user"
	password := "pass"
	var sensorId int64
	var columns []string
	client := NewClient(server, username, password)

	// Check sensor list within id 9301
	sensorId = 0
	columns = []string{"objid","probe","group","name","downsens","partialdownsens","downacksens",
						"upsens","warnsens","pausedsens","unusualsens","undefinedsens"}
	groupList, err := client.GetGroupList(sensorId, columns)
	if err != nil {
		t.Errorf("It should be success but error: %v", err)
	}
	if len(groupList) <= 0 {
		t.Errorf("It should be not empty.")
	}

	// Check sensor list within id 9000 (empty)
	sensorId = 9000
	groupList, err = client.GetGroupList(sensorId, columns)
	if len(groupList) > 0 {
		t.Errorf("It should be empty.")
	}

	// Check sensor id less than zero
	sensorId = -1
	groupList, err = client.GetGroupList(sensorId, columns)
	if err == nil {
		t.Errorf("Since the id is less than zero, an error should occur.")
	}

	// Check columns is nil
	sensorId = 0
	columns = nil
	groupList, err = client.GetGroupList(sensorId, columns)
	if err != nil {
		t.Errorf("Columns should turn to default column if the columns are nil, but error: %v", err)
	}
	// Check columns contain so many random string
	columns = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s",
						"t", "u", "v", "w", "x", "y", "z"}
	groupList, err = client.GetGroupList(sensorId, columns)
	if err != nil {
		t.Errorf("Columns should turn to default column if the column's values are too much, but error: %v", err)
	}
	// for sensor id 9321
	sensorId = 9321
	groupList, err = client.GetGroupList(sensorId, columns)
	if err == nil {
		t.Errorf("Since the response's body is XML, an error should occur.")
	}
}