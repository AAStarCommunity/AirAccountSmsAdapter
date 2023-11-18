package infra

import (
	"AirAccountSmsAdapter/conf"
	"bytes"
	"fmt"
	"github.com/totoval/framework/helpers/log"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Qb struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Balance string `json:"balance"`
		Status  int    `json:"status"`
	} `json:"data"`
	Cost string `json:"cost"`
}

const (
	BindWallet   = "create"                            // 绑定钱包指令
	QueryBalance = "query"                             // 查询钱包余额
	TransferTo   = `transfer\s+([\d\.]+)\s+to\s+(\d+)` // 转账
)

// Retrieve 读取短信内容
func Retrieve(chip *Sim800c, smsIndex uint) error {
	if err := chip.Write([]byte(fmt.Sprintf("AT+CMGR=%d\r\n", smsIndex))); err != nil {
		return err
	}
	return nil
}

// InstructionOp 指令操作
func InstructionOp(chip *Sim800c, from string, rawMsg string) error {
	cfg := conf.GetAirCenterHost()
	// bind: C
	if strings.EqualFold(rawMsg, BindWallet) {
		from = strings.TrimPrefix(from, "+")
		if resp, err := http.Post(cfg+"/api/instructions/bind?id="+from, "application/json", bytes.NewBuffer([]byte("{}"))); err != nil {
			return log.Error(err)
		} else {
			log.Info("bind:" + resp.Status)
			go SendMessage(chip, from, "Congratulations! Your AirAccount Created!")
		}
	} else if strings.EqualFold(rawMsg, QueryBalance) {
		from = strings.TrimPrefix(from, "+")
		if resp, err := http.Get(cfg + "/api/instructions/balance?id=" + from); err != nil {
			return log.Error(err)
		} else {
			log.Info("query balance:" + resp.Status)
			if data, err := io.ReadAll(resp.Body); err == nil {
				b := Qb{}
				json.Unmarshal(data, &b)
				go SendMessage(chip, from, fmt.Sprintf("Your balance is %s eth", b.Data.Balance))
			}
		}
	} else {
		re := regexp.MustCompile(TransferTo)
		if matches := re.FindStringSubmatch(rawMsg); len(matches) == 3 {
			value := matches[1]
			receiver := matches[2]
			body, _ := json.Marshal(struct {
				Receiver string `json:"receiver"`
				Value    string `json:"value"`
			}{
				Receiver: receiver,
				Value:    value,
			})
			if resp, err := http.Post(cfg+"/api/instructions/transfer?id="+from,
				"application/json",
				bytes.NewBuffer(body)); err != nil {
				return log.Error(err)
			} else {
				log.Info("transfer:" + resp.Status)
				b := struct {
					Op string `json:"op"`
				}{}
				go SendMessage(chip, from, "transfer accepted")
				go CheckTransfer(chip, from, b.Op)
			}
		} else {
			log.Info("invalid indication: " + rawMsg)
		}
	}

	return nil
}

func CheckTransfer(chip *Sim800c, from string, op string) {
	cfg := conf.GetAirCenterHost()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 10)

		if resp, err := http.Get(cfg + "/api/instructions/transfer/check?id=" + from + "&op=" + op); err != nil {
			log.Error(err)
		} else {
			log.Info("transfer check result:" + resp.Status)

			if resp.StatusCode == 200 {
				go SendMessage(chip, from, "transfer successful")
				return
			}
		}
	}
}
