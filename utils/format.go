package utils

import (
	"fmt"
	"strconv"
	"strings"
)

//Used to avoid precision issues with float arithmetic
func FormatAmtStrToInt64(str string, precison int) (int64, error) {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0, nil
	}

	parts := strings.Split(str, ".") //split on decimal point

	if len(parts) == 1 { //no decimal point
		parts = append(parts, strings.Repeat("0", precison)) //

	} else if diff := precison - len(parts[1]); diff > 0 { //less numbers after decimal point than precision specified
		parts[1] += strings.Repeat("0", diff)
	}

	return strconv.ParseInt(parts[0]+parts[1], 10, 64)
}

func FormatInt64AmtToString(num int64, precison int) string {
	return fmt.Sprintf("%d.%0"+fmt.Sprint(precison)+"d", num/100, num%100)
}
