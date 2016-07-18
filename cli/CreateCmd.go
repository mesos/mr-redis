package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/codegangsta/cli"
)

//CreateCmd sub-command implementation, for now it takes name, memory and number of slaves
//If wait is bool flag is enabled then we will keep polling scheduler/framework until we get a "Created OK" from it.
func CreateCmd(c *cli.Context) {

	name := c.String("name")
	mem := c.Int("memory")
	slaves := c.Int("slaves")
	isWait := c.Bool("wait")

	if name == "" {
		fmt.Printf("Error: Should have a valid name")
	}

	if mem < 100 {
		mem = 100
	}

	if slaves < 0 || slaves > 100 {
		slaves = 0
	}
	fmt.Printf("Attempting to Creating a Redis Instance (%s) with mem=%d slaves=%d\n", name, mem, slaves)

	url := fmt.Sprintf("%s/v1/CREATE/%s/%d/1/%d", MrRedisFW, name, mem, slaves)
	res, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Printf("Error: Creating the Instance error=%v\n", err)
		return
	}

	if res.StatusCode == http.StatusCreated {

		fmt.Printf("Instance Creation accepted..")

		if isWait {
			cnt := 0
			for !IsRunning(name) {
				time.Sleep(100 * time.Millisecond)
				cnt++
				p := cnt % 10
				if p == 0 {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\nInstance Created.\n")
		} else {
			fmt.Printf("\nCheck $mrr status -n %s for status\n", name)
		}
	} else {
		fmt.Printf("Error Creating the instance response = %v\n", res)
	}

}
