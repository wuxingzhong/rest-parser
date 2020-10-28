package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)


func main() {
	var (
		restIndex int
		restFile  string
	)
	app := &cli.App{
		Name: "rest-parse",
		Usage: "解析rest文件",
		ArgsUsage: "rest文件",
		Flags: []cli.Flag {
			&cli.IntFlag{
				Name:        "index",
				Aliases: []string{
					"i",
				},
				Usage:       "",
				Destination: &restIndex,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() >0 {
				restFile = c.Args().First()
			}
			return nil
		},

	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 处理文件
	restFileBuf, err:= ioutil.ReadFile(restFile)
	if err !=nil {
		return
	}
	m := map[string]string{
		"host": "",
		"Authorization": "",
	}
	reg := regexp.MustCompile("{{(.+)}}")
	str := reg.ReplaceAllStringFunc(string(restFileBuf), func(old string ) string {
		old = strings.ReplaceAll(old,"{{","")
		old = strings.ReplaceAll(old,"}}","")
		v, b := m[old]
		if !b {
			return ""
		}
		return v
	})

	// 解析
	lines := strings.Split(str,"\n")
	if len(lines) < 2 {
		lines = strings.Split(str,"\r\n")
		if len(lines) < 2 {
			log.Fatalf("read %s err, too short", restFile)
		}
	}
	restList := make([]RestInfo,0)
	lastRest := RestInfo{Header:map[string]string{}}
	flag := 0
	index := 1
	for _, v := range lines {
		v = strings.Trim(v," ")
		if len(v) == 0 {
			if flag == 3 {
				flag = 4
			}
			continue
		}
		if strings.HasPrefix(v , "###") {
			if len(lastRest.Path) > 0 {
				lastRest.Index = index
				index = index + 1
				restList = append(restList, lastRest)
				lastRest = RestInfo{Header:map[string]string{}}
				flag = 0
			}
			flag = 1
		}
		switch flag {
		case 1:
			// comment
			lastRest.Comment = strings.Trim(v, "#")
			lastRest.Comment = strings.Trim(lastRest.Comment, " ")
			flag = 2
		case 2:
			 tmpStr := strings.Split(v," ")
			 lastRest.Method = tmpStr[0]
			 lastRest.Path = tmpStr[1]
			 flag = 3
		case 3:
			tmpStr := strings.Split(v,": ")
			if len(tmpStr) != 2 {
				lastRest.Body =  lastRest.Body + v
				flag = 4
			}else{
				if lastRest.Header == nil {
					lastRest.Header = map[string]string{}
				}
				lastRest.Header[tmpStr[0]] = tmpStr[1]
			}
		case 4:
			lastRest.Body =  lastRest.Body + v
		}

	}
	if len(lastRest.Path) > 0 {
		lastRest.Index = index
		restList = append(restList, lastRest)
		lastRest = RestInfo{}
	}

	tmpMap := map[string][]RestInfo{}
	if restIndex > 0 {
		if len(restList) >= restIndex {
			tmpMap["restList"] = []RestInfo{restList[restIndex-1]}
		}else{
			return
		}
	}else {
		tmpMap["restList"] = restList
	}
	data, _:= json.Marshal(tmpMap)
	fmt.Printf("%v", string(data))
	return
}
