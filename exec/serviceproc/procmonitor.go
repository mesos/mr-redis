package serviceproc

import (
	"log"
	"time"
)

type ProcStats struct {
	Mem int
	Cpu int
	BW  int
}

type ProcMonitor struct {
	ID   string //monitor id
	proc *RedisProc
}

func NewProcMonitor(proc *RedisProc) *ProcMonitor {
	return &ProcMonitor{ID: "", proc: proc}
}

func monitorStats( /*connected client details*/ ) {

	//read related stats from redis server in question
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

	for {
		select {

		case waitErr = <-waitChan:
			if waitErr != nil {
				log.Printf("Instance %s failed with error %v", pm.proc.ID, waitErr)
			}
			//tbd: other things like state update if any
			return waitErr

		case <-time.After(1 * time.Second):
			monitorStats()
		}

	}
}

func (pm *ProcMonitor) Stop() error {

	//tbd:disconnect from redis server
	//tbd:stop the monitored server
	return nil

}
