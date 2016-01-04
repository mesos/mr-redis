package main

//initial helper for run the executor functionality from the command line without mesos
//the same (or similar :)) code will go into mesos-go executor launchtask function

import (
	"../common/store/etcd"
	//	"../common/types"
	"./serviceproc"
	"log"
	"time"
	"fmt"
	"github.com/mesos/mesos-go/executor"
)

var procMap map[string](*serviceproc.RedisProc)
var driver executor.ExecutorDriver

func main() {
	log.Printf("Starting Redis server on given port\n")

	procMap = make(map[string](*serviceproc.RedisProc))

	serviceproc.Store = etcd.New()

	/* Setup etcd with the etcd endpoint*/
	err := serviceproc.Store.Setup("http://127.0.0.1:2379")
	/* Test if this is setup */
	log.Printf("IsSetup %v\n with error :%s", serviceproc.Store.IsSetup(), err)

	//tbd: only the service instance id needs to be passed here; how to get it?
	//tbd: who gets you the port value?
	for i := 0; i < 2; i++ {
		redisproc := serviceproc.NewRedisProc("ServiceInstnsID", (6379 + i), fmt.Sprintf("%d", i))
		procMap[fmt.Sprintf("%d", i)] = redisproc
		log.Printf("spawning a new server with id:%d and proc:%v", i, redisproc)
		monitor := serviceproc.NewProcMonitor(redisproc, &driver)

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
		err = procMap[k].GetfromStore()
		if err != nil {
			log.Printf("error getting the expected value from store\n")
			log.Println(err)
		}
		log.Println("stopping the redis server")
		procMap[k].Stop()
	}

	for {

		log.Printf("sleeping forever\n")
		time.Sleep(time.Second)

	}
}
