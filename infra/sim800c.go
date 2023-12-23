package infra

import (
	"bytes"
	"github.com/tarm/serial"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
	"io"
	"strings"
	"time"
)

type Sim800c struct {
	writer    *serial.Port
	conf      *serial.Config
	chB       chan []byte
	chErr     chan error
	isReading bool
}

func NewSim800c(comPort string, baudRate int, readTimeout zone.Duration) (*Sim800c, error) {
	s := &Sim800c{
		chB:   make(chan []byte, 50*3), // "+CPMS: \"SM_P\",50,50,\"SM_P\",50,50,\"SM_P\",50,50"
		chErr: make(chan error, 50*3),
	}
	s.conf = &serial.Config{Name: comPort, Baud: baudRate, ReadTimeout: readTimeout}

	var err error
	if s.writer, err = serial.OpenPort(s.conf); err != nil {
		return nil, err
	}

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}
func (s *Sim800c) init() error {
	if err := s.Write([]byte("AT+CMGF=1\r\n")); err != nil {
		return err
	}
	if err := s.Write([]byte("AT+CSCS=\"GSM\"\r\n")); err != nil {
		return err
	}
	if err := s.Write([]byte("AT+CNMI=2,1\r\n")); err != nil {
		return err
	}
	return nil
}
func (s *Sim800c) Close() error {
	defer close(s.chB)
	defer close(s.chErr)
	return s.writer.Close()
}
func (s *Sim800c) flush() error {
	return s.writer.Flush()
}
func (s *Sim800c) read(b []byte) (int, error) {
	return s.writer.Read(b)
}
func (s *Sim800c) write(b []byte) (int, error) {
	return s.writer.Write(b)
}
func (s *Sim800c) Read() {
	if s.isReading {
		return
	}

	s.isReading = true // only read once
	defer func() {
		s.isReading = false
	}()

	var b []byte
	for {
		_b := make([]byte, 128)

		_n, err := s.read(_b)
		if err != nil {
			if err == io.EOF {
				if len(b) <= 0 {
					// no message, continue receiving
					continue
				}

				if _n <= 0 {
					// received finished
					return
				}

				// len(b) > 0 && _n > 0 received aborted
				s.chErr <- io.EOF
				s.chB <- b
				return
			}
			// received error
			s.chErr <- err
			return
		}

		if _n > 0 {
			b = append(b, _b[:_n]...)
		}

		if bytes.HasSuffix(b, []byte("\r\n")) {
			log.Info("\r\n{[recv]" + strings.Trim(string(b), "\r\n}\r\n"))
		}

		time.Sleep(2 * time.Second)
		//if bytes.HasPrefix(b, []byte("\r\n")) && bytes.HasSuffix(b, []byte("\r\n")) {
		//	if bytes.Contains(b, []byte("\r\n\r\n")) {
		//		__b := bytes.Trim(b, "\r\n")
		//		msgArr := bytes.Split(__b, []byte("\r\n\r\n")) // [][data]
		//		log.Warn(len(msgArr))
		//		for _, msg := range msgArr {
		//			log.Warn(msg)
		//			s.chB <- msg
		//		}
		//	} else {
		//		// single msg bytes
		//		__b := bytes.Trim(b, "\r\n") // data
		//		s.chB <- __b
		//	}
		//
		//	break // msg end
		//} else if bytes.HasPrefix(b, []byte("AT+CMGR=")) && bytes.Contains(b, []byte("\r\n")) {
		//	bytes.Split(b, []byte("\r\n"))
		//	msgArr := bytes.Split(b, []byte("\r\n")) // [][data]
		//	for idx, msg := range msgArr {
		//		if bytes.HasPrefix(msg, []byte("+CMGR:")) {
		//			msg = []byte(string(msg) + "<br />" + string(msgArr[idx+1]))
		//			s.chB <- msg
		//		}
		//	}
		//	break
		//}
	}
	return
}

func (s *Sim800c) Write(b []byte) error {
	if err := s.flush(); err != nil {
		return err
	}
	n, err := s.write(b)
	if err != nil {
		return err
	}
	log.Info("Send Bytes", toto.V{"bytes": string(b[:]), "length": n})
	time.Sleep(time.Second)
	return nil
}

func (s *Sim800c) Error() <-chan error {
	return s.chErr
}
func (s *Sim800c) Bytes() <-chan []byte {
	return s.chB
}
