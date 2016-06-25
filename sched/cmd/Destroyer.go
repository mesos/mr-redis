package cmd

import (
	"fmt"
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
		select {

		case msg := <-typ.Dchan:
			switch msg.MSG {
			case typ.TASK_MSG_DESTROY:
				msg.P.Msg = "SHUTDOWN"
			case typ.TASK_MSG_MAKEMASTER:
				msg.P.Msg = "MASTER"
				msg.P.SyncSlaveOf()
			case typ.TASK_MSG_SLAVEOF:
				msg.P.Msg = fmt.Sprintf("SLAVEOF %s", msg.P.SlaveOf)
				msg.P.SyncSlaveOf()
			}
			msg.P.SyncMsg()
			log.Printf("Destorying proc %v from Instance %v", msg.P.ID, msg.P.Instance)
			break
		}
	}

}
