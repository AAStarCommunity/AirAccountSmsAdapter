package main

import (
	"AirAccountSmsAdapter/conf"
	"AirAccountSmsAdapter/infra"
	"github.com/totoval/framework/helpers/zone"
)

func main() {
	port, baud := conf.GetSim800c()
	if c, err := infra.NewSim800c(port, baud, 5*zone.Second); err == nil {
		gw := infra.New(c)

		gw.Listen()
	}
}
