package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var MrRedisFW string //Frameworks IP and Port number

func Init() {
	//Check if we have a ~/.MrRedis config file in the system already,
	//If yes then open it and read the content (first line)
	//It should have IP:Port format

}

func main() {
	Init()
	app := cli.NewApp()
	app.Name = "mrr"
	app.Usage = "MrRedis Command Line Interface"
	app.HideVersion = true
	app.Action = func(c *cli.Context) {
		println("MrRedis Command Line")
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initalize the cli",
			Action: func(c *cli.Context) {
				println("Initalized : ", c.Args().First())
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create a Redis Instance",
			Action: func(c *cli.Context) {
				println("Create Redis instance: ", c.Args().First())
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Status of a Redis Instance",
			Action: func(c *cli.Context) {
				println("Status of redis instance is..: ", c.Args().First())
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete a Redis Instance",
			Action: func(c *cli.Context) {
				println("Redis Instnace is deleted..: ", c.Args().First())
			},
		},
	}

	app.Run(os.Args)
}
