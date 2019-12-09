package prtg

type prtgSensorDetailsResponse struct {
	PrtgVersion string         `json:"prtgversion" xml:"prtg-version"`
	SensorData  PrtgSensorData `json:"sensordata"`
}

type prtgSensorDetailsResponseXML struct {
	PrtgVersion string `xml:"prtg-version"`
	PrtgSensorData
}

// PrtgSensorData contains property for each sensor, device, and group object within detail API.
type PrtgSensorData struct {
	Name             string `json:"name" xml:"name"`
	SensorType       string `json:"sensortype" xml:"sensortype"`
	Interval         string `json:"interval" xml:"interval"`
	ProbeName        string `json:"probename" xml:"probename"`
	ParentGroupName  string `json:"parentgroupname" xml:"parentgroupname"`
	ParentDeviceName string `json:"parentdevicename" xml:"parentdevicename"`
	ParentDeviceId   string `json:"parentdeviceid" xml:"parentdeviceid"`
	LastValue        string `json:"lastvalue" xml:"lastvalue"`
	LastMessage      string `json:"lastmessage" xml:"lastmessage"`
	Favorite         string `json:"favorite" xml:"favorite"`
	StatusText       string `json:"statustext" xml:"statustext"`
	StatusId         string `json:"statusid" xml:"statusid"`
	LastUp           string `json:"lastup" xml:"lastup"`
	LastDown         string `json:"lastdown" xml:"lastdown"`
	LastCheck        string `json:"lastcheck" xml:"lastcheck"`
	Uptime           string `json:"uptime" xml:"uptime"`
	UptimeTime       string `json:"uptimetime" xml:"uptimetime"`
	Downtime         string `json:"downtime" xml:"downtime"`
	DowntimeTime     string `json:"downtimetime" xml:"downtimetime"`
	UpDownTotal      string `json:"updowntotal" xml:"updowntotal"`
	UpDownSince      string `json:"updownsince" xml:"updownsince"`
	Info             string `json:"info" xml:"info"`
}

type prtgTableListResponse struct {
	PrtgVersion string          `json:"prtgversion" xml:"prtg-version"`
	TreeSize    int64           `json:"treesize" xml:"treesize"`
	Groups      []PrtgTableList `json:"groups, xml:"groups,omitempty"`
	Devices     []PrtgTableList `json:"devices, xml:"devices,omitempty"`
	Sensors     []PrtgTableList `json:"sensors, xml:"sensors,omitempty"`
}

// PrtgTableList contains property for each sensor, device, and group object within list API.
type PrtgTableList struct {
	ObjectId           int64  `json:"objid" xml:"objid"`
	Probe              string `json:"probe" xml:"probe"`
	Group              string `json:"group" xml:"group"`
	Name               string `json:"name" xml:"name"`
	Device             string `json:"device" xml:"device"`
	Host               string `json:"host" xml:"host"`
	Sensor             string `json:"sensor" xml:"sensor"`
	DownSensors        int64  `json:"downsens_raw" xml:"downsens_raw"`
	PartialDownSensors int64  `json:"partialdownsens_raw" xml:"partialdownsens_raw"`
	DownAckSensors     int64  `json:"downacksens_raw" xml:"downacksens_raw"`
	UpSensors          int64  `json:"upsens_raw" xml:"upsens_raw"`
	WarningSensors     int64  `json:"warnsens_raw" xml:"warnsens_raw"`
	PausedSensors      int64  `json:"pausedsens_raw" xml:"pausedsens_raw"`
	UnusualSensors     int64  `json:"unusualsens_raw" xml:"unusualsens_raw"`
	UndefinedSensors   int64  `json:"undefinedsens_raw" xml:"undefinedsens_raw"`
}

type prtgHistoricDataResponse struct {
	PrtgVersion  string             `json:"prtgversion" xml:"prtg-version"`
	TreeSize     int64              `json:"treesize" xml:"treesize"`
	HistoricData []PrtgHistoricData `json:"histdata" xml:"histdata"`
}

// PrtgHistoricData contains historic data param and value for each series.
type PrtgHistoricData map[string]interface{}

