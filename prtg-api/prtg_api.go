package prtg

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Client contains PRTG's server data, and http client parameters.
type Client struct {
	// PRTG's Server URL (mandatory)
	// It should be in the form of http://host:port
	Server string

	// Any account's username of PRTG (mandatory)
	Username string

	// Any account's password of PRTG (mandatory)
	Password string

	// Any account's password hash of PRTG
	PasswordHash string

	// Timeout Context in millisecond
	Timeout int64
}

var (
	defaultTimeout         int64    = 10000
	deltaHistoricThreshold int64    = 31 * 24 * 60 * 60 // 31 days
	dateFormat             string   = "2006-01-02-15-04-05"
	defaultSensorListCols  []string = []string{"objid", "probe", "group", "device", "sensor", "status", "message",
		"lastvalue", "priority", "favorite"}
	defaultSensorListColsLen int      = len(defaultSensorListCols)
	defaultDeviceListCols    []string = []string{"objid", "probe", "group", "device", "host", "downsens", "partialdownsens",
		"downacksens", "upsens", "warnsens", "pausedsens", "unusualsens",
		"undefinedsens"}
	defaultDeviceListColsLen int      = len(defaultDeviceListCols)
	defaultGroupListCols     []string = []string{"objid", "probe", "group", "name", "downsens", "partialdownsens", "downacksens",
		"upsens", "warnsens", "pausedsens", "unusualsens", "undefinedsens"}
	defaultGroupListColsLen int = len(defaultGroupListCols)
)

const (
	// GetSensorDetailsEndpoint contains path to sensor detail API endpoint
	GetSensorDetailsEndpoint = "/api/getsensordetails.json"
	// GetTableListsEndpoint contains path to table list API endpoint for listing groups, devices, and sensors.
	GetTableListsEndpoint = "/api/table.json"
	// GetHistoricDatasEndpoint contains path to historic data API endpoint
	GetHistoricDatasEndpoint = "/api/historicdata.json"
	// GetSensorTreesEndpoint contains path to serson tree API endpoint
	GetSensorTreesEndpoint = "/api/table.xml"
	// Some Private constant.
	userAgent = "golang-prtg-api"
)

// NewClient takes server, username, password and returns client's instance.
// input format:
// server := "http://localhost"
// username := "user"
// password := "pass"
func NewClient(server, username, password string) *Client {
	instance := new(Client)
	instance.Server = server
	instance.Username = username
	instance.Password = password
	instance.Timeout = 10000
	return instance
}

// NewClientWithHashedPass takes server, username, passwordHash and returns client's instance.
// input format:
// server := "http://localhost"
// username := "user"
// passwordHash := "000000000"
func NewClientWithHashedPass(server, username, passwordHash string) *Client {
	instance := new(Client)
	instance.Server = server
	instance.Username = username
	instance.PasswordHash = passwordHash
	instance.Timeout = 10000
	return instance
}

// SetContextTimeout configures the client timeout value in millisecond format.
func (c *Client) SetContextTimeout(timeout int64) {
	if (timeout <= 0) || (timeout > 30000) {
		c.Timeout = defaultTimeout
	} else {
		c.Timeout = timeout
	}
}

func (c *Client) getTemplateUrlQuery() *url.Values {
	q := url.Values{}
	q.Set("username", c.Username)
	if c.Password != "" {
		q.Set("password", c.Password)
	}
	if c.PasswordHash != "" {
		q.Set("passhash", c.PasswordHash)
	}
	return &q
}

func (c *Client) getCompleteUrl(p string, q *url.Values) (string, error) {
	u, err := url.Parse(c.Server)
	if err != nil {
		return "", fmt.Errorf("Unable to parse url: %v", err)
	}
	u.Path = p
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *Client) getSensorDetail(q *url.Values) (*prtgSensorDetailsResponse, error) {
	p := GetSensorDetailsEndpoint

	// Complete URL
	u, err := c.getCompleteUrl(p, q)
	if err != nil {
		return nil, err
	}

	var sensorDetailResp prtgSensorDetailsResponse
	err = getPrtgResponse(u, c.Timeout, &sensorDetailResp)
	if err != nil {
		return nil, err
	}
	return &sensorDetailResp, nil
}

