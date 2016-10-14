package main

import (
	"github.com/LuvLetter/hc"
	"github.com/LuvLetter/hc/accessory"
	"log"
)

func GetTempture(){

}

func main() {
	info := accessory.Info{
		Name:         "Generic Tempture Sensor",
		Manufacturer: "HDU LUG",
	}
  for i := 0; true; i++ {
    TemptureSensor := accessory.NewTemperatureSensor(info, GetTempture, -35, 100, 0.5)
    t, err := hc.NewIPTransport(hc.Config{Pin: "32191123"}, TemptureSensor.Accessory)
  	if err != nil {
  		log.Fatal(err)
  	}
  }
}