type prtgHistoricDataResponseXML struct {
	PrtgVersion  string       `json:"prtgversion" xml:"prtg-version"`
	HistoricData []ItemTagXML `json:"histdata" xml:"item"`
}

// ItemTagXML contains the tag informations
type ItemTagXML struct {
	Datetime    string     `xml:"datetime"`
	DatetimeRAW string     `xml:"datetime_raw"`
	Coverage    string     `xml:"coverage"`
	CoverageRAW string     `xml:"coverage_raw"`
	Value       []ValueXML `xml:"value"XML`
	ValueRAW    []ValueXML `xml:"value_raw"`
}

type ValueXML struct {
	Key   string `xml:"channel,attr"`
	Value string `xml:",chardata"`
}

// PrtgSensorTreeResponse contains parsed xml format of sensor tree API response.
type PrtgSensorTreeResponse struct {
	PrtgVersion string                `xml:"prtg-version"`
	Groups      []SensorTreeGroup     `xml:"sensortree>nodes>group"`
	ProbeNodes  []SensorTreeProbeNode `xml:"sensortree>nodes>probenode"`
	Devices     []SensorTreeDevice    `xml:"sensortree>nodes>device"`
	Sensors     []SensorTreeSensor    `xml:"sensortree>nodes>sensor"`
}

// SensorTreeGroup contains Group's Tree structure.
type SensorTreeGroup struct {
	GroupId     int64                 `xml:"id"`
	GroupName   string                `xml:"name"`
	GroupTags   string                `xml:"tags"`
	GroupActive bool                  `xml:"active"`
	Groups      []SensorTreeGroup     `xml:"group"`
	ProbeNodes  []SensorTreeProbeNode `xml:"probenode"`
	Devices     []SensorTreeDevice    `xml:"device"`
	Sensors     []SensorTreeSensor    `xml:"sensor"`
}

// SensorTreeProbeNode contains Probe's Tree structure.
type SensorTreeProbeNode struct {
	ProbeId       int64              `xml:"id,attr"`
	ProbeName     string             `xml:"name"`
	ProbeNoAccess int64              `xml:"noaccess,attr"`
	Groups        []SensorTreeGroup  `xml:"group"`
	Devices       []SensorTreeDevice `xml:"device"`
	Sensors       []SensorTreeSensor `xml:"sensor"`
}

// SensorTreeDevice contains Device's Tree structure.
type SensorTreeDevice struct {
	DeviceId     int64              `xml:"id"`
	DeviceName   string             `xml:"name"`
	DeviceTags   string             `xml:"tags"`
	DeviceHost   string             `xml:"host"`
	DeviceActive bool               `xml:"active"`
	Sensors      []SensorTreeSensor `xml:"sensor"`
}

// SensorTreeSensor contains Sensor's Tree structure.
type SensorTreeSensor struct {
	SensorId                int64   `xml:"id"`
	SensorName              string  `xml:"name"`
	SensorTags              string  `xml:"tags"`
	SensorType              string  `xml:"sensortype"`
	SensorKind              string  `xml:"sensorkind"`
	SensorInterval          int64   `xml:"interval"`
	SensorStatus            string  `xml:"status"`
	SensorLastValue         float64 `xml:"lastvalue_raw"`
	SensorStatusMessage     string  `xml:"statusmessage"`
	SensorStatusSince       float64 `xml:"statussince_raw_utc"`
	SensorLastTime          float64 `xml:"lasttime_raw_utc"`
	SensorLastOk            float64 `xml:"lastok_raw_utc"`
	SensorLastError         float64 `xml:"lasterror_raw_utc"`
	SensorLastUp            float64 `xml:"lastup_raw_utc"`
	SensorLastDown          float64 `xml:"lastdown_raw_utc"`
	SensorCumulatedDownTime float64 `xml:"cumulateddowntime_raw"`
	SensorCumulatedUpTime   float64 `xml:"cumulateduptime_raw"`
	SensorCumulatedSince    float64 `xml:"cumulatedsince_raw"`
	SensorActive            bool    `xml:"active"`
}
