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
	outputEnv()

	port, baud := conf.GetSim800c()
	if c, err := infra.NewSim800c(port, baud, 5*zone.Second); err == nil {
		gw := infra.New(c)

		//gw.Listen()
		gw.PollUnreadMessages()
	}
}

func getSerialComlist() {
	ports, _ := serial.GetPortsList()

	fmt.Printf("%#v", ports)
	for _, port := range ports {
		fmt.Printf("Find Serial Com: %v\n", port)
	}
}

func outputEnv() {
	fmt.Printf("airaccount host: %s\r\n", conf.GetAirCenterHost())
	a, b := conf.GetSim800c()
	fmt.Printf("serial com: %s %d\n", a, b)
}
