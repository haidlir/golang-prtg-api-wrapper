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

	var sensorId int64 = 9301 // Device Object ID
	sensorDetail, err := client.GetSensorDetail(sensorId)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Device Detail:")
	fmt.Printf(". Device Id: %v\n", sensorId)
	fmt.Printf(". Device Name: %v\n", sensorDetail.Name)
	fmt.Printf(". Device Type: %v\n", sensorDetail.SensorType)
	fmt.Printf(". Device's Parent Device: %v\n", sensorDetail.ParentDeviceName)
	fmt.Printf(". Device Status: %v\n", sensorDetail.StatusText)

	sensorList, err := client.GetSensorList(sensorId, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("List of sensor under sensor %v:", sensorId)
	fmt.Printf(". Total sensor: %v\n", len(sensorList))
	for _, sensor := range sensorList {
		fmt.Printf(". %v - %v - %v\n", sensor.ObjectId, sensor.Sensor, sensor.Probe)
	}
}
