package main

import (
	"cdncheck/config"
	check "cdncheck/modules"
	"encoding/json"
	"fmt"
	"strings"
)

type content struct {
	Message string       `json:"message"`
	Data    check.Result `json:"data"`
	Status  int          `json:"status"`
}

func printJson(ct content) {
	b, err := json.Marshal(ct)
	if err != nil {
		return
	}
	fmt.Printf("%s", b)
}

func main() {
	var ct content
	var Args config.CommandLineArgs
	config.Flag(&Args)

	if Args.Targets == "" {
		ct.Message = "Targets not provided"
		ct.Status = -1
		printJson(ct)
		return
	}

	if Args.Filepath == "" {
		ct.Message = "CDN data not provided"
		ct.Status = -1
		printJson(ct)
		return
	}

	cli, err := check.InitCdnClient(Args.Filepath)
	if err != nil {
		ct.Message = err.Error()
		ct.Status = -1
		printJson(ct)
		return
	}

	r := cli.RangeCheck(strings.Split(Args.Targets, ","))
	ct.Data = r
	ct.Status = 1
	printJson(ct)
}
