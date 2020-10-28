package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// restInfo
type RestInfo struct {
	// 序号
	Index int
	// 注释
	Comment string
	// 请求头
	Header map[string]string
	// 请求方法
	Method string
	// 路径
	Path string
	// body
	Body string
}

const (
	_restNull = iota
	_restComment
	_restPath
	_restHeader
	_restBody
)

// RestParser 解析rest文件  filename: 文件名 , varMap: 替换变量列表
func RestParser(filename string, varMap map[string]string) (restInfoList []RestInfo, err error) {
	// 处理文件
	restFileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	// 替换变量列表
	reg := regexp.MustCompile("{{(.+)}}")
	str := reg.ReplaceAllStringFunc(string(restFileBuf), func(old string) string {
		old = strings.ReplaceAll(old, "{{", "")
		old = strings.ReplaceAll(old, "}}", "")
		if varMap == nil {
			return ""
		}
		v, b := varMap[old]
		if !b {
			return ""
		}
		return v
	})

	// 解析
	lines := strings.Split(str, "\n")
	if len(lines) < 2 {
		lines = strings.Split(str, "\r\n")
		if len(lines) < 2 {
			err = fmt.Errorf("read %s err, too short", filename)
			return
		}
	}
	restInfoList = make([]RestInfo, 0, 10)
	lastRest := RestInfo{Header: map[string]string{}}
	flag := _restNull
	index := 1
	for _, v := range lines {
		v = strings.Trim(v, " ")
		if len(v) == 0 {
			if flag == _restHeader {
				flag = _restBody
			}
			continue
		}
		if strings.HasPrefix(v, "###") {
			if len(lastRest.Path) > 0 {
				lastRest.Index = index
				index = index + 1
				restInfoList = append(restInfoList, lastRest)
				lastRest = RestInfo{Header: map[string]string{}}
				flag = _restNull
			}
			flag = _restComment
		}
		switch flag {
		case _restComment:
			// comment
			lastRest.Comment = strings.Trim(v, "#")
			lastRest.Comment = strings.Trim(lastRest.Comment, " ")
			flag = _restPath
		case _restPath:
			tmpStr := strings.Split(v, " ")
			lastRest.Method = tmpStr[0]
			lastRest.Path = tmpStr[1]
			flag = _restHeader
		case _restHeader:
			tmpStr := strings.Split(v, ": ")
			if len(tmpStr) != 2 {
				lastRest.Body = lastRest.Body + v
				flag = _restBody
			} else {
				if lastRest.Header == nil {
					lastRest.Header = map[string]string{}
				}
				lastRest.Header[tmpStr[0]] = tmpStr[1]
			}
		case _restBody:
			lastRest.Body = lastRest.Body + v
		}
	}
	if len(lastRest.Path) > 0 {
		lastRest.Index = index
		restInfoList = append(restInfoList, lastRest)
		lastRest = RestInfo{}
	}
	return
}
