package RedMon

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	typ "github.com/mesos/mr-redis/common/types"
	redisclient "gopkg.in/redis.v3"
)

//RedMon This structure is used to implement a monitor thread/goroutine for a running Proc(redisProc)
//This structure should be extended only if more functionality is required on the Monitoring functionality
//A Redis Proc's objec is created within this and monitored hence forth
type RedMon struct {
	P       *typ.Proc //The Proc structure that should be used
	Pid     int       //The Pid of the running proc
	IP      string    //IP address the redis instance should bind to
	Port    int       //The port number of this redis instance to be started
	Ofile   io.Writer //Stdout log file to be re-directed to this io.writer
	Efile   io.Writer //stderr of the redis instance should be re-directed to this file
	MS_Sync bool      //Make this as master after sync
	MonChan chan int
	Cmd     *exec.Cmd
	Client  *redisclient.Client //redis client library connection handler
	L       *log.Logger         //to redirect log outputs to a file
	//cgroup *CgroupManager		//Cgroup manager/cgroup connection pointer
}

//NewRedMon Create a new monitor based on the Data sent along with the TaskInfo
//The data could have the following details
//Capacity Master                 => Just start this PROC send update as TASK_RUNNING and monitor henceforth
//Capacity SlaveOf IP:Port        => This is a redis slave so start it as a slave, sync and then send TASK_RUNNING update then Monitor
//Capacity Master-SlaveOf IP:Port => This is a New master of the instance with an upgraded memory value so
//                          Start as slave, Sync data, make it as master, send TASK_RUNNING update and start to Monitor
func NewRedMon(tskName string, IP string, Port int, data string, L *log.Logger) *RedMon {

	var R RedMon
	var P *typ.Proc

	R.MonChan = make(chan int)
	R.Port = Port
	R.IP = IP

	//ToDo does this need error handling
	R.L = L

	R.L.Printf("Split data received is %v\n", data)

	splitData := strings.Split(data, " ")
	if len(splitData) < 1 || len(splitData) > 4 {
		//Print an error this is not suppose to happen
		R.L.Printf("RedMon Splitdata error %v\n", splitData)
		return nil
	}

	Cap, _ := strconv.Atoi(splitData[0])

	switch splitData[1] {
	case "Master":
		P = typ.NewProc(tskName, Cap, "M", "")
		R.L.Printf("created proc for new MASTER\n")
		break
	case "SlaveOf":
		P = typ.NewProc(tskName, Cap, "S", splitData[2])
		break
	case "Master-SlaveOf":
		P = typ.NewProc(tskName, Cap, "MS", splitData[2])
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
	}

	if !R.MS_Sync {
		return R.StartSlave()
	}
	//Posibly a scale request so start it as a slave, sync then make as master
	return R.StartSlaveAndMakeMaster()

}

//StartMaster Start the Proc as a master
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

//StartSlave start the proc as a slave, should be called to point to a valid master
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

//StartSlaveAndMakeMaster Start is as a slave and make it as a master, should be done for replication or adding a new slave
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

func fetchSubSection(value string, SubSection string) string {
	arr := strings.Split(value, "\r\n")

	for _, key := range arr {
		if strings.Contains(key, SubSection) {
			subArr := strings.Split(key, ":")
			if len(subArr) != 2 {
				return ""
			}
			return subArr[1]
		}
	}
	return ""
}

//GetRedisInfo Connect to the Redis Proc and collect info we need
func (R *RedMon) GetRedisInfo(Section string, Subsection string) string {

	value, err := R.Client.Info(Section).Result()
	if err != nil {
		R.L.Printf("STATS collection returned error on IP:%s and PORT:%d Err:%v for section %s subsection %s", R.IP, R.Port, err, Section, Subsection)
		return ""
	}
	return fetchSubSection(value, Subsection)
}

