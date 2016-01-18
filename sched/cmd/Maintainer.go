package cmd

import (
	"log"

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

		switch ts.State {

		case "TASK_STAGING":
			log.Printf("Task %s is Staging", ts.Name)
			break
		case "TASK_STARTING":
			log.Printf("Task %s is Starting", ts.Name)
			break
		case "TASK_RUNNING":
			log.Printf("Task %s is Running", ts.Name)
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
