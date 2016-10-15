package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"log"
	"sync"
 	"encoding/json"
)

func getJson(url string, target interface{}) error {
    r, err := http.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
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

type sensorResponse struct {
    Temperature float64
		hasPeopleMoved bool
		CO int
		heartRate int
}

func GetTempture() float64 {
	responseData := new(sensorResponse)
	getJson(devIP+"json", responseData)
	log.Println(responseData.Temperature)
	return sensorResponse.Temperature
}
func GetIRsensor() bool {
	responseData := new(sensorResponse)
	getJson(devIP+"json", responseData)
	log.Println(responseData.hasPeopleMoved)
	return responseData.hasPeopleMoved
}
func GetCO() int {
	responseData := new(sensorResponse)
	getJson(devIP+"json", responseData)
	log.Println(responseData.CO)
	return responseData.CO
}
func GetHeartRate() int {
	responseData := new(sensorResponse)
	getJson(devIP+"json", responseData)
	log.Println(responseData.heartRate)
	return responseData.heartRate
}

func TemptureSensor_setup() {
	info := accessory.Info{
		Name:         "Generic Tempture Sensor",
		Manufacturer: "HDU LUG",
	}
	TemptureSensor := accessory.NewTemperatureSensor(info, GetTempture(), -35, 100, 0.5)
	t, err := hc.NewIPTransport(hc.Config{Pin: "00000001"}, TemptureSensor.Accessory)
	log.Println("0000001")
	if err != nil {
		log.Fatal(err)
	}
	hc.OnTermination(func() {
		t.Stop()
	})
	t.Start()

}

func main() {
	devIP := "10.221.65.124/"
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		Lightbulb_setup()
	}()
	go func() {
		defer wg.Done()
		TemptureSensor_setup()
	}()

	wg.Wait()

	// or
	// signals := make(chan uint8, 2)
	// go func() {
	// 	Lightbulb_setup()
	// 	signals <- 1
	// }()
	// go func() {
	// 	TemptureSensor_setup()
	// 	signals <- 1
	// }()
	// <-signals
	// <-signals
}