//UpdateStats Update the stats structure and flush it to the Store/DB
func (R *RedMon) UpdateStats() bool {

	var redisStats typ.Stats
	var txt string
	var err error

	txt = R.GetRedisInfo("Memory", "used_memory")
	redisStats.Mem, err = strconv.ParseInt(txt, 10, 64)
	if err != nil {
		R.L.Printf("UpdateStats(Mem) Unable to convert %s to number %v", txt, err)
	}

	txt = R.GetRedisInfo("Server", "uptime_in_seconds")
	redisStats.Uptime, err = strconv.ParseInt(txt, 10, 64)
	if err != nil {
		R.L.Printf("UpdateStats(Uptime) Uptime Unable to convert %s to number %v", txt, err)
	}

	txt = R.GetRedisInfo("Clients", "connected_clients")
	redisStats.Clients, err = strconv.Atoi(txt)
	if err != nil {
		R.L.Printf("UpdateStats(Clients) Unable to convert %s to number %v", txt, err)
	}

	txt = R.GetRedisInfo("Replication", "master_last_io_seconds_ago")
	redisStats.LastSyced, err = strconv.Atoi(txt)
	if err != nil && txt != "" {
		R.L.Printf("UpdateStats(master_last_io) Unable to convert %s to number %v", txt, err)
	}

	txt = R.GetRedisInfo("Replication", "slave_repl_offset")
	redisStats.SlaveOffset, err = strconv.ParseInt(txt, 10, 64)
	if err != nil && txt != "" {
		R.L.Printf("UpdateStats(slave_repl_offset) Unable to convert %s to number %v", txt, err)
	}

	txt = R.GetRedisInfo("Replication", "slave_priority")
	redisStats.SlavePriority, err = strconv.Atoi(txt)
	if err != nil && txt != "" {
		R.L.Printf("UpdateStats(slave_priority) Unable to convert %s to number %v", txt, err)
	}

	errSync := R.P.SyncStats(redisStats)
	if !errSync {
		R.L.Printf("Error syncing stats to store")
		return false
	}
	return true
}

//Monitor Primary monitor thread started for every PROC
func (R *RedMon) Monitor() bool {

	//wait for a second for the server to start
	//ToDo: is it needed

	CheckMsgCh := time.After(time.Second)
	UpdateStatsCh := time.After(2 * time.Second)

	for {
		if R.P.State == "Running" {
			select {

			case <-R.MonChan:
				//ToDo:update state if needed
				//signal to stop monitoring this
				R.L.Printf("Stopping RedMon for %s %s", R.P.IP, R.P.Port)
				return false

			case <-CheckMsgCh:
				//this is to check communication from scheduler; mesos messages are not reliable
				R.CheckMsg()
				CheckMsgCh = time.After(time.Second)

			case <-UpdateStatsCh:
				R.UpdateStats()
				UpdateStatsCh = time.After(2 * time.Second)
			}
		} else {
			time.Sleep(time.Second)
		}

	}

}

//Stop we have been told to stop the Redis
func (R *RedMon) Stop() bool {

	//send SHUTDOWN command for a gracefull exit of the redis-server
	//the server exited graceful will reflect at the task status FINISHED
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

//Die Kill the Redis Proc
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

//CheckMsg constantly keep checking if there is a new message for this Proc
func (R *RedMon) CheckMsg() {
	//check message from scheduler
	//currently we do it to see if scheduler asks us to quit

	//ToDo better error handling needed
	err := R.P.LoadMsg()
	if !err {
		R.L.Printf("Failed While Loading msg for proc %v from node %v", R.P.ID, R.P.Nodename)
		return
	}

	switch {
	case R.P.Msg == "SHUTDOWN":
		err = R.Stop()
		if err {

			R.L.Printf("failed to stop the REDIS server")
		}
		//in any case lets stop monitoring
		R.MonChan <- 1
		return
	case R.P.Msg == "MASTER":
		R.MakeMaster()
	case strings.Contains(R.P.Msg, "SLAVEOF"):
		R.TargetNewMaster(R.P.Msg)
		//If this is the message then this particular redis proc will become slave of a different master
	}
	//Once you have read the message delete the message.
	R.P.Msg = ""
	R.P.SyncMsg()

}

//IsSyncComplete Should be called by the Monitors on Slave Procs, this gives the boolien anser if the sync is complegted or not
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
			}
			return true
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

//MakeMaster Make a Proc as a master (ie: supply the command "slaveof no on" to the Proc
func (R *RedMon) MakeMaster() bool {

	//send the slaveof no one command to server
	_, err := R.Client.SlaveOf("no", "one").Result()
	if err != nil {
		R.L.Printf("Error turning the slave to Master at IP:%s and port:%d", R.IP, R.Port)
		return false
	}

	R.L.Printf("Slave of NO ONE worked")
	return true
}

//TargetNewMaster Make this Proc now target a new master, should be done when a new slave is promoted
func (R *RedMon) TargetNewMaster(Msg string) bool {

	SlaveofArry := strings.Split(Msg, " ") //Split it with space as while we are sending fromt the sheduler we send it of the format SLAVEOF<SPACE>IP<SPACE>PORT
	if len(SlaveofArry) != 3 {             //This should have three elements otherwise its an error

		R.L.Printf("Writing SLAVE of COMMAND %s", Msg)
		return false

	}

	//send the slaveof IP (Arry[1]) and PORT (Array[1])
	_, err := R.Client.SlaveOf(SlaveofArry[1], SlaveofArry[2]).Result()
	if err != nil {
		R.L.Printf("Error turning the slave to Master at IP:%s and port:%d", R.IP, R.Port)
		return false
	}

	R.L.Printf("Slave of %s %s worked", SlaveofArry[1], SlaveofArry[2])
	return true
}
