package main

//initial helper for run the executor functionality from the command line without mesos
//the same (or similar :)) code will go into mesos-go executor launchtask function

import (
	"./serviceproc"
	"fmt"
)

func main() {
	fmt.Println("Starting Redis server on given port\n")

	//tbd: only the service instance id needs to be passed here; how to get it?
	//tbd: who gets you the port value?
	redisproc := serviceproc.NewRedisProc("ServiceInstnsID", 6379)
	monitor := serviceproc.NewProcMonitor(redisproc)

	monitor.SpawnandMonitor()

}
