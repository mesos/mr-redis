package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	typ "../common/types"
	"github.com/codegangsta/cli"
)

func IsRunning(name string) bool {

	ret := StatusOf(name)

	if ret == nil {
		return false
	} else {
		if ret.Status == typ.INST_STATUS_RUNNING {
			return true
		}
	}

	return false
}

func StatusOf(name string) *typ.Instance_Json {
	var ret typ.Instance_Json

	url := fmt.Sprintf("%s/v1/STATUS/%s", MrRedisFW, name)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: Gettin status of the instance=%v\n", err)
		return nil
	}

	if res.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: Unable to Read response body err=%v\n", err)
			return nil
		}

		err = json.Unmarshal(b, &ret)
		if err != nil {
			fmt.Printf("Error: Unable to unmarshal the json err=%v\n", err)
			return nil
		}

		return &ret

	}

	return nil

}

func StatusCmd(c *cli.Context) {

	name := c.String("name")

	if name == "" {
		fmt.Printf("Getting Status for all the running Instances\n")
		fmt.Printf("Status all is not implemented yet\n")
		return
	}

	inst := StatusOf(name)

	if inst == nil {
		fmt.Printf("Status not available\n")
		return
	}

	fmt.Printf("Status = %s\nType = %s\nCapacity = %d\nMaster = %s:%s\n", inst.Status, inst.Type, inst.Capacity, inst.Master.IP, inst.Master.Port)

	for i, s := range inst.Slaves {
		fmt.Printf("\tSlave%d = %s:%s\n", i, s.IP, s.Port)
	}
}
