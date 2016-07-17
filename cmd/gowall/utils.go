package main

import (
	"net/url"
	"strings"
	"regexp"
)

func getEscapedString(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

var r1, _ = regexp.Compile(`[^\w ]+`)
var r2, _ = regexp.Compile(` +`)

/**
	preparing id
 */

func slugify (str string) string {
	str = strings.ToLower(str)
	str = r1.ReplaceAllString(str, "")
	str = r2.ReplaceAllString(str, "-")
	return str
}

func slugifyName(str string) string {
	str = strings.TrimSpace(str)
	return r2.ReplaceAllString(str, " ")
}
