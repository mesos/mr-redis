package cmd

import (
	"log"
	"strings"

	typ "../../common/types"
)

//This is the main function that handles all the task updates
func Maintainer() {

	log.Printf("Scheduler Maintainer is startring")

	var ts *typ.TaskUpdate

	for {

		select {

		case ts = <-typ.Mchan:
			log.Printf("Recived a Task update from the channel %v", ts)
			break

		}

		tnames := strings.SplitN(ts.Name, "::", 2)
		if len(tnames) != 2 {
			log.Printf("Bad task name formate %v", tnames)
			continue
		}
		InstName := tnames[0]
		ProcID := tnames[1]

		//Check in the memory if there is such an instance runnig

		Inst := typ.MemDb.Get(InstName)

		if Inst == nil {
			Inst = typ.LoadInstance(InstName)
			if Inst == nil {
				log.Printf("No such Task(%v) in our records, Ignoring", ts)
				continue
			} else {
				typ.MemDb.Add(InstName, Inst)
			}
		}

		proc := typ.LoadProc(ts.Name)

		switch ts.State {

		case "TASK_STAGING":
			log.Printf("Task %s is Staging", ts.Name)
			break
		case "TASK_STARTING":
			log.Printf("Task %s is Starting", ts.Name)
			break
		case "TASK_RUNNING":
			log.Printf("Task %s is Running", ts.Name)
			switch proc.Type {
			case "M":
				if Inst.Masters <= Inst.ExpMasters {
					//
					Inst.Masters++
				} else {
					//Now this means that we have master task when we already have all our masters running
					//This could mean that a new master is available with a scaled capacity
					//OldMaster := Inst.Mname
					//ToDo Send old master id to the Destoryer

					//Mark all the old slave to be deleted send the slave id to the destroyer
					Inst.Slaves = 0
					Inst.SyncSlaves()

				}
				Inst.Mname = proc.ID
				Inst.SyncMasters()
				Inst.Procs[proc.ID] = proc
				//Send the instance detail to Create so that slaves can be created now
				typ.Cchan <- Inst
				break
			case "S":
				if Inst.Slaves <= Inst.ExpSlaves {
					Inst.Slaves++
					Inst.Snames = append(Inst.Snames, ProcID)
					Inst.SyncSlaves()
				} else {
					log.Printf("Unknown Slave %v  created for this instnace", ts.Name)
				}
				break
			}
			if Inst.Masters == Inst.ExpMasters && Inst.Slaves == Inst.ExpSlaves {
				Inst.Status = typ.INST_STATUS_RUNNING
				Inst.SyncStatus()
			}
			break
		case "TASK_FINISHED":
			log.Printf("Task %s is Finished", ts.Name)
			break
		case "TASK_FAILED":
			log.Printf("Task %s is Failed", ts.Name)
			break
		case "TASK_KILLED":
			log.Printf("Task %s is Killed", ts.Name)
			break
		case "TASK_LOST":
			log.Printf("Task %s is Lost", ts.Name)
			break
		case "TASK_ERROR":
			log.Printf("Task %s is Error", ts.Name)
			break

		}
	}

	log.Printf("Scheduler Maintainer is stopped")

}
