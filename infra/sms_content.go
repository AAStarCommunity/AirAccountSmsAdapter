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
	re, err := regexp.Compile(REGEXP_CMGR)
	if err != nil {
		return false, "", "", err
	}

	if m := re.FindStringSubmatch(string(msgStrArr[0][:])); m != nil && len(m) > 1 {
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
