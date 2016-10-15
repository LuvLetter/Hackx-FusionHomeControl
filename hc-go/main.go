package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"log"
)

func turnLightOn() {
	log.Println("Turn Light On")
}

func turnLightOff() {
	log.Println("Turn Light Off")
}
func Lightbulb_setup(){
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

	t, err := hc.NewIPTransport(hc.Config{Pin: "32191123"}, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
func GetTempture() float64{
		return 0
}

func TemptureSensor_setup() {
	info := accessory.Info{
		Name:         "Generic Tempture Sensor",
		Manufacturer: "HDU LUG",
	}
  TemptureSensor := accessory.NewTemperatureSensor(info, GetTempture(), -35, 100, 0.5)
  t, err := hc.NewIPTransport(hc.Config{Pin: "32191123"}, TemptureSensor.Accessory)
	log.Println("32191123")
  if err != nil {
  	log.Fatal(err)
  }
	hc.OnTermination(func() {
		t.Stop()
	})
	t.Start()

}

func main() {
	Lightbulb_setup()
	TemptureSensor_setup()
}
