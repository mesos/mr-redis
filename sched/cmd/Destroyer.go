package cmd

import (
	"log"

	typ "github.com/mesos/mr-redis/common/types"
)

//This shouldbe started as a goroutine as it runs unconditionally
//This listens on the Dchan declared globally, A pointer to Proc variable is sent across through this channel
//The only job of this goroutine is to update MSG structure of the Proc entry
//The MSG will be read by RedMon (Redis Monitoring goroutine) and perform the action
//The primary duty of this goroutine is to initate send SHUTDOWN of a redis PROC
func Destoryer() bool {

	for {
		var proc *typ.Proc
		select {

		case proc = <-typ.Dchan:
			proc.Msg = "SHUTDOWN"
			proc.SyncMsg()
			log.Printf("Destorying proc %v from Instance %v", proc.ID, proc.Instance)
			break
		}

	}
}
