package main

import (
	"fmt"
	"log"

	"github.com/haidlir/golang-prtg-api-wrapper/prtg-api"
)

func main() {
	// Configuration
	server := "https://prtg.paessler.com"
	username := "demo"
	password := "demodemo"
	client := prtg.NewClient(server, username, password)

	var sensorId int64 = 9217 // Group Object ID
	sensorDetail, err := client.GetSensorDetail(sensorId)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Group Detail:")
	fmt.Printf(". Group Id: %v\n", sensorId)
	fmt.Printf(". Group Name: %v\n", sensorDetail.Name)
	fmt.Printf(". Group Type: %v\n", sensorDetail.SensorType)
	fmt.Printf(". Group's Parent Device: %v\n", sensorDetail.ParentDeviceName)
	fmt.Printf(". Group Status: %v\n", sensorDetail.StatusText)

	deviceList, err := client.GetDeviceList(sensorId, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("List of device under sensor %v:", sensorId)
	fmt.Printf(". Total device: %v\n", len(deviceList))
	for _, device := range deviceList {
		fmt.Printf(". %v - %v - %v\n", device.ObjectId, device.Device, device.Host)
	}
}
