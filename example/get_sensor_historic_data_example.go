package main

import (
	"fmt"
	"log"
	"time"

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
	log.Println("Sensor Detail:")
	fmt.Printf(". Sensor Id: %v\n", sensorId)
	fmt.Printf(". Sensor Name: %v\n", sensorDetail.Name)
	fmt.Printf(". Sensor Type: %v\n", sensorDetail.SensorType)
	fmt.Printf(". Sensor's Parent Device: %v\n", sensorDetail.ParentDeviceName)
	fmt.Printf(". Sensor Status: %v\n", sensorDetail.StatusText)
	fmt.Printf(". Sensor Interval: %v\n", sensorDetail.Interval)

	startDate := time.Date(2018, time.May, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, time.May, 1, 0, 35, 0, 0, time.UTC)
	var average int64 = 0
	histData, err := client.GetHistoricData(sensorId, average, startDate, endDate)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("History Data Record from %v to %v:", startDate, endDate)
	fmt.Printf(". Total Data: %v\n", len(histData))
	for _, data := range(histData) {
		fmt.Printf("%v - %vms\n", data["datetime"], data["Loading time"])
	}
}