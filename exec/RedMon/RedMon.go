package RedMon

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	typ "github.com/mesos/mr-redis/common/types"
	redisclient "gopkg.in/redis.v3"
)

//This structure is used to implement a monitor thread/goroutine for a running Proc(redisProc)
//This structure should be extended only if more functionality is required on the Monitoring functionality
//A Redis Proc's objec is created within this and monitored hence forth
type RedMon struct {
	P       *typ.Proc //The Proc structure that should be used
	Pid     int       //The Pid of the running proc
	IP      string    //IP address the redis instance should bind to
	Port    int       //The port number of this redis instance to be started
	Ofile   io.Writer //Stdout log file to be re-directed to this io.writer
	Efile   io.Writer //stderr of the redis instnace should be re-directed to this file
	MS_Sync bool      //Make this as master after sync
	monChan chan int
	Cmd     *exec.Cmd
	Client  *redisclient.Client //redis client library connection handler
	L       *log.Logger         //to redirect log outputs to a file
	//cgroup *CgroupManager		//Cgroup manager/cgroup connection pointer
}

//Create a new monitor based on the Data sent along with the TaskInfo
//The data could have the following details
//Capacity Master                 => Just start this PROC send update as TASK_RUNNING and monitor henceforth
//Capacity SlaveOf IP:Port        => This is a redis slave so start it as a slave, sync and then send TASK_RUNNING update then Monitor
//Capacity Master-SlaveOf IP:Port => This is a New master of the instance with an upgraded memory value so
//                          Start as slave, Sync data, make it as master, send TASK_RUNNING update and start to Monitor
func NewRedMon(tskName string, IP string, Port int, data string) *RedMon {

	var R RedMon
	var P *typ.Proc
	var out io.Writer = ioutil.Discard

	R.monChan = make(chan int)
	R.Port = Port
	R.IP = IP

	//initialise logger also
	out, _ = os.Create("/tmp/MrRedisExecutor.log")
	//ToDo does this need error handling
	R.L = log.New(out, "[Info]", log.Lshortfile)

	R.L.Printf("Split data recived is %v\n", data)

	split_data := strings.Split(data, " ")
	if len(split_data) < 1 || len(split_data) > 4 {
		//Print an error this is not suppose to happen
		R.L.Printf("RedMon Splitdata error %v\n", split_data)
		return nil
	}

	Cap, _ := strconv.Atoi(split_data[0])

	switch split_data[1] {
	case "Master":
		P = typ.NewProc(tskName, Cap, "M", "")
		R.L.Printf("created proc for new MASTER\n")
		break
	case "SlaveOf":
		P = typ.NewProc(tskName, Cap, "S", split_data[2])
		break
	case "Master-SlaveOf":
		P = typ.NewProc(tskName, Cap, "MS", split_data[2])
		R.MS_Sync = true
		break
	}
	R.P = P
	//ToDo each instance should be started with its own dir and specified config file
	//ToDo Stdout file to be tskname.stdout
	//ToDo stderere file to be tskname.stderr

	return &R
}

