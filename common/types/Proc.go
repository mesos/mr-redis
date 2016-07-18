package types

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	//"../redis"
	//"../store"
	"github.com/mesos/mr-redis/common/store/etcd"
)

//Proc A standalone redis KV store is usually started in any slave (Linux) like below
//$./redis-server -p <PORT> ..... {OPTIONS}
//This stand alone redis-server will be an actual unix process bound to a particular port witha PID
//A redis Master Slave setup will have two such "redis-server" processes running in either the same machine or two different machines
//A redis KV with one master and 3 slaves will have a total of 4 "redis-server" processes running
//The below structure "Proc" is a representation of such a running 'redis-server' process started via this framework
//This started/running "Proc" could be a Master/Standalone instance or could be a Slave of another "redis-server' running as a master
type Proc struct {
	Instance string //Name of the instance it belongs to
	Nodename string //Node name at which this should start syncing its details to
	MemCap   int    //Maximum Memory this instance can go to
	MemUsed  int    //Current usage of the memory
	Pid      int    //Unix Process id of this running instance
	ID       string //UUID that was generated for this PROC
	State    string //Current state of the process Active/Dead/Crashed etc.,
	Type     string //Type of the PROC master/Slave etc.,
	SlaveOf  string //Slave of which redis master
	Stats    string //All other statistics apart from Memory usage to be stored as a json/string
	Msg      string //Message we will revive fromt he scheduler and action to be taken on it
	IP       string //IP address of the slave at which this redis-server proc is running
	Port     string //Port number at which this PROC is bound to
	EID      string //Executor ID of this PROC  .. Just in case we need to send a framework messsage
	SID      string //Slave ID of this PROC .. Just in case we need to send a framework message
	//cli    redis.Cli
}

//Stats the whole stats structure could be a json structure reflecting all the fields what redis info returns
//currently one field has many new line saperated values;ToDO will this work if returned in API?
type Stats struct {
	Uptime        int64
	Mem           int64
	Clients       int
	LastSyced     int
	SlaveOffset   int64 //Offset of the slave
	SlavePriority int
}

//ProcJson Fields to be packed in a json when replied to a HTTP REST query.
type ProcJson struct {
	IP                 string
	Port               string
	MemoryCapacity     int
	MemoryUsed         int64
	Uptime             int64
	ClientsConnected   int
	LastSyncedToMaster int
}

//NewProc Constructor for a PROC struct, this does not sync anything to the DB
func NewProc(TskName string, Capacity int, Type string, SlaveOf string) *Proc {

	var tmpProc Proc
	Tids := strings.Split(TskName, "::")

	if len(Tids) != 2 {
		//Something wrong the TaskID should be of the format <InstanceName>::<UUID of the PROC>
		//Throw an error and ignore
		log.Printf("Wrong format Task Name %s", TskName)
		return nil
	}

	tmpProc.Instance = Tids[0]
	tmpProc.ID = Tids[1]
	tmpProc.MemCap = Capacity
	tmpProc.Type = Type
	tmpProc.SlaveOf = SlaveOf

	tmpProc.Nodename = etcd.ETC_INST_DIR + "/" + tmpProc.Instance + "/Procs/" + tmpProc.ID
	return &tmpProc
}

//LoadProc Load a Proc information from the store and populate a structure
func LoadProc(TskName string) *Proc {

	var P Proc

	Tids := strings.Split(TskName, "::")

	if len(Tids) != 2 {
		log.Printf("Proc.Load() Wrong format Task Name %s", TskName)
		return nil
	}

	P.Instance = Tids[0]
	P.ID = Tids[1]

	P.Nodename = etcd.ETC_INST_DIR + "/" + P.Instance + "/Procs/" + P.ID

	P.Load()

	return &P
}

//Load the latest from KV store
func (P *Proc) Load() bool {

	var err error
	var tmpStr string
	if Gdb.IsSetup() != true {
		return false
	}

	if ok, _ := Gdb.IsKey(P.Nodename); !ok {
		log.Printf("Invalid Key %v, Cannot load", P.Nodename)
		return false
	}

	tmpStr, err = Gdb.Get(P.Nodename + "/Capacity")
	P.MemCap, err = strconv.Atoi(tmpStr)

	tmpStr, err = Gdb.Get(P.Nodename + "/MemUsed")
	P.MemUsed, err = strconv.Atoi(tmpStr)

	P.Port, err = Gdb.Get(P.Nodename + "/Port")
	P.IP, err = Gdb.Get(P.Nodename + "/IP")

	tmpStr, err = Gdb.Get(P.Nodename + "/Pid")
	P.Pid, err = strconv.Atoi(tmpStr)

	P.State, err = Gdb.Get(P.Nodename + "/State")
	P.Type, err = Gdb.Get(P.Nodename + "/Type")
	P.EID, err = Gdb.Get(P.Nodename + "/EID")
	P.SID, err = Gdb.Get(P.Nodename + "/SID")
	P.Msg, err = Gdb.Get(P.Nodename + "/Msg")
	P.Stats, err = Gdb.Get(P.Nodename + "/Stats")
	P.SlaveOf, err = Gdb.Get(P.Nodename + "/SlaveOf")

	if err != nil {
		log.Printf("Error occured %v", err)
		return false
	}

	return true

}

