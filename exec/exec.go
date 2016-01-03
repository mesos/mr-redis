/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"

	"./serviceproc"
	exec "github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

//tbd: the monitoring needs to be improved and should be monitoring the launched
// process pid in the system, and the launched process should not exit
// in case there is a problem with the executor and it exits
// this also brings in the need of init-executor which will load the info
//about the existing processes on this slave

var procMap map[string](*serviceproc.RedisProc)

type exampleExecutor struct {
	tasksLaunched int
}

func newExampleExecutor() *exampleExecutor {
	return &exampleExecutor{tasksLaunched: 0}
}

func (exec *exampleExecutor) Registered(driver exec.ExecutorDriver, execInfo *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *exampleExecutor) Reregistered(driver exec.ExecutorDriver, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Re-registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *exampleExecutor) Disconnected(exec.ExecutorDriver) {
	fmt.Println("Executor disconnected.")
}

func (exec *exampleExecutor) LaunchTask(driver exec.ExecutorDriver, taskInfo *mesos.TaskInfo) {
	fmt.Println("Launching task", taskInfo.GetName(), "with command", taskInfo.Command.GetValue())

	runStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_RUNNING.Enum(),
	}
	_, err := driver.SendStatusUpdate(runStatus)
	if err != nil {
		fmt.Println("Got error", err)
	}

	exec.tasksLaunched++
	fmt.Println("Total tasks launched ", exec.tasksLaunched)

	// this is where one would perform the requested task
	fmt.Println("Starting Redis server on given port\n")

	//tbd: only the service instance id needs to be passed here; how to get it?
	//tbd: who gets you the port value?
	redisproc, uidStr := serviceproc.NewRedisProc("ServiceInstnsID", (6379 + exec.tasksLaunched))
	procMap[uidStr] = redisproc
	fmt.Println("spawning a new server with id:%s and proc:%v", uidStr, redisproc)
	monitor := serviceproc.NewProcMonitor(redisproc)

	//tbd: this needs to be tested by invoking multiple servers on the same executor
	//if the same executor is used for all launchtasks, we cannot block here
	//and thus will need a goroutine; (but then why does example put task finished status here)
	monitor.SpawnandMonitor()
	// finished the task (success or failure) if returned from the above

	fmt.Println("Finishing task", taskInfo.GetName())
	finStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_FINISHED.Enum(),
	}
	_, err = driver.SendStatusUpdate(finStatus)
	if err != nil {
		fmt.Println("Got error", err)
	}
	fmt.Println("Task finished", taskInfo.GetName())
}

func (exec *exampleExecutor) KillTask(exec.ExecutorDriver, *mesos.TaskID) {
	fmt.Println("Kill task")
}

func (exec *exampleExecutor) FrameworkMessage(driver exec.ExecutorDriver, msg string) {
	fmt.Println("Got framework message: ", msg)
}

func (exec *exampleExecutor) Shutdown(exec.ExecutorDriver) {
	fmt.Println("Shutting down the executor")
}

func (exec *exampleExecutor) Error(driver exec.ExecutorDriver, err string) {
	fmt.Println("Got error message:", err)
}

// -------------------------- func inits () ----------------- //
func init() {
	flag.Parse()
}

func main() {
	fmt.Println("Starting Example Executor (Go)")

	dconfig := exec.DriverConfig{
		Executor: newExampleExecutor(),
	}

	//maps which stores the proc pointer wrt UID
	procMap = make(map[string](*serviceproc.RedisProc))

	//tbd: after new, an init on the executor has to be called
	//which will read info about existing redis processes
	//also need to check how mesos handles executor errors

	driver, err := exec.NewMesosExecutorDriver(dconfig)

	if err != nil {
		fmt.Println("Unable to create a ExecutorDriver ", err.Error())
	}

	_, err = driver.Start()
	if err != nil {
		fmt.Println("Got error:", err)
		return
	}
	fmt.Println("Executor process has started and running.")
	_, err = driver.Join()
	if err != nil {
		fmt.Println("driver failed:", err)
	}
	fmt.Println("executor terminating")
}
