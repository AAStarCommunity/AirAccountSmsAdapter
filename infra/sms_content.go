package infra

import (
	"bytes"
	"errors"
	"regexp"
)

const REGEXP_CMGR = `^\+CMGR: "[^"]+","(\+?\d+)","[^"]*","[^"]+"$`

func ParseSmsContent(msg []byte) (matched bool, sender string, content string, err error) {
	if msg == nil {
		return false, "", "", nil
	}
	matched = false
	msgStrArr := bytes.Split(msg, []byte("<br />"))
	re := regexp.MustCompile(REGEXP_CMGR)

	if m := re.FindStringSubmatch(string(msgStrArr[0][:])); len(m) > 1 {
		sender = m[1]
		content = string(msgStrArr[1][:])
		matched = true
		return
	}

	if !matched {
		return matched, "", "", errors.New("not a message content")
	}

	return matched, "", "", err
}
