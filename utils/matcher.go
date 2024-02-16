package utils

import "regexp"

func MatchString(regex, str string) bool {
	expr, err := regexp.Compile(regex)
	if err != nil {
		return false
	}

	return expr.MatchString(str)
}
