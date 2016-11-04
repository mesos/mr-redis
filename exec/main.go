package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	exec "github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"

	typ "github.com/mesos/mr-redis/common/types"
	"github.com/mesos/mr-redis/exec/RedMon"
)

//DbType Flag for dbtype like etcd/zookeeper
var DbType = flag.String("DbType", "etcd", "Type of the database etcd/zookeeper etc.,")

//DbEndPoint The actuall endpoint of the database.
var DbEndPoint = flag.String("DbEndPoint", "", "Endpoint of the database")

var Image = flag.String("Image", "redis:3.0-alpine", "Image of the Redis Proc to be downloaded")

//MrRedisLogger A global Logger pointer for the executor all the RedMon will write to the same logger
var MrRedisLogger *log.Logger

//MrRedisExecutor Basic strucutre for the executor
type MrRedisExecutor struct {
	tasksLaunched int
	HostIP        string
	ExecutorID    string
	monMap        map[string](*RedMon.RedMon)
}

//GetLocalIP A function to look up the exposed local IP such that the executor can bind to
func GetLocalIP() string {

	if libprocessIP := os.Getenv("LIBPROCESS_IP"); libprocessIP != "" {
		address := net.ParseIP(libprocessIP)
		if address != nil {
			//If its a valid IP address return the string
			fmt.Printf("LibProess IP = %s", libprocessIP)
			return libprocessIP
		}

	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Printf("InterfaceAddress = %s", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

//NewMrRedisExecutor Constructor for the executor structure
func NewMrRedisExecutor() *MrRedisExecutor {
	return &MrRedisExecutor{tasksLaunched: 0}
}

//Registered Call back for registered driver
func (exec *MrRedisExecutor) Registered(driver exec.ExecutorDriver, execInfo *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo, slaveInfo *mesos.SlaveInfo) {
	exec.ExecutorID = execInfo.ExecutorId.GetValue()
	fmt.Println("Registered Executor on slave ") //, slaveInfo.GetHostname())
}

//Reregistered call back for the re-registered driver
func (exec *MrRedisExecutor) Reregistered(driver exec.ExecutorDriver, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Re-registered Executor on slave ") //, slaveInfo.GetHostname())
}

//Disconnected Call back for disconnected
func (exec *MrRedisExecutor) Disconnected(exec.ExecutorDriver) {
	fmt.Println("Executor disconnected.")
}

//LaunchTask Call back implementation when a Launch task request comes from Slave/Agent
func (exec *MrRedisExecutor) LaunchTask(driver exec.ExecutorDriver, taskInfo *mesos.TaskInfo) {
	fmt.Println("Launching task", taskInfo.GetName(), "with command", taskInfo.Command.GetValue())

	var runStatus *mesos.TaskStatus
	exec.tasksLaunched++
	M := RedMon.NewRedMon(taskInfo.GetTaskId().GetValue(), exec.HostIP, exec.tasksLaunched+6379, string(taskInfo.Data), MrRedisLogger, *Image)

	fmt.Printf("The Redmon object = %v\n", *M)

	tid := taskInfo.GetTaskId().GetValue()
	exec.monMap[tid] = M

	go func() {
		if M.Start() {
			runStatus = &mesos.TaskStatus{
				TaskId: taskInfo.GetTaskId(),
				State:  mesos.TaskState_TASK_RUNNING.Enum(),
			}
		} else {
			runStatus = &mesos.TaskStatus{
				TaskId: taskInfo.GetTaskId(),
				State:  mesos.TaskState_TASK_ERROR.Enum(),
			}
		}
		_, err := driver.SendStatusUpdate(runStatus)
		if err != nil {
			fmt.Println("Got error", err)
		}

		fmt.Println("Total tasks launched ", exec.tasksLaunched)

		//our server is now running, lets start monitoring it also
		go func() {
			M.Monitor()
		}()

		exitState := mesos.TaskState_TASK_FINISHED.Enum()

		exitErr := M.Container.Wait() //TODO: Collect the return value of the process and send appropriate TaskUpdate eg:TaskFinished only on clean shutdown others will get TaskFailed
		if exitErr != 0 || M.P.Msg != "SHUTDOWN" {
			//If the redis-server proc finished either with a non-zero or its not suppose to die then mark it as Task filed
			exitState = mesos.TaskState_TASK_FAILED.Enum()
			//Signal the monitoring thread to stop monitoring from now on
			M.MonChan <- 1
		}

		// finish task
		fmt.Println("Finishing task", taskInfo.GetName())
		finStatus := &mesos.TaskStatus{
			TaskId: taskInfo.GetTaskId(),
			State:  exitState,
		}
		_, err = driver.SendStatusUpdate(finStatus)
		if err != nil {
			fmt.Println("Got error", err)
		}
		fmt.Println("Task finished", taskInfo.GetName())
	}()
}

//KillTask When a running task needs to be killed should come from the Agent/Slave its a call back implementation
func (exec *MrRedisExecutor) KillTask(driver exec.ExecutorDriver, taskID *mesos.TaskID) {
	tid := taskID.GetValue()
	//tbd: is there any error check needed
	exec.monMap[tid].Die()

	fmt.Println("Killed task with task id:", tid)
}

//FrameworkMessage Any message sent from the scheduelr , not sued for this project
func (exec *MrRedisExecutor) FrameworkMessage(driver exec.ExecutorDriver, msg string) {
	fmt.Println("Got framework message: ", msg)
}

//Shutdown Not implemented yet
func (exec *MrRedisExecutor) Shutdown(exec.ExecutorDriver) {
	fmt.Println("Shutting down the executor")
	fmt.Printf("Killing all the containers")
}

//Error not implemented yet
func (exec *MrRedisExecutor) Error(driver exec.ExecutorDriver, err string) {
	fmt.Println("Got error message:", err)
}

// -------------------------- func inits () ----------------- //
func init() {
	flag.Parse()
}

func main() {
	fmt.Println("Starting MrRedis Executor")

	typ.Initialize(*DbType, *DbEndPoint)

	var out io.Writer
	out = ioutil.Discard

	out, _ = os.Create("/tmp/MrRedisExecutor.log")
	//ToDo does this need error handling
	MrRedisLogger = log.New(out, "[Info]", log.Lshortfile)

	MrRedisExec := NewMrRedisExecutor()
	MrRedisExec.HostIP = GetLocalIP()
	MrRedisExec.monMap = make(map[string](*RedMon.RedMon))

	dconfig := exec.DriverConfig{
		BindingAddress: net.ParseIP(MrRedisExec.HostIP),
		Executor:       MrRedisExec,
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
	fmt.Println("Executor Finished, Delete all the containers")
	for _, M := range MrRedisExec.monMap {
		M.Die()
	}
	fmt.Println("executor terminated")
}
