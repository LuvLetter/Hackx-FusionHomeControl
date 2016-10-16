package main

import (
        "github.com/brutella/hc"
        "github.com/brutella/hc/accessory"
        "log"
        "net/http"
)

func turnLight(stat bool) {
  var url := "10.221.64.122/gpio"
  if(stat) r, err := http.Get(url + "/1")
  else r, err := http.Get(url)
  if err != nil {
    return err
}


func main() {
        info := accessory.Info{
                Name:         "Personal Light Bulb",
                Manufacturer: "HDU LUG",
        }

        acc := accessory.NewLightbulb(info)

        acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
                if on == true {
                        turnLight(true)
                } else {
                        turnLight(false)
                }
        })

        t, err := hc.NewIPTransport(hc.Config{Pin: "00000076"}, acc.Accessory)
        if err != nil {
                log.Fatal(err)
        }

        hc.OnTermination(func() {
                t.Stop()
        })

        t.Start()
}
