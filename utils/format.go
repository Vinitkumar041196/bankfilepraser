package utils

import (
	"strconv"
	"strings"
)

//assumption decimal precison is 2
func FormatAmtToInt64(numStr string) (int64, error) {
	numStr = strings.TrimSpace(numStr)
	if numStr == "" {
		return 0, nil
	}
	if strings.Contains(numStr, ".") {
		return strconv.ParseInt(strings.ReplaceAll(numStr, ".", ""), 10, 64)
	}
	return strconv.ParseInt(numStr+"00", 10, 64)
}
