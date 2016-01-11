package main

import (
	"log"

	"../common/types"
	"./cmd"
	"./httplib"
)

//Declare all the Constants to be used in this file

const (
	HTTP_SERVER_PORT = "8080"
)

//Main execution point

func main() {

	log.Printf("*****************************************************************")
	log.Printf("*********************Starting MrRedis-Scheduler******************")
	log.Printf("*****************************************************************")

	//Command line argument parsing

	//Initalize the common entities like store, store configuration etc.
	isInit, err := types.Initialize("etcd", "http://10.11.12.24:2379")
	if err != nil || isInit != true {
		log.Fatalf("Failed to intialize Error:%v return %v", err, isInit)
	}

	//Start the creator

	go cmd.Creator()

	//Start HTTP server and related things to handle restfull calls to the scheduler
	httplib.Run(HTTP_SERVER_PORT)

	//Start mesos scheduler and related things

	//Start Creater

	//Start the Maintainer

	//Start the destroyer

	//Wait for termination signal

	log.Printf("*****************************************************************")
	log.Printf("*********************Finished MrRedis-Scheduler******************")
	log.Printf("*****************************************************************")

}