//Sync everything thats in-memory to the KV store, should be used only whne you have more attributes to be synced/written to the DB at the same time, Disk intensive operation
func (P *Proc) Sync() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	//Attempt to create the directory/section for storing the PROC relevant information in the instance
	//P.Nodename = etcd.ETC_INST_DIR + "/" + P.Instance + "/PROC/" + P.ID
	Gdb.CreateSection(P.Nodename)

	Gdb.Set(P.Nodename+"/Instance", P.Instance)
	Gdb.Set(P.Nodename+"/Nodename", P.Nodename)
	Gdb.Set(P.Nodename+"/Capacity", fmt.Sprintf("%d", P.MemCap))
	Gdb.Set(P.Nodename+"/MemUsed", fmt.Sprintf("%d", P.MemUsed))
	Gdb.Set(P.Nodename+"/Pid", fmt.Sprintf("%d", P.Pid))
	Gdb.Set(P.Nodename+"/IP", P.IP)
	Gdb.Set(P.Nodename+"/Port", P.Port)
	Gdb.Set(P.Nodename+"/State", P.State)
	Gdb.Set(P.Nodename+"/Stats", P.Stats)
	Gdb.Set(P.Nodename+"/SlaveOf", P.SlaveOf)
	Gdb.Set(P.Nodename+"/Msg", P.Msg)
	Gdb.Set(P.Nodename+"/EID", P.EID)
	Gdb.Set(P.Nodename+"/SID", P.SID)
	Gdb.Set(P.Nodename+"/Type", P.Type)

	return true
}

//SyncStats Updates only the statistic related information to the disk, used by RedMon every second
func (P *Proc) SyncStats(s Stats) bool {
	if Gdb.IsSetup() != true {
		return false
	}

	sBytes, err := json.Marshal(s)

	if err != nil {
		log.Printf("SyncStats() unbale to marshal error %v", err)
		return false
	}

	P.Stats = string(sBytes)

	Gdb.Set(P.Nodename+"/Stats", P.Stats)

	return true
}

//SyncType should be used to only update the TYPE attribute of a PROC
func (P *Proc) SyncType() bool {
	if Gdb.IsSetup() != true {
		return false
	}
	Gdb.Set(P.Nodename+"/Type", P.Type)

	return true
}

//SyncMsg Should be used only to update a MSG with respect to a PROC
func (P *Proc) SyncMsg() bool {
	if Gdb.IsSetup() != true {
		return false
	}
	Gdb.Set(P.Nodename+"/Msg", P.Msg)

	return true
}

//SyncSlaveOf should be used to make a slave (proc) point to a new master
func (P *Proc) SyncSlaveOf() bool {
	if Gdb.IsSetup() != true {
		return false
	}
	Gdb.Set(P.Nodename+"/SlaveOf", P.SlaveOf)

	return true
}

//LoadStats get the latest STATS info from the store, usually called by the scheduler to make a decision
func (P *Proc) LoadStats() *Stats {
	var err error
	var s Stats
	if Gdb.IsSetup() != true {
		return nil
	}
	P.Stats, err = Gdb.Get(P.Nodename + "/Stats")

	if err != nil {
		log.Printf("Error occured %v", err)
		return nil
	}

	err = json.Unmarshal([]byte(P.Stats), &s)

	if err != nil {
		log.Printf("Error occured un-marshalling stats LoadStats() %v stats=%s", err, P.Stats)
		return nil
	}
	return &s
}

//LoadType Update what is the current type of a PROC, usually called by the scheduler
func (P *Proc) LoadType() bool {
	var err error
	if Gdb.IsSetup() != true {
		return false
	}
	P.Type, err = Gdb.Get(P.Nodename + "/Type")
	if err != nil {
		log.Printf("Error occured %v", err)
		return false
	}
	return true
}

//LoadMsg Get the latest MSG from the scheduler, usually called by the executor(RedMon)
func (P *Proc) LoadMsg() bool {
	var err error
	if Gdb.IsSetup() != true {
		return false
	}

	P.Msg, err = Gdb.Get(P.Nodename + "/Msg")
	if err != nil {
		log.Printf("Error occured %v", err)
		return false
	}

	return true
}

//ToJson Convert a PROC to json but only selected fileds
func (P *Proc) ToJson() *ProcJson {

	var pJson ProcJson
	pJson.IP = P.IP
	pJson.Port = P.Port
	pJson.MemoryCapacity = P.MemCap
	stats := P.LoadStats()
	if stats == nil {
		return nil
	}
	pJson.MemoryUsed = stats.Mem
	pJson.ClientsConnected = stats.Clients
	pJson.Uptime = stats.Uptime
	pJson.LastSyncedToMaster = stats.LastSyced

	return &pJson

	/*
		ret_bytes, err := json.Marshal(p_json)

		if err != nil {
			return "{LoadStats Failed PROC}"
		}

		return string(ret_bytes)
	*/
}

//ToJsonStats The stats are always store in JSON format in the DB/Store against a single key
func (P *Proc) ToJsonStats(stats Stats) string {

	retBytes, err := json.Marshal(stats)

	if err != nil {
		return "{Invalid Stats}"
	}

	return string(retBytes)
}
