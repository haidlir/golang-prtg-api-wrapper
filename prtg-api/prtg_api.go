package prtg

import (
	"fmt"
	"net/url"
	"time"
)

// Client's fields are read-only onece instantiated.
// So it's safe to use it in concurrent condition.
type client struct {
	// PRTG's Server URL (mandatory)
	// It should be in the form of http://host:port
	server				string

	// Any account's username of PRTG (mandatory)
	username			string

	// Any account's password of PRTG (mandatory)
	password			string

	// Timeout Context in milisecond
	timeout				int64
}

// Client is a client interface for querying PRTG's server.
type Client interface {
	// Get PRTG's version
	GetPrtgVersion() (string, error)

	// Get details of specific PRTG's sensor
	GetSensorDetail(id int64) (*PrtgSensorData, error)
	
	// Get records of data from specific sensor
	// The reponse's format depends on the sensor's type
	GetHistoricData(id, average int64, startDate, endDate *time.Time) ([]map[string]interface{}, error)

	// Get all sensors under specific devices or groups
	// GetSensorList(id int64, columns []string) error

	// Get all devices under specific groups
	// GetDeviceList(id int64, columns []string) error

	// Get all groups under specific groups
	// Since in PRTG, it's possible to have nested group
	// GetGroupList(id int64, columns []string) error

	// It's possible to capture the whole relation of sensors, devices, and groups
	// in tree format, instead of getting the information separately using GetSensorList,
	// GetDeviceList, and GetGroupList.
	// If id is not zero, it will capture the sensortree from specific group or device.
	// GetSensorTree(id int64) error

	// Set Context HTTP Request Timeout
	SetContextTimeout(timeout int64)
}

var instance *client
var (
	defaultTimeout int64 = 10000
	deltaHistoricThreshold int64 = 31 * 24 * 60 * 60 // 31 days
	dateFormat string = "2006-01-02-15-04-05"
)
const (
	GetSensorDetailsEndpoint = "/api/getsensordetails.json"
	GetSensorListsEndpoint = "/api/table.json"
	GetHistoricDatasEndpoint = "/api/historicdata.json"
	GetSensorTreesEndpoint = "/api/table.xml"
	userAgent = "golang-prtg-api"
)

// Create new Client that later used to request data from PRTG's server
func NewClient(server, username, password string) *client {
	if instance == nil {
		instance = new(client)
	}
	instance.server = server
	instance.username = username
	instance.password = password
	instance.timeout = 10000
	return instance
}

// Set Context HTTP Request Timeout in milisecond
func (c *client) SetContextTimeout(timeout int64) {
	if (timeout <= 0) || (timeout >30000) {
		instance.timeout = defaultTimeout
	} else {
		instance.timeout = timeout
	}
}

func (c *client) getTemplateUrlQuery() (*url.Values) {
	q := url.Values{}
	q.Set("username", c.username)
	q.Set("password", c.password)
	return &q
}

func (c *client) getCompleteUrl(p string, q *url.Values) (string, error) {
	u, err := url.Parse(c.server)
	if err != nil {
		return "", fmt.Errorf("Unable to parse url: %v", err)
	}
	u.Path = p
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *client) getSensorDetail(q *url.Values) (*PrtgSensorDetailsResponse, error) {
	p := GetSensorDetailsEndpoint

	// Complete URL
	u, err := c.getCompleteUrl(p, q)
	if err != nil {
		return nil, err
	}

	sensorDetail, err := getSensorDetail(u, c.timeout)
	if err != nil {
		return nil, err
	}
	return sensorDetail, nil
}


// Get PRTG's version.
// Take nothing as input.
// Return PRTG's version in string.
func (c *client) GetPrtgVersion() (string, error) {
	// Set the query
	q := c.getTemplateUrlQuery()
	q.Set("id", "0")

	sensorDetail, err := c.getSensorDetail(q)
	if err != nil {
		return "", err
	}
	return sensorDetail.PrtgVersion, nil
}

// Get details for specific sensor.
// Take sensor's id as input.
// Return sensor structure.
func (c *client) GetSensorDetail(id int64) (*PrtgSensorData, error) {
	// Set the query
	q := c.getTemplateUrlQuery()
	q.Set("id", fmt.Sprintf("%v", id))

	sensorDetail, err := c.getSensorDetail(q)
	if err != nil {
		return nil, err
	}
	return &sensorDetail.SensorData, nil
}

func (c *client) getHistoricData(id, average int64, startDate, endDate time.Time) (*PrtgHistoricDataResponse, error) {
	// Compose queries
	q := c.getTemplateUrlQuery()
	q.Set("id", fmt.Sprintf("%v", id))
	q.Set("avg", fmt.Sprintf("%v", average))
	q.Set("sDate", fmt.Sprintf("%v", startDate.Format(dateFormat)))
	q.Set("eDate", fmt.Sprintf("%v", endDate.Format(dateFormat)))
	q.Set("usecaption", fmt.Sprintf("%v", 1))
	p := GetHistoricDatasEndpoint
	// Complete URL
	u, err := c.getCompleteUrl(p, q)
	if err != nil {
		return nil, err
	}

	histDataResp, err := getHistoricData(u, c.timeout)
	if err != nil {
		return nil, err
	}
	return histDataResp, nil
}

func getDeltaSecond(sDate, eDate time.Time) int64 {
	return eDate.Unix() - sDate.Unix()
}

func (c *client) GetHistoricData(id, average int64, startDate, endDate time.Time) ([]PrtgHistoricData, error) {
	// Validate Input
	// Make sure that average is not less than 0
	if id < 0 || average < 0{
		return nil, fmt.Errorf("Id should be more than or equals to zero")
	}
	// Make sure that data range less than 31 days
	if deltaSecond := getDeltaSecond(startDate, endDate); (deltaSecond < 0) || (deltaSecond > deltaHistoricThreshold) {
		return nil, fmt.Errorf("Data range is more than 31 days")
	}

	// Get Historic Data using PRTG's API
	histData, err := c.getHistoricData(id, average, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("Unable to get historic data: %v", err)
	}
	if len(histData.HistoricData) <= 0 {
		return histData.HistoricData, fmt.Errorf("No Data Found")
	}

	// Return the historic data
	return histData.HistoricData, nil
}
