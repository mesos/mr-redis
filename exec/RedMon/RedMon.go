package RedMon

import (
	"fmt"
	"io"
	//	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	//Redis Client
	typ "../../common/types"
)

//This structure is used to implement a monitor thread/goroutine for a running Proc(redisProc)
//This structure should be extended only if more functionality is required on the Monitoring functionality
//A Redis Proc's objec is created within this and monitored hence forth
type RedMon struct {
	P       *typ.Proc //The Proc structure that should be used
	Pid     int       //The Pid of the running proc
	Port    int       //The port number of this redis instance to be started
	Ofile   io.Writer //Stdout log file to be re-directed to this io.writer
	Efile   io.Writer //stderr of the redis instnace should be re-directed to this file
	MS_Sync bool      //Make this as master after sync
	C       *exec.Cmd
	//Cli *Redis.Cli 		//redis cli client library connection handler
	//cgroup *CgroupManager		//Cgroup manager/cgroup connection pointer
}

//Create a new monitor based on the Data sent along with the TaskInfo
//The data could have the following details
//Capacity Master                 => Just start this PROC send update as TASK_RUNNING and monitor henceforth
//Capacity SlaveOf IP:Port        => This is a redis slave so start it as a slave, sync and then send TASK_RUNNING update then Monitor
//Capacity Master-SlaveOf IP:Port => This is a New master of the instance with an upgraded memory value so
//                          Start as slave, Sync data, make it as master, send TASK_RUNNING update and start to Monitor

func NewRedMon(tskName string, Type string, data string) *RedMon {

	var R RedMon
	var P *typ.Proc

	split_data := strings.Split(data, " ")
	if len(split_data) < 1 || len(split_data) > 4 {
		//Print an error this is not suppose to happen
		return nil
	}

	Cap, _ := strconv.Atoi(split_data[0])

	switch split_data[1] {
	case "Master":
		P = typ.NewProc(tskName, Cap, Type, "")
		break
	case "SlaveOf":
		P = typ.NewProc(tskName, Cap, Type, split_data[2])
		break
	case "Master-SlaveOf":
		P = typ.NewProc(tskName, Cap, Type, split_data[2])
		R.MS_Sync = true
		break
	}
	R.P = P
	//ToDo Stdout file to be tskname.stdout
	//ToDo stderere file to be tskname.stderr

	return &R
}

//Start the redis Proc be it Master or Slave
func (R *RedMon) Start() bool {

	if R.P.SlaveOf == "" {
		return R.StartMaster()
	} else {

		if !R.MS_Sync {
			return R.StartSlave()
		} else {
			//Posibly a scale request so start it as a slave, sync then make as master
			return R.StartSlaveAndMakeMaster()
		}
	}
	return false
}

func (R *RedMon) StartMaster() bool {
	//Command Line
	R.C = exec.Command("/home/ubuntu/progs/redis-3.0.6/src/redis-server", "--port", fmt.Sprintf("%d", R.Port))
	err := R.C.Start()

	if err != nil {
		//Print some error
		return false
	}

	R.Pid = R.C.Process.Pid
	R.P.Sync()

	return true
}

func (R *RedMon) StartSlave() bool {
	//Command Line
	R.C = exec.Command("/home/ubuntu/progs/redis-3.0.6/src/redis-server", "--port", fmt.Sprintf("%d", R.Port), "--SlaveOf", R.P.SlaveOf)
	err := R.C.Start()

	if err != nil {
		//Print some error
		return false
	}

	R.Pid = R.C.Process.Pid
	//Monitor the redis PROC to check if the sync is complete
	for !R.IsSyncComplete() {
		time.Sleep(time.Second)
	}

	R.P.Sync()

	return true
}

func (R *RedMon) StartSlaveAndMakeMaster() bool {
	//Command Line
	R.C = exec.Command("/home/ubuntu/progs/redis-3.0.6/src/redis-server", "--port", fmt.Sprintf("%d", R.Port), "--SlaveOf", R.P.SlaveOf)
	err := R.C.Start()

	if err != nil {
		//Print some error
		return false
	}

	R.Pid = R.C.Process.Pid

	//Monitor the redis PROC to check if the sync is complete
	for !R.IsSyncComplete() {
		time.Sleep(time.Second)
	}
	//Make this Proc as master
	R.MakeMaster()

	R.P.Sync()

	return true
}

func (R *RedMon) StatsUpdate() bool {
	//Contact the redis instace
	//collecgt the stats

	//sync

	R.P.SyncStats()
	return true
}

func (R *RedMon) Die() bool {
	//send SHUTDOWN cli command for a gracefull exit of the redis-server
	return true
}

//Should be called by the Monitors on Slave Procs, this gives the boolien anser if the sync is complegted or not
func (R *RedMon) IsSyncComplete() bool {

	//Keep checking if the sync of data freom master is completed or not
	return true
}

func (R *RedMon) MakeMaster() bool {

	//Send a cli config comamnd to make a current Proc a master from a slave
	return true
}
