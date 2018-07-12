package main

import (
    "fmt"
    "log"

    "github.com/haidlir/golang-prtg-api-wrapper/prtg-api"
)

type SensorDetail struct {
    GroupLadder   []string
    Device        string
    Sensor        string
}

func convertTreeToFlat(i interface{}) ([]SensorDetail) {
    switch v := i.(type) {
    case prtg.SensorTreeGroup:
        groupInstance := v
        sensorDetails := []SensorDetail{}
        newSensorDetails := []SensorDetail{}
        for _, group := range groupInstance.Groups {
            newSensorDetails = convertTreeToFlat(group)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        for _, probe := range groupInstance.ProbeNodes {
            newSensorDetails = convertTreeToFlat(probe)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        for _, device := range groupInstance.Devices {
            newSensorDetails = convertTreeToFlat(device)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        newSensorDetails = []SensorDetail{}
        for _, sensor := range groupInstance.Sensors {
            newSensorDetails = append(newSensorDetails, SensorDetail {
                Sensor: fmt.Sprintf("%v (%v)", sensor.SensorName, sensor.SensorId),
            })
        }
        sensorDetails = append(sensorDetails, newSensorDetails...)
        for i, _ := range(sensorDetails) {
            sensorDetails[i].GroupLadder = append([]string{fmt.Sprintf("%v (%v)", groupInstance.GroupName, groupInstance.GroupId)},
                sensorDetails[i].GroupLadder...)
        }
        return sensorDetails
    case prtg.SensorTreeProbeNode:
        probeInstance := v
        sensorDetails := []SensorDetail{}
        newSensorDetails := []SensorDetail{}
        for _, group := range probeInstance.Groups {
            newSensorDetails = convertTreeToFlat(group)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        for _, device := range probeInstance.Devices {
            newSensorDetails = convertTreeToFlat(device)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        newSensorDetails = []SensorDetail{}
        for _, sensor := range probeInstance.Sensors {
            newSensorDetails = append(newSensorDetails, SensorDetail {
                Sensor: fmt.Sprintf("%v (%v)", sensor.SensorName, sensor.SensorId),
            })
        }
        sensorDetails = append(sensorDetails, newSensorDetails...)
        for i, _ := range(sensorDetails) {
            sensorDetails[i].GroupLadder = append([]string{fmt.Sprintf("%v (%v)", probeInstance.ProbeName, probeInstance.ProbeId)},
                sensorDetails[i].GroupLadder...)
        }
        return sensorDetails
    case prtg.SensorTreeDevice:
        deviceInstance := v
        sensorDetails := []SensorDetail{}
        newSensorDetails := []SensorDetail{}
        for _, sensor := range deviceInstance.Sensors {
            newSensorDetails = append(newSensorDetails, SensorDetail {
                Sensor: fmt.Sprintf("%v (%v)", sensor.SensorName, sensor.SensorId),
            })
        }
        sensorDetails = append(sensorDetails, newSensorDetails...)
        for i, _ := range(sensorDetails) {
            sensorDetails[i].Device = fmt.Sprintf("%v (%v)", deviceInstance.DeviceName, deviceInstance.DeviceId)
        }
        return sensorDetails
    case *prtg.PrtgSensorTreeResponse:
        sensorTree := v
        sensorDetails := []SensorDetail{}
        newSensorDetails := []SensorDetail{}
        for _, group := range sensorTree.Groups {
            newSensorDetails = convertTreeToFlat(group)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        for _, probe := range sensorTree.ProbeNodes {
            newSensorDetails = convertTreeToFlat(probe)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        for _, device := range sensorTree.Devices {
            newSensorDetails = convertTreeToFlat(device)
            sensorDetails = append(sensorDetails, newSensorDetails...)
        }
        newSensorDetails = []SensorDetail{}
        for _, sensor := range sensorTree.Sensors {
            newSensorDetails = append(newSensorDetails, SensorDetail {
                Sensor: fmt.Sprintf("%v (%v)", sensor.SensorName, sensor.SensorId),
            })
        }
        sensorDetails = append(sensorDetails, newSensorDetails...)
        return sensorDetails
    default:
        // sensorDetails := convertTreeToFlat(v)
        return nil
    }
    return nil
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
    sensorDetails := convertTreeToFlat(sensorTree)
    fmt.Printf("%v sensors found:\n",len(sensorDetails))
    for i, s := range(sensorDetails) {
        fmt.Printf("%v - %v - %v - %v\n", i, s.GroupLadder, s.Device, s.Sensor)
    }
}
