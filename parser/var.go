package parser

import "strings"

type VarMap map[string]string

func (v VarMap) ReplaceFunc(old string) string {
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
