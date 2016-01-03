package main

//initial helper for run the executor functionality from the command line without mesos
//the same (or similar :)) code will go into mesos-go executor launchtask function

import (
	"./serviceproc"
	"log"
	"time"
)

var procMap map[string](*serviceproc.RedisProc)

func main() {
	log.Printf("Starting Redis server on given port\n")

	procMap = make(map[string](*serviceproc.RedisProc))

	//tbd: only the service instance id needs to be passed here; how to get it?
	//tbd: who gets you the port value?
	for i := 0; i < 2; i++ {
		redisproc, uidStr := serviceproc.NewRedisProc("ServiceInstnsID", (6379 + i))
		procMap[uidStr] = redisproc
		log.Printf("spawning a new server with id:%s and proc:%v", uidStr, redisproc)
		monitor := serviceproc.NewProcMonitor(redisproc)

		//the launch task returns after spawning this
		go func(monitor *serviceproc.ProcMonitor) {

			monitor.SpawnandMonitor()

		}(monitor)
	}

	//put a logic to see what got filled in the map

	for i := 0; i < 10; i++ {

		log.Printf("sleeping for a sec\n")
		time.Sleep(time.Second)

	}

	for k := range procMap {
		log.Println("stopping the redis server")
		procMap[k].Stop()
	}

	for {

		log.Printf("sleeping forever\n")
		time.Sleep(time.Second)

	}
}
