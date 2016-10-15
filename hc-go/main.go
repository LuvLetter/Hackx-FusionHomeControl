package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"log"
	"sync"
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

func GetSensorJson(){
	getJson(NodeIP, responseData)
}

func turnLightOn() {
	log.Println("Turn Light On")
}

func turnLightOff() {
	log.Println("Turn Light Off")
}
func Lightbulb_setup() {
	info := accessory.Info{
		Name:         "Personal Light Bulb",
		Manufacturer: "Matthias",
	}

	acc := accessory.NewLightbulb(info)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			turnLightOn()
		} else {
			turnLightOff()
		}
	})

	t, err := hc.NewIPTransport(hc.Config{Pin: "00000002"}, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
func GetTempture() float64 {
	return 0
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
	hc.OnTermination(func() {
		t.Stop()
	})
	t.Start()

}

func main() {
	GetSensorJson()
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		GetSensorJson()
	}()

	go func() {
		defer wg.Done()
		TemptureSensor_setup()

	}()
	go func() {
		defer wg.Done()
		Lightbulb_setup()
	}()

	wg.Wait()


	// signals := make(chan uint8, 2)
	// go func() {
	// 	Lightbulb_setup()
	// 	log.Println("light")
	// 	signals <- 1
	// }()
	// go func() {
	// 	TemptureSensor_setup()
	// 	log.Println("temperature")
	// 	signals <- 1
	//
	// }()
	// <-signals
	// <-signals
}