func (R *RedMon) getConnectedClient() *redisclient.Client {

	R.L.Printf("Monitoring stats")

	client := redisclient.NewClient(&redisclient.Options{
		Addr:     R.IP + ":" + fmt.Sprintf("%d", R.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	R.L.Printf(pong, err)

	return client
}

func (R *RedMon) launchRedisServer(isSlave bool, IP string, port string) bool {

	if isSlave {
		R.Cmd = exec.Command("./redis-server", "--port", fmt.Sprintf("%d", R.Port), "--SlaveOf", IP, port)
	} else {
		R.Cmd = exec.Command("./redis-server", "--port", fmt.Sprintf("%d", R.Port))
	}

	err := R.Cmd.Start()

	if err != nil {
		//Print some error
		return false
	}

	//hack otherwise its too quick to have the server receiving connections
	time.Sleep(time.Second)

	//get the connected client immediately after for monitoring and other functions
	R.Client = R.getConnectedClient()
	return true
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

	var ret = false
	//Command Line
	ret = R.launchRedisServer(false, "", "")
	if ret != true {
		return ret
	}

	R.Pid = R.Cmd.Process.Pid
	R.P.Pid = R.Cmd.Process.Pid
	R.P.Port = fmt.Sprintf("%d", R.Port)
	R.P.IP = R.IP
	R.P.State = "Running"
	R.P.Sync()

	return true
}

func (R *RedMon) StartSlave() bool {
	var ret = false
	//Command Line
	slaveof := strings.Split(R.P.SlaveOf, ":")
	if len(slaveof) != 2 {
		R.L.Printf("Unacceptable SlaveOf value %vn", R.P.SlaveOf)
		return false
	}

	//Command Line
	ret = R.launchRedisServer(true, slaveof[0], slaveof[1])
	if ret != true {
		return ret
	}

	//Monitor the redis PROC to check if the sync is complete
	for !R.IsSyncComplete() {
		time.Sleep(time.Second)
	}
	R.Pid = R.Cmd.Process.Pid
	R.P.Pid = R.Cmd.Process.Pid
	R.P.Port = fmt.Sprintf("%d", R.Port)
	R.P.IP = R.IP
	R.P.State = "Running"

	R.P.Sync()

	return true
}

func (R *RedMon) StartSlaveAndMakeMaster() bool {
	var ret = false
	//Command Line
	slaveof := strings.Split(R.P.SlaveOf, ":")
	if len(slaveof) != 2 {
		R.L.Printf("Unacceptable SlaveOf value %vn", R.P.SlaveOf)
		return false
	}

	ret = R.launchRedisServer(true, slaveof[0], slaveof[1])
	if ret != true {
		return ret
	}

	R.Pid = R.Cmd.Process.Pid

	//Monitor the redis PROC to check if the sync is complete
	for !R.IsSyncComplete() {
		time.Sleep(time.Second)
	}
	//Make this Proc as master
	R.MakeMaster()

	R.Pid = R.Cmd.Process.Pid
	R.P.Pid = R.Cmd.Process.Pid
	R.P.Port = fmt.Sprintf("%d", R.Port)
	R.P.IP = R.IP
	R.P.State = "Running"
	R.P.Sync()

	return true
}

func (R *RedMon) UpdateStats() bool {

	var redisStats typ.Stats
	var err error

	redisStats.Mem, err = R.Client.Info("memory").Result()
	if err != nil {
		R.L.Printf("STATS collection returned error on IP:%s and PORT:%d Err:%v", R.IP, R.Port, err)
		return false
	}

	redisStats.Cpu, err = R.Client.Info("cpu").Result()
	if err != nil {
		R.L.Printf("STATS collection returned error on IP:%s and PORT:%d Err:%v", R.IP, R.Port, err)
		return false
	}

	redisStats.Others, err = R.Client.Info("stats").Result()
	if err != nil {
		R.L.Printf("STATS collection returned error on IP:%s and PORT:%d Err:%v", R.IP, R.Port, err)
		return false
	}

	R.P.Stats = R.P.ToJsonStats(redisStats)

	errSync := R.P.SyncStats()
	if !errSync {
		R.L.Printf("Error syncing stats to store")
		return false
	}
	return true
}

func (R *RedMon) Monitor() bool {

	//wait for a second for the server to start
	//ToDo: is it needed
	time.Sleep(1 * time.Second)

	for {
		select {

		case <-R.monChan:
			//ToDo:update state if needed
			//signal to stop monitoring this
			return false

		case <-time.After(time.Second):
			//this is to check communication from scheduler; mesos messages are not reliable
			R.CheckMsg()

		case <-time.After(1 * time.Second):
			R.UpdateStats()
		}

	}

}

func (R *RedMon) Stop() bool {

	//send SHUTDOWN command for a gracefull exit of the redis-server
	//the server exited gracefully will reflect at the task status FINISHED
	_, err := R.Client.Shutdown().Result()
	if err != nil {
		R.L.Printf("problem shutting down the server at IP:%s and port:%d with error %v", R.IP, R.Port, err)

		//in this error case the scheduler will get a task killed notification
		//but will also see that the status it updated was SHUTDOWN, thus will handle it as OK

		errMsg := R.Die()
		if !errMsg { //message should be read by scheduler
			R.L.Printf("Killing the redis server also did not work for  IP:%s and port:%d", R.IP, R.Port)
		}
		return false
	}

	return true

}

func (R *RedMon) Die() bool {
	//err := nil
	err := R.Cmd.Process.Kill()
	if err != nil {
		R.L.Printf("Unable to kill the process %v", err)
		return false
	}

	//either the shutdown or a kill will stop the monitor also
	return true
}

func (R *RedMon) CheckMsg() {
	//check message from scheduler
	//currently we do it to see if scheduler asks us to quit

	//ToDo better error handling needed
	err := R.P.LoadMsg()
	if !err {
		R.L.Printf("Failed While Loading msg for proc %v from node %v", R.P.ID, R.P.Nodename)
		return
	}

	switch R.P.Msg {
	case "SHUTDOWN":
		err = R.Stop()
		if err {

			R.L.Printf("failed to stop the REDIS server")
		}
		//in any case lets stop monitoring
		R.monChan <- 1
	case "MASTER":
		R.MakeMaster()
	case "SLAVEOF":
		//If this is the message then this particular redis proc will become slave of a different master
	}

}

//Should be called by the Monitors on Slave Procs, this gives the boolien anser if the sync is complegted or not
func (R *RedMon) IsSyncComplete() bool {

	//time.Sleep(1 * time.Second)

	if R.Client == nil {
		return false
	}

	respStr, err := R.Client.Info("replication").Result()
	if err != nil {
		R.L.Printf("getting the repication stats from server at IP:%s and port:%d", R.IP, R.Port)
		//dont return but try next time in another second/.1 second
	}

	respArr := strings.Split(respStr, "\n")
	for _, resp := range respArr {
		R.L.Printf("resp = %v", resp)
		r := strings.Split(resp, ":")
		switch r[0] {
		case "role":
			if !strings.Contains(r[1], "slave") {
				R.L.Printf("Trying to call is sync, but this server is not really a slave IP:%s, port:%d", R.IP, R.Port)
				return false
			}
			continue
		case "master_sync_in_progress":
			if !strings.Contains(r[1], "0") {
				R.L.Printf("Sync not complete yet in slave IP:%s, port:%d", R.IP, R.Port)
				return false
			} else {
				return true
			}
		case "master_sync_last_io_seconds_ago":
			//If the sync is completed then return true
			return true
		default:
			continue
		}

	}

	//if we did not find a master_sync_in_progress or slave in return at all, then some other problem, try again
	return false
}

func (R *RedMon) MakeMaster() bool {

	//send the slaveof no one command to server
	_, err := R.Client.SlaveOf("no", "one").Result()
	if err != nil {
		R.L.Printf("Error turning the slave to Master at IP:%s and port:%d", R.IP, R.Port)
		return false
	}

	return true
}
