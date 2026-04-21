package device

import (
	"net/url"
	"strconv"
	"strings"
)

func buildAttachmentDisposition(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		fileName = "firmware.bin"
	}
	return "attachment; filename*=UTF-8''" + url.QueryEscape(fileName)
}

func int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}
