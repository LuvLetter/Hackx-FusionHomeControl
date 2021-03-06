package main

import (
        "github.com/brutella/hc"
        "github.com/brutella/hc/accessory"
        "log"
        "encoding/json"
        "net/http"
)
var responseData = lightReturn {"high"}

func getJson(url string, target interface{}) error {
  r, err := http.Get(url)
  if err != nil {
    return err
  }
  defer r.Body.Close()
	log.Println(r.Body)
  return json.NewDecoder(r.Body).Decode(target)
}
type lightReturn struct {
  successful string
}
// func turnLightOn() {
//
// }
// func turnLightOff() {
// 	getJson("10.221.64.122/gpio/0", responseData)
//   log.Println(responseData.successful)
// }



func main() {
        info := accessory.Info{
                Name:         "Personal Light Bulb",
                Manufacturer: "HDU LUG",
        }

        acc := accessory.NewLightbulb(info)

        acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
                if on == true {
                  //10.221.64.122
                  getJson("10.221.64.122/gpio/1", responseData)
                  log.Println(responseData.successful)
                } else {
                  getJson("10.221.64.122/gpio/0", responseData)
                  log.Println(responseData.successful)
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
