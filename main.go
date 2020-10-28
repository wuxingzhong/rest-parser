package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/wuxingzhong/rest-parser/parser"
	"log"
	"os"
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

	restInfoList, err := parser.RestParser(restFile, nil)
	tmpMap := map[string][]parser.RestInfo{}
	if restIndex > 0 {
		if len(restInfoList) >= restIndex {
			tmpMap["restList"] = []parser.RestInfo{restInfoList[restIndex-1]}
		}else{
			return
		}
	}else {
		tmpMap["restList"] = restInfoList
	}
	data, _:= json.Marshal(tmpMap)
	fmt.Printf("%v", string(data))
	return
}
