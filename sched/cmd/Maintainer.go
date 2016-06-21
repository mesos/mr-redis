package cmd

import (
	"log"
	"strings"

	typ "github.com/mesos/mr-redis/common/types"
)

//This is the main function that handles all the task updates
func Maintainer() {

	log.Printf("Scheduler Maintainer is starting")

	var ts *typ.TaskUpdate

	for {

		select {

		case ts = <-typ.Mchan:
			log.Printf("Received a Task update from the channel %v", ts)
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

		if Inst.Procs == nil {
			Inst.LoadProcs()
		}

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
				if Inst.Masters < Inst.ExpMasters {
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
				typ.Cchan <- typ.CreateSlaves(Inst, Inst.ExpSlaves)
				break
			case "S":
				if Inst.Slaves < Inst.ExpSlaves {
					Inst.Slaves++
					Inst.Snames = append(Inst.Snames, ProcID)
					Inst.SyncSlaves()
					Inst.Procs[proc.ID] = proc
				} else {
					log.Printf("Unknown Slave %v  created for this instance", ts.Name)
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
			switch proc.Type {
			case "M":
				if Inst.Masters > 0 {
					Inst.Masters--
					Inst.SyncMasters()

				}

			case "S":
				if Inst.Slaves > 0 {
					Inst.Slaves--
					//Remove this lsave from the list of slaves
					Inst.SyncSlaves()
				}
			}
			if Inst.Slaves == 0 && Inst.Masters == 0 {
				Inst.Status = typ.INST_STATUS_DELETED
				Inst.SyncStatus()
			}
			break
		case "TASK_FAILED":
			log.Printf("Task %s is Failed", ts.Name)
			switch proc.Type {
			case "M":
				//If the task lost is a master then we must select a most updated slave as the next master
				//Make rest of the slave to start connectin to this new master
				//Send Request to creator to bring back one more slave
				//For now lets just start a master for single instance master
				if Inst.Type == typ.INST_TYPE_SINGLE {
					if Inst.Masters > 0 {
						Inst.Masters--
						Inst.Mname = ""
						Inst.SyncMasters()
						typ.Cchan <- typ.CreateMaster(Inst)
					}
				}
				break
			case "S":
				//Just send requst to bring back one more slave to the creator
				if Inst.Slaves > 0 {
					Inst.Slaves--
					//Remove this lsave from the list of slaves
					var tmp_Snames []string
					for _, pid := range Inst.Snames {
						if pid != ProcID {
							tmp_Snames = append(tmp_Snames, pid)
						}
					}
					Inst.Snames = tmp_Snames
					Inst.SyncSlaves()
					typ.Cchan <- typ.CreateSlaves(Inst, 1)
				}
				break
			}
			break
		case "TASK_KILLED":
			log.Printf("Task %s is Killed", ts.Name)
			break
		case "TASK_LOST":
			log.Printf("Task %s is Lost", ts.Name)
			switch proc.Type {
			case "M":
				//If the task lost is a master then we must select a most updated slave as the next master
				//Make rest of the slave to start connectin to this new master
				//Send Request to creator to bring back one more slave
				//For now lets just start a master for single instance master
				if Inst.Type == typ.INST_TYPE_SINGLE {
					if Inst.Masters > 0 {
						Inst.Masters--
						Inst.Mname = ""
						Inst.SyncMasters()
						typ.Cchan <- typ.CreateMaster(Inst)
					}
				} else {  //The Master has died in  a Master Slave Setup

					//We need to Elect a new master among the slaves, below are the steps we need to perform
					//1) Loop through the slaves and select the one with MAX slave_repl_offset 

				}
				break
			case "S":
				//Just send requst to bring back one more slave to the creator
				if Inst.Slaves > 0 {
					Inst.Slaves--
					//Remove this lsave from the list of slaves
					var tmp_Snames []string
					for _, pid := range Inst.Snames {
						if pid != ProcID {
							tmp_Snames = append(tmp_Snames, pid)
						}
					}
					Inst.Snames = tmp_Snames
					Inst.SyncSlaves()
					typ.Cchan <- typ.CreateSlaves(Inst, 1)
				}
				break
			}
			break
		case "TASK_ERROR":
			log.Printf("Task %s is Error", ts.Name)
			break
		}
	}

	log.Printf("Scheduler Maintainer is stopped")

}
