package parser

import (
	"regexp"
	"strings"
)

type VarMap map[string]string

func (v VarMap) Replace(old string) string {
	reg := regexp.MustCompile("{{(.+?)}}")
	str := reg.ReplaceAllStringFunc(old, v.replaceFunc)
	return str
}

func (v VarMap) replaceFunc(old string) string {
	old = strings.ReplaceAll(old, "{{", "")
	old = strings.ReplaceAll(old, "}}", "")
	if v == nil {
		return ""
	}
	value, b := v[old]
	if !b {
		return ""
	}
	return value
}
