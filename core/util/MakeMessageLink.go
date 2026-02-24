package util

import (
	"fmt"
	"strconv"
	"strings"
)

func MakeMessageLink(groupId int64, messageId int) string {
	// Private supergroup: strip leading -100
	raw := strconv.FormatInt(groupId, 10)
	if strings.HasPrefix(raw, "-100") {
		raw = raw[4:]
	} else {
		raw = strings.TrimPrefix(raw, "-")
	}
	return fmt.Sprintf("https://t.me/c/%s/%d", raw, messageId)
}
