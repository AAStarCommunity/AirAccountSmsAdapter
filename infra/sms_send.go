package infra

import (
	"github.com/totoval/framework/helpers/log"
)

func SendMessage(chip *Sim800c, sender string, msg string) error {
	if err := chip.Write([]byte("AT+CMGF=1\r\n")); err != nil {
		return log.Error(err)
	}

	if err := chip.Write([]byte("AT+CMGS=\"" + sender + "\"\r\n")); err != nil {
		return log.Error(err)
	}

	if err := chip.Write([]byte(msg + string(rune(26)))); err != nil {
		return log.Error(err)
	}

	return nil
}
