package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)


func main() {
	var (
		line int
		rest_file string
	)
	app := &cli.App{
		Name: "rest-parse",
		Usage: "解析rest文件",
		Flags: []cli.Flag {
			&cli.IntFlag{
				Name:        "line",
				Aliases: []string{
					"l",
				},
				Usage:       "",
				Destination: &line,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() >0 {
				rest_file = c.Args().First()
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
	restFileBuf, err:= ioutil.ReadFile(rest_file)
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
			log.Fatalf("read %s err, too short", rest_file)
		}
	}
	restList := make([]RestInfo,0)
	lastRest := RestInfo{Header:map[string]string{}}
	flag := 0
	for _, v := range lines {
		v = strings.Trim(v," ")
		if len(v) == 0 {
			continue
		}
		if strings.HasPrefix(v , "###") {
			if len(lastRest.Path) > 0 {
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
			tmpStr := strings.Split(v,":")
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
		restList = append(restList, lastRest)
		lastRest = RestInfo{}
	}

	tmpMap := map[string][]RestInfo{
		"restList": restList,
 	}
	data, _:= json.Marshal(tmpMap)
	fmt.Printf("%v", string(data))
	return

}
