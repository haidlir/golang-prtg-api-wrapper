package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/haidlir/golang-prtg-api-wrapper/prtg-api"
)

func printChild(i interface{}, level int) {
	switch v := i.(type) {
	case prtg.SensorTreeGroup:
		groupInstance := v
		for _, group := range groupInstance.Groups {
			fmt.Printf("%v%vGroup: %v (%v)\n", "|", strings.Repeat("--", level), group.GroupName, group.GroupId)
			printChild(group, level+1)
		}
		for _, probe := range groupInstance.ProbeNodes {
			fmt.Printf("%v%vProbe: %v (%v)\n", "|", strings.Repeat("--", level), probe.ProbeName, probe.ProbeId)
			printChild(probe, level+1)
		}
		for _, device := range groupInstance.Devices {
			fmt.Printf("%v%vDevice: %v (%v)(IpAddr: %v)\n", "|", strings.Repeat("--", level), device.DeviceName, device.DeviceId, device.DeviceHost)
			printChild(device, level+1)
		}
		for _, sensor := range groupInstance.Sensors {
			fmt.Printf("%v%vsensor: %v (%v)\n", "|", strings.Repeat("--", level), sensor.SensorName, sensor.SensorId)
		}
	case prtg.SensorTreeProbeNode:
		probeInstance := v
		for _, group := range probeInstance.Groups {
			fmt.Printf("%v%vGroup: %v (%v)\n", "|", strings.Repeat("--", level), group.GroupName, group.GroupId)
			printChild(group, level+1)
		}
		for _, device := range probeInstance.Devices {
			fmt.Printf("%v%vDevice: %v (%v)\n", "|", strings.Repeat("--", level), device.DeviceName, device.DeviceId)
			printChild(device, level+1)
		}
		for _, sensor := range probeInstance.Sensors {
			fmt.Printf("%v%vsensor: %v (%v)\n", "|", strings.Repeat("--", level), sensor.SensorName, sensor.SensorId)
		}
	case prtg.SensorTreeDevice:
		deviceInstance := v
		for _, sensor := range deviceInstance.Sensors {
			fmt.Printf("%v%vsensor: %v (%v)\n", "|", strings.Repeat("--", level), sensor.SensorName, sensor.SensorId)
		}
	case *prtg.PrtgSensorTreeResponse:
		sensorTree := v
		for _, group := range sensorTree.Groups {
			fmt.Printf("%v%vGroup: %v (%v)\n", "", strings.Repeat("--", level), group.GroupName, group.GroupId)
			printChild(group, level+1)
		}
		for _, probe := range sensorTree.ProbeNodes {
			fmt.Printf("%v%vProbe: %v (%v)\n", "|", strings.Repeat("--", level), probe.ProbeName, probe.ProbeId)
			printChild(probe, level+1)
		}
		for _, device := range sensorTree.Devices {
			fmt.Printf("%v%vDevice: %v (%v)\n", "|", strings.Repeat("--", level), device.DeviceName, device.DeviceId)
			printChild(device, level+1)
		}
		for _, sensor := range sensorTree.Sensors {
			fmt.Printf("%v%vsensor: %v (%v)\n", "|", strings.Repeat("--", level), sensor.SensorName, sensor.SensorId)
		}
	default:
		fmt.Printf("%T\n", v)
	}

}

func main() {
	// Configuration
	server := "https://prtg.paessler.com"
	username := "demo"
	password := "demodemo"
	client := prtg.NewClient(server, username, password)

	var sensorId int64 = 0 // Device Object ID
	sensorTree, err := client.GetSensorTree(sensorId)
	if err != nil {
		log.Println(err)
		return
	}
	printChild(sensorTree, 0)
}
