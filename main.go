package main

import (
	"AirAccountSmsAdapter/conf"
	"AirAccountSmsAdapter/infra"
	"fmt"
	"github.com/totoval/framework/helpers/zone"
	"go.bug.st/serial"
)

func main() {
	getSerialComlist()
	port, baud := conf.GetSim800c()
	if c, err := infra.NewSim800c(port, baud, 5*zone.Second); err == nil {
		gw := infra.New(c)

		gw.Listen()
	}
}

func getSerialComlist() {
	ports, _ := serial.GetPortsList()

	fmt.Printf("%#v", ports)
	for _, port := range ports {
		fmt.Printf("Find Serial Com: %v\n", port)
	}
}
