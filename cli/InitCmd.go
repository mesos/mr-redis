package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
)

func InitCmd(c *cli.Context) {
	//
	EP := c.Args().First()
	if !(strings.Contains(EP, "http:")) {
		fmt.Printf("Error: The end point should contain http or https")
		return
	}

	var conf_file_path string
	if runtime.GOOS == "windows" {
		conf_file_path = ".MrRedis"
	} else {
		conf_file_path = "/tmp/.MrRedis"
	}

	f, err := os.Create(conf_file_path)

	if err != nil {
		fmt.Printf("Error: Unable to create config file err=%v\n", err)
		return
	}

	if _, err := http.Get(fmt.Sprintf("%s/v1/STATUS", EP)); err != nil {
		fmt.Printf("Error: Testing the End Point err=%v\n", err)
		return
	}

	if _, err := f.WriteString(EP); err != nil {
		fmt.Printf("Error: Unable to write to config file err=%v\n", err)
	}
	f.Close()
	return
}
