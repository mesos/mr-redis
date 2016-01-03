package serviceproc

import (
	"fmt"
	redisclient "gopkg.in/redis.v3"
	"log"
	"time"
)

type ProcStats struct {
	Mem int
	Cpu int
	BW  int //there are way too many stats which can be picked from redis client info command
}

type ProcMonitor struct {
	ID        string //monitor id
	proc      *RedisProc
	redClient *redisclient.Client
}

func NewProcMonitor(proc *RedisProc) *ProcMonitor {
	return &ProcMonitor{ID: "", proc: proc}
}

func (pm *ProcMonitor) monitorStats( /*connected client details*/ ) {

	_, err := pm.redClient.Info().Result()
	if err != nil {
		log.Printf("error:", err)
	}
	log.Printf("MONITORING STATS")

	//log.Printf(info)
	//read related stats from redis server in question
}

func (pm *ProcMonitor) GetConnectedClient() *redisclient.Client {

	log.Printf("Monitoring stats")

	client := redisclient.NewClient(&redisclient.Options{
		Addr:     "localhost:" + fmt.Sprintf("%d", pm.proc.Portno),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Printf(pong, err)
	//read related stats from redis server in question
	return client
}

func (pm *ProcMonitor) SpawnandMonitor() error {

	var waitChan chan error
	var waitErr error

	defer pm.proc.Stop()

	//goroutine code here
	go func() {

		cmdErr := pm.proc.Start(pm.proc.Portno)
		log.Printf("Process finsihed error=%v", cmdErr)
		waitChan <- cmdErr

	}()

	//wait for a second for the server to start
	time.Sleep(1 * time.Second)
	//then initiate a connection to it; for stats
	pm.redClient = pm.GetConnectedClient()

	for {
		select {

		case waitErr = <-waitChan:
			if waitErr != nil {
				log.Printf("Instance %s failed with error %v", pm.proc.ID, waitErr)
			}
			//tbd: other things like state update if any
			return waitErr

		case <-time.After(1 * time.Second):
			pm.monitorStats()
		}

	}
}

func (pm *ProcMonitor) Stop() error {

	//tbd:disconnect from redis server
	//tbd:stop the monitored server
	return nil

}
