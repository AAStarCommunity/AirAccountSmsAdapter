package infra

import (
	"github.com/totoval/framework/helpers/log"
)

func SendMessage(chip *Sim800c, sender string, msg string) {
	if err := chip.Write([]byte("AT+CMGF=1\r\n")); err != nil {
		log.Error(err)
		return
	}

	if err := chip.Write([]byte("AT+CMGS=\"" + sender + "\"\r\n")); err != nil {
		log.Error(err)
		return
	}

	if err := chip.Write([]byte(msg + string(rune(26)))); err != nil {
		log.Error(err)
		return
	}
}
