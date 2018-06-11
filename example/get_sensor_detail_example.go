package main

import (
	"log"

	"github.com/haidlir/golang-prtg-api-wrapper/prtg-api"
)

func main() {
	// Configuration
	server := "https://prtg.paessler.com"
	username := "demo"
	password := "demodemo"
	client := prtg.NewClient(server, username, password)

	var sensorId int64 = 7986 // Sensor Homepage
	sensorDetail, err := client.GetSensorDetail(sensorId)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Sensor Id: %v", sensorId)
	log.Printf("Sensor Name: %v", sensorDetail.Name)
	log.Printf("Sensor Type: %v", sensorDetail.SensorType)
	log.Printf("Sensor's Parent Device: %v", sensorDetail.ParentDeviceName)
	log.Printf("Sensor Status: %v", sensorDetail.StatusText)

}