package infra

import (
	"errors"
	"fmt"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"time"
)

type Gateway struct {
	chip *Sim800c
}

func New(chip *Sim800c) *Gateway {
	return &Gateway{
		chip: chip,
	}
}
func (gw *Gateway) Listen() {

	go func() {
		for {
			gw.chip.Read()
		}
	}()

	go gw.chip.Write([]byte("AT+CSQ\r\n"))

	for {
		select {
		case b := <-gw.chip.Bytes():
			log.Info("Incoming data", toto.V{"data": string(b[:])})
			err := log.Error(gw.parse(b))
			if err != nil {
				go log.Debug(fmt.Sprintf("error: %s | %s", err.Error(), string(b[:])))
			}
		case err := <-gw.chip.Error():
			log.Panic(err)
		default:
			time.Sleep(time.Second)
		}

	}
}

func (gw *Gateway) parse(msg []byte) error {

	if msg == nil {
		return nil
	}
	if ParseOk(msg) {
		return nil
	}

	// parse sms index +CMTI: "SM",2
	if matched, smsIndex, err := ParseSmsIndex(msg); matched {
		if err != nil {
			return err
		}
		// sms receive event
		if err := Retrieve(gw.chip, smsIndex); err != nil {
			return err
		}

		return nil
	}

	if matched, sender, content, err := ParseSmsContent(msg); matched && err == nil {
		log.Info(sender + ":" + content)
		if err := InstructionOp(gw.chip, sender, content); err != nil {
			return err
		}
		return nil
	}

	return errors.New(fmt.Sprintf("Not a normal message: %s", string(msg[:]))) // not a valid
}
