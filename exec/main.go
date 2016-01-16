package main

import (
	"flag"
	"fmt"
	exe "os/exec"

	exec "github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"

	typ "../common/types"
)

type MrRedisExecutor struct {
	tasksLaunched int
}

func NewMrRedisExecutor() *MrRedisExecutor {
	return &MrRedisExecutor{tasksLaunched: 0}
}

func (exec *MrRedisExecutor) Registered(driver exec.ExecutorDriver, execInfo *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *MrRedisExecutor) Reregistered(driver exec.ExecutorDriver, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Re-registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *MrRedisExecutor) Disconnected(exec.ExecutorDriver) {
	fmt.Println("Executor disconnected.")
}

func (exec *MrRedisExecutor) LaunchTask(driver exec.ExecutorDriver, taskInfo *mesos.TaskInfo) {
	fmt.Println("Launching task", taskInfo.GetName(), "with command", taskInfo.Command.GetValue())

	cmd := exe.Command("/home/ubuntu/progs/redis-3.0.6/src/redis-server", "--port", fmt.Sprintf("%d", 6379+exec.tasksLaunched))
	err := cmd.Start()
	if err != nil {
		return
	}
	runStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_RUNNING.Enum(),
	}
	_, err = driver.SendStatusUpdate(runStatus)
	if err != nil {
		fmt.Println("Got error", err)
	}

	exec.tasksLaunched++
	fmt.Println("Total tasks launched ", exec.tasksLaunched)
	//
	// this is where one would perform the requested task
	//

	go func() {
		cmd.Wait()
		// finish task
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
	}()
}

func (exec *MrRedisExecutor) KillTask(exec.ExecutorDriver, *mesos.TaskID) {
	fmt.Println("Kill task")
}

func (exec *MrRedisExecutor) FrameworkMessage(driver exec.ExecutorDriver, msg string) {
	fmt.Println("Got framework message: ", msg)
}

func (exec *MrRedisExecutor) Shutdown(exec.ExecutorDriver) {
	fmt.Println("Shutting down the executor")
}

func (exec *MrRedisExecutor) Error(driver exec.ExecutorDriver, err string) {
	fmt.Println("Got error message:", err)
}

// -------------------------- func inits () ----------------- //
func init() {
	flag.Parse()
}

func main() {
	fmt.Println("Starting MrRedis Executor")

	DbType := flag.String("DbType", "etcd", "Type of the database etcd/zookeeper etc.,")
	DbEndPoint := flag.String("DbEndPoint", "", "Endpoint of the database")

	typ.Initialize(*DbType, *DbEndPoint)
	dconfig := exec.DriverConfig{
		Executor: NewMrRedisExecutor(),
	}
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
