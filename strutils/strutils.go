package strutils

import "strings"

// StringTrim .
func StringTrim(s string) string {
	return strings.Trim(s, "\r\n\t ")
}

// StringTrimBlank .
func StringTrimBlank(s string) string {
	return strings.Trim(s, " ")
}
