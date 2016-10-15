package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"log"
	"encoding/json"
	"net/http"
)

type sensorResponse struct {
  Temperature float64
	hasPeopleMoved bool
	CO int
	heartRate int
}

var NodeIP string = "192.168.0.103/json"
var responseData = sensorResponse {26.312, true, 183, 0}

func getJson(url string, target interface{}) error {
  r, err := http.Get(url)
  if err != nil {
    return err
  }
  defer r.Body.Close()
	log.Println(r.Body)
  return json.NewDecoder(r.Body).Decode(target)
}

func GetHasPeple() float64{
	getJson(NodeIP, responseData)
  return responseData.hasPeopleMoved
}

func TemptureSensor_setup() {
	info := accessory.Info{
		Name:         "Generic Temperature Sensor",
		Manufacturer: "HDU LUG",
	}
	TemptureSensor := accessory.NewTemperatureSensor(info, GetTempture(), -35, 100, 0.5)
	t, err := hc.NewIPTransport(hc.Config{Pin: "00000001"}, TemptureSensor.Accessory)
	if err != nil {
		log.Fatal(err)
	}
	t.Start()

}

func main() {
  log.Println("Start:")
	TemptureSensor_setup()
}
