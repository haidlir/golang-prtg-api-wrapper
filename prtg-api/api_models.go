package prtg

type prtgSensorDetailsResponse struct {
	PrtgVersion				string				`json:"prtgversion"`
	SensorData				PrtgSensorData		`json:"sensordata"`
}

type PrtgSensorData struct {
	Name					string				`json:"name"`
	SensorType				string				`json:"sensortype"`
	Interval				string				`json:"interval"`
	ProbeName				string				`json:"probename"`
	ParentGroupName			string				`json:"parentgroupname"`
	ParentDeviceName		string				`json:"parentdevicename"`
	ParentDeviceId			string				`json:"parentdeviceid"`
	LastValue				string				`json:"lastvalue"`
	LastMessage				string				`json:"lastmessage"`
	Favorite				string				`json:"favorite"`
	StatusText				string				`json:"statustext"`
	StatusId				string				`json:"statusid"`
	LastUp					string				`json:"lastup"`
	LastDown				string				`json:"lastdown"`
	LastCheck				string				`json:"lastcheck"`
	Uptime					string				`json:"uptime"`
	UptimeTime				string				`json:"uptimetime"`
	Downtime				string				`json:"downtime"`
	DowntimeTime			string				`json:"downtimetime"`
	UpDownTotal				string				`json:"updowntotal"`
	UpDownSince				string				`json:"updownsince"`
	Info					string				`json:"info"`
}

type prtgTableListResponse struct {
	PrtgVersion				string				`json:"prtg-version"`
	TreeSize				int64				`json:"treesize"`
	Groups					[]PrtgTableList		`json:"groups,omitempty"`
	Devices					[]PrtgTableList		`json:"devices,omitempty"`
	Sensors					[]PrtgTableList		`json:"sensors,omitempty"`
}

type PrtgTableList struct {
	ObjectId				int64				`json:"objid"`
	Probe					string				`json:"probe"`
	Group					string				`json:"group"`
	Name					string				`json:"name"`
	DownSensors				int64				`json:"downsens_raw"`
	PartialDownSensors		int64				`json:"partialdownsens_raw"`
	DownAckSensors			int64				`json:"downacksens_raw"`
	UpSensors				int64				`json:"upsens_raw"`
	WarningSensors			int64				`json:"warnsens_raw"`
	PausedSensors			int64				`json:"pausedsens_raw"`
	UnusualSensors			int64				`json:"unusualsens_raw"`
	UndefinedSensors		int64				`json:"undefinedsens_raw"`
}

type prtgHistoricDataResponse struct {
	PrtgVersion				string				`json:"prtg-version"`
	TreeSize				int64				`json:"treesize"`
	HistoricData			[]PrtgHistoricData	`json:"histdata"`
}

type PrtgHistoricData map[string]interface{}

type prtgSensorTreeResponse struct {
	PrtgVersion				string				`xml:"prtg>prtg-version"`
	SensorTreeGroups		[]SensorTreeGroup	`xml:"prtg>sensortree>nodes>group id"`
}

type SensorTreeGroup struct {
	GroupId					int64				`xml:"id"`
	GroupName				string				`xml:"name"`
	GroupTags				string				`xml:"tags"`
	GroupActive				bool				`xml:"active"`
	SensorTreeGroups		[]SensorTreeGroup	`xml:"probenode>group id"`
	SensorTreeDevices		[]SensorTreeDevice	`xml:"name"`
}

type SensorTreeDevice struct {
	DeviceId				int64				`xml:"id"`
	DeviceName				string				`xml:"name"`
	DeviceTags				string				`xml:"tags"`
	DeviceHost				string				`xml:"host"`
	DeviceActive			bool				`xml:"active"`
	SensorTreeSensors		[]SensorTreeSensor	`xml:"sensor"`
}

type SensorTreeSensor struct {
	SensorId				int64				`xml:"id"`
	SensorName				string				`xml:"name"`
	SensorTags				string				`xml:"tags"`
	SensorType				string				`xml:"sensortype"`
	SensorKind				string				`xml:"sensorkind"`
	SensorInterval			int64				`xml:"interval"`
	SensorStatus			string				`xml:"status"`
	SensorLastValue			float64				`xml:"lastvalue_raw"`
	SensorStatusMessage		string				`xml:"statusmessage"`
	SensorStatusSince		float64				`xml:"statussince_raw_utc"`
	SensorLastTime			float64				`xml:"lasttime_raw_utc"`
	SensorLastOk			float64				`xml:"lastok_raw_utc"`
	SensorLastError			float64				`xml:"lasterror_raw_utc"`
	SensorLastUp			float64				`xml:"lastup_raw_utc"`
	SensorLastDown			float64				`xml:"lastdown_raw_utc"`
	SensorCumulatedDownTime	float64				`xml:"cumulateddowntime_raw"`
	SensorCumulatedUpTime	float64				`xml:"cumulateduptime_raw"`
	SensorCumulatedSince	float64				`xml:"cumulatedsince_raw"`
	SensorActive			bool				`xml:"active"`
}