// GetPrtgVersion returns PRTG's version of the specified server.
func (c *Client) GetPrtgVersion() (string, error) {
	// Set the query
	q := c.getTemplateUrlQuery()
	q.Set("id", "0")

	sensorDetailResp, err := c.getSensorDetail(q)
	if err != nil {
		return "", err
	}
	return sensorDetailResp.PrtgVersion, nil
}

// GetSensorDetail returns the detail of specified sensor.
func (c *Client) GetSensorDetail(id int64) (*PrtgSensorData, error) {
	// Set the query
	q := c.getTemplateUrlQuery()
	q.Set("id", fmt.Sprintf("%v", id))

	sensorDetailResp, err := c.getSensorDetail(q)
	if err != nil {
		return nil, err
	}
	return &sensorDetailResp.SensorData, nil
}

func (c *Client) getHistoricData(id, average int64, startDate, endDate time.Time) (*prtgHistoricDataResponse, error) {
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

	var histDataResp prtgHistoricDataResponse
	err = getPrtgResponse(u, c.Timeout, &histDataResp)
	if err != nil {
		return nil, err
	}
	return &histDataResp, nil
}

func getDeltaSecond(sDate, eDate time.Time) int64 {
	return eDate.Unix() - sDate.Unix()
}

// GetHistoricData returns series of recorded data of specified sensor.
// Take start and end of date's boundaries.
func (c *Client) GetHistoricData(id, average int64, startDate, endDate time.Time) ([]PrtgHistoricData, error) {
	// Validate Input
	// Make sure that id and average is not less than 0
	if id < 0 || average < 0 {
		return nil, fmt.Errorf("Id and average should be more than or equals to zero")
	}
	// Make sure that data range less than 31 days
	if deltaSecond := getDeltaSecond(startDate, endDate); (deltaSecond < 0) || (deltaSecond > deltaHistoricThreshold) {
		return nil, fmt.Errorf("Data range is more than 31 days")
	}

	// Get Historic Data using PRTG's API
	histDataResp, err := c.getHistoricData(id, average, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("Unable to get historic data: %v", err)
	}
	if len(histDataResp.HistoricData) <= 0 {
		return histDataResp.HistoricData, fmt.Errorf("No Data Found")
	}

	// Return the historic data
	return histDataResp.HistoricData, nil
}

func (c *Client) getTableList(id int64, content string, columns []string) (*prtgTableListResponse, error) {
	// Compose queries
	q := c.getTemplateUrlQuery()
	q.Set("id", fmt.Sprintf("%v", id))
	q.Set("content", fmt.Sprintf("%v", content))
	colStr := strings.Join(columns, ",")
	q.Set("columns", fmt.Sprintf("%v", colStr))
	p := GetTableListsEndpoint
	// Complete URL
	u, err := c.getCompleteUrl(p, q)
	if err != nil {
		return nil, err
	}

	var tableListResp prtgTableListResponse
	err = getPrtgResponse(u, c.Timeout, &tableListResp)
	if err != nil {
		return nil, err
	}
	return &tableListResp, nil
}

// GetSensorList returns list of sensor within specified device.
// The id retrive sensor and group's object id, but it will return empty list.
// The default columns's value is nil.
func (c *Client) GetSensorList(id int64, columns []string) ([]PrtgTableList, error) {
	// Validate input
	// Make sure that id is not less than 0
	if id < 0 {
		return nil, fmt.Errorf("Id should be more than or equals to zero")
	}
	// if columns is nil, use the default column's entry instead
	if (columns == nil) || (len(columns) > defaultSensorListColsLen) {
		columns = defaultSensorListCols
	}

	// Get sensor list within this group or device
	content := "sensors"
	sensorListResp, err := c.getTableList(id, content, columns)
	if err != nil {
		return nil, fmt.Errorf("Unable to get sensor list data: %v", err)
	}
	if len(sensorListResp.Sensors) <= 0 {
		return sensorListResp.Sensors, fmt.Errorf("No Data Found")
	}

	// Return sensor list
	return sensorListResp.Sensors, nil
}

