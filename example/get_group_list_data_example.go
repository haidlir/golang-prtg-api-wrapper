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

	var sensorId int64 = 0 // Root Group Object ID
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

	groupList, err := client.GetGroupList(sensorId, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("List of group under sensor %v:", sensorId)
	fmt.Printf(". Total group: %v\n", len(groupList))
	for _, group := range(groupList) {
		fmt.Printf(". %v - %v - %v\n", group.ObjectId, group.Group, group.Name)
	}
}