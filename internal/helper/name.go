package helper

import "strings"

func NormaliseName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", ""))
}
