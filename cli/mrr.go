package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

var MrRedisFW string //Frameworks IP and Port number

func Init() {
	//Check if we have a ~/.MrRedis config file in the system already,
	f, err := os.Open("/tmp/.MrRedis")
	if err != nil {
		fmt.Printf("Cli is not initalized err=%v\n", err)
		fmt.Printf("$mrr init <http://MrRedisEndPoint>\n")
		return
	}
	//If yes then open it and read the content (first line)
	d := make([]byte, 512)
	if c, err := f.Read(d); err != nil && c != 0 {
		fmt.Printf("Unable to read the config file err=%v\n", err)
		return
	}
	//It should have IP:Port format
	//MrRedisFW = string(d)
	MrRedisFW = strings.Trim(string(d), "\x00")
	f.Close()
	return
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
			Usage:   "$mrr init <http://MrRedisEndPoint>",
			Action:  InitCmd,
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create a Redis Instance",
			Action:  CreateCmd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of the Redis Instance",
				},
				cli.IntFlag{
					Name:  "memory, m",
					Usage: "Memory in MB",
				},
				cli.IntFlag{
					Name:  "slaves, s",
					Usage: "Number of Slaves",
				},
				cli.BoolFlag{
					Name:  "wait, w",
					Usage: "Wait for the Instnace to be create (by default the command is async)",
				},
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Status of a Redis Instance",
			Action:  StatusCmd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of the Redis Instance",
				},
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
