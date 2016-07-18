package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/cli"
)

func httpDelete(url string) (*http.Response, error) {
	var client http.Client
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		fmt.Printf("Error: Unable to create NewRequsterr=%v", err)
		return nil, err
	}
	// ...
	//req.Header.Add("If-None-Match", `W/"wyzzy"`)
	resp, err := client.Do(req)

	return resp, err
}

//DeleteCmd sub command implementation for DELETE instance
//Simply fires the DELTE REST api to the scheduler
func DeleteCmd(c *cli.Context) {

	name := c.String("name")

	if name == "" {
		fmt.Printf("Error: Should have a valid name")
	}

	url := fmt.Sprintf("%s/v1/DELETE/%s", MrRedisFW, name)
	res, err := httpDelete(url)
	if err != nil {
		fmt.Printf("Error: Deleting the Instance error=%v\n", err)
		return
	}

	if res.StatusCode == http.StatusOK {

		fmt.Printf("Instance Deletion initated..")

	} else {
		fmt.Printf("Error Creating the instance response = %v\n", res)
	}

}
