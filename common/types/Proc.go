package types

import (
	"log"
	"strconv"

	//"../redis"
	"../store"
	"../store/etcd"
)

//A standalone redis KV store is usually started in any slave (Linux) like below
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
	Stats    string //All other statistics apart from Memory usage to be stored as a json/string
	Msg      string //Message we will revive fromt he scheduler and action to be taken on it
	IP       string //IP address of the slave at which this redis-server proc is running
	Port     string //Port number at which this PROC is bound to
	EID      string //Executor ID of this PROC  .. Just in case we need to send a framework messsage
	SID      string //Slave ID of this PROC .. Just in case we need to send a framework message
	//cli    redis.Cli
}

func NewProc(TskName string, Capacity int, Type string, SlaveOf string) *Proc {

	var tmpProc Proc
	Tids := strings.Split(TskName, "::")

	if len(Tid) != 2 {
		//Something wrong the TaskID should be of the format <InstanceName>::<UUID of the PROC>
		//Throw an error and ignore
		log.Printf("Wrong format Task Name %s", TskName)
		return nil
	}

	tmpProc.Instance = Tids[0]
	tmpProc.ID = Tids[1]
	tmpProc.MemCap = Capacity
	tmpProc.Type = Type
	tmpPorc.SlaveOf = SlaveOf

	return &tmpProc
}

//Load a Proc information from the store to structure and return
func Load(TskName string) *Proc {

	var P Proc

	Tids := strings.Split(TskName, "::")

	if len(Tid) != 2 {
		log.Printf("Proc.Load() Wrong format Task Name %s", TskName)
		return false
	}

	P.Instance = Tids[0]
	P.ID = Tids[1]

	P.Nodename = etcd.ETC_INST_DIR + "/" + P.Instance + "/PROC/" + I.ID

	P.Load()

	return &P
}

//Load the latest from ETC store
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

	tmpStr, err = Gdb.Get(P.Nodename + "/Port")
	P.Port, err = strconv.Atoi(tmpStr)

	tmpStr, err = Gdb.Get(P.Nodename + "/Pid")
	P.Pid, err = strconv.Atoi(tmpStr)

	P.State, err = Gdb.Get(P.Nodename + "/State")
	P.Type, err = Gdb.Get(P.Nodename + "/Type")
	P.EID, err = Gdb.Get(P.Nodename + "/EID")
	P.SID, err = Gdb.Get(P.Nodename + "/SID")
	P.Msg, err = Gdb.Get(P.Nodename + "/Msg")
	P.Stats, err = Gdb.Get(P.Nodename + "/Stats")

}

//Sync everything thats in-memory to the the central store
func (P *Proc) Sync() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	//Attempt to create the directory/section for storing the PROC relevent information in the instance
	Gdb.CreateSection(P.Nodename)

	Gdb.Set(P.Nodename+"/Instance", P.Instance)
	Gdb.Set(P.Nodename+"/Nodename", P.Nodename)
	Gdb.Set(P.MemCap+"/Capacity", fmt.Sprintf("%s", P.MemCap))
	Gdb.Set(P.MemCap+"/MemUsed", fmt.Sprintf("%s", P.MemUsed))
	Gdb.Set(P.MemCap+"/Pid", fmt.Sprintf("%s", P.Pid))
	Gdb.Set(P.MemCap+"/State", P.State)
	Gdb.Set(P.MemCap+"/Stats", P.Stats)
	Gdb.Set(P.MemCap+"/Msg", P.Msg)
	Gdb.Set(P.MemCap+"/EID", P.EID)
	Gdb.Set(P.MemCap+"/SID", P.SID)
	Gdb.Set(P.MemCap+"/Type", P.Type)

	return true
}

func (P *Proc) SyncStats() bool {
	if Gdb.IsSetup() != true {
		return false
	}

	Gdb.Set(P.Nodename+"/Stats", P.Stats)

}

func (P *Proc) SyncType() bool {
	if Gdb.IsSetup() != true {
		return false
	}
	Gdb.Set(P.NodeName+"/Type", P.Type)
}

func (P *Proc) LoadStats() bool {
	var err error
	if Gdb.IsSetup() != true {
		return false
	}
	P.Stats, err = Gdb.Get(P.Nodename + "/Stats")
}

func (P *Proc) LoadType() bool {
	var err error
	if Gdb.IsSetup() != true {
		return false
	}
	P.Type, err = Gdb.Get(P.Nodename + "/Type")
}
