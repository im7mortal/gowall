package main

import (
	"net/url"
	"strings"
)

func getEscapedString(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

