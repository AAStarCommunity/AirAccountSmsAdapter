package infra

import (
	"bytes"
	"fmt"
	"github.com/totoval/framework/helpers/log"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/tarm/serial"
	"github.com/totoval/framework/helpers/zone"
)

const REGEXP_CMGL = `\+CMGL: ([0-9]+),".*?","(\+[0-9]+)",".*?",".*?"`

type Sim800c struct {
	writer       *serial.Port
	conf         *serial.Config
	chB          chan []byte
	chErr        chan error
	isReading    bool
	SmsThreshold int
}

func NewSim800c(comPort string, baudRate int, readTimeout zone.Duration, smsThreshold int) (*Sim800c, error) {
	s := &Sim800c{
		chB:          make(chan []byte, 50*3), // "+CPMS: \"SM_P\",50,50,\"SM_P\",50,50,\"SM_P\",50,50"
		chErr:        make(chan error, 50*3),
		SmsThreshold: smsThreshold,
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
	re := regexp.MustCompile(REGEXP_CMGL)
	for {
		_b := make([]byte, 512)

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

		if bytes.HasSuffix([]byte(strings.TrimRight(string(b), "\r\n")), []byte("OK")) {
			//str := string(bytes.Trim(b, "\x00"))
			//log.Info("recv raw: " + str)
			__b := bytes.Trim(b, "\r\n")
			msgArr := bytes.Split(__b, []byte("\r\n")) // [][data]
			line := 1
			for _, msg := range msgArr {
				str := string(msg)
				matches := re.FindStringSubmatch(str)
				if matches != nil && len(matches) > 2 {
					from := matches[2]
					text := msgArr[line]
					msg = []byte(fmt.Sprintf("%s<br />%s", from, text))
					s.chB <- msg
				}
				line++
			}
			b = make([]byte, 512)
		}

		time.Sleep(2 * time.Second)
	}
}

func (s *Sim800c) Write(b []byte) error {
	if err := s.flush(); err != nil {
		return err
	}
	_, err := s.write(b)
	if err != nil {
		return err
	}
	log.Info("send msg: " + string(b[:]))
	time.Sleep(time.Second)
	return nil
}

func (s *Sim800c) Error() <-chan error {
	return s.chErr
}
func (s *Sim800c) Bytes() <-chan []byte {
	return s.chB
}