// GetDeviceList returns list of sensor within specified device.
// The id retrive sensor and device's object id, but it will return empty list.
// The default columns's value is nil.
func (c *Client) GetDeviceList(id int64, columns []string) ([]PrtgTableList, error) {
	// Validate input
	// Make sure that id is not less than 0
	if id < 0 {
		return nil, fmt.Errorf("Id should be more than or equals to zero")
	}
	// if columns is nil, use the default column's entry instead
	if (columns == nil) || (len(columns) > defaultDeviceListColsLen) {
		columns = defaultDeviceListCols
	}

	// Get sensor list within this group or device
	content := "devices"
	sensorListResp, err := c.getTableList(id, content, columns)
	if err != nil {
		return nil, fmt.Errorf("Unable to get device list data: %v", err)
	}
	if len(sensorListResp.Devices) <= 0 {
		return sensorListResp.Devices, fmt.Errorf("No Data Found")
	}

	// Return sensor list
	return sensorListResp.Devices, nil
}

// GetGroupList returns list of sensor within specified device.
// The id retrive sensor and device's object id, but it will return empty list.
// The default columns's value is nil.
// Since in PRTG, it's possible to have nested group
func (c *Client) GetGroupList(id int64, columns []string) ([]PrtgTableList, error) {
	// Validate input
	// Make sure that id is not less than 0
	if id < 0 {
		return nil, fmt.Errorf("Id should be more than or equals to zero")
	}
	// if columns is nil, use the default column's entry instead
	if (columns == nil) || (len(columns) > defaultGroupListColsLen) {
		columns = defaultGroupListCols
	}

	// Get sensor list within this group or device
	content := "groups"
	sensorListResp, err := c.getTableList(id, content, columns)
	if err != nil {
		return nil, fmt.Errorf("Unable to get group list data: %v", err)
	}
	if len(sensorListResp.Groups) <= 0 {
		return sensorListResp.Groups, fmt.Errorf("No Data Found")
	}

	// Return sensor list
	return sensorListResp.Groups, nil
}

func (c *Client) getTableTree(id int64) (*PrtgSensorTreeResponse, error) {
	// Compose queries
	q := c.getTemplateUrlQuery()
	q.Set("id", fmt.Sprintf("%v", id))
	q.Set("content", "sensortree")
	p := GetSensorTreesEndpoint
	// Complete URL
	u, err := c.getCompleteUrl(p, q)
	if err != nil {
		return nil, err
	}

	var tableTreeResp PrtgSensorTreeResponse
	err = getPrtgResponse(u, c.Timeout, &tableTreeResp)
	if err != nil {
		return nil, err
	}
	return &tableTreeResp, nil
}

// GetSensorTree return tree structures within the specified object id.
// It's possible to capture the whole relation of sensors, devices, and groups
// in tree format, instead of getting the information separately using GetSensorList,
// GetDeviceList, and GetGroupList.
// If id is not zero, it will capture the sensortree from specific group or device.
func (c *Client) GetSensorTree(id int64) (*PrtgSensorTreeResponse, error) {
	// Validate input
	// Make sure that id is not less than 0
	if id < 0 {
		return nil, fmt.Errorf("Id should be more than or equals to zero")
	}

	sensorTree, err := c.getTableTree(id)
	if err != nil {
		return nil, fmt.Errorf("Unable to get sensor tree data: %v", err)
	}
	return sensorTree, nil
}
