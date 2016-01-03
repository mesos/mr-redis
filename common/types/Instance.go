package types

import (
	"../store"
	"strconv"
)

//A structure that will be able to store a tree of data

type Instance struct {
	Name    string      //Name of the instance
	Type    string      //Type of the instance "Single Instance = S; Master-Slave  = MS; Cluster = C
	Masters int         //Number of masters in this Instance
	Slaves  int         //Number of slaves in this instance
	Status  string      //Status of this instance "CREATING/ACTIVE/DELETED/DISABLED"
	PNames  string      //A list of comma seperated ids of Redis Procs
	Procs   []RedisProc //An array of redis procs to be filled later
}

// Creates a new instance variable
// Fills up the structure and updates the central store
// Returns an instnace pointer
// Returns nil if the instance already exists

func NewInstance(Name string, Type string, Masters int, Slaves int) *Instance {

	p := &Instance{Name: Name, Type: Type, Masters: Masters, Slaves: Slaves}
	return p
}

// Loads up the datastructure for the given Service Name to the struture
// If the Instance cannot be loaded the it returns an error

func (I *Instance) Load() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := ETC_INST_DIR + "/" + P.Name + "/"
	I.Type = Gdb.Get(node_name + "Type")
	I.Masters = strconv.Itoa(Gdb.Get(node_name + "Masters"))
	I.Slaves = strconv.Itoa(Gdb.Get(node_name + "Slaves"))
	I.Status = Gdb.Get(node_name + "Status")
	I.PNames = Gdb.Get(node_name + "PNames")

	return true
}

//Writes the entier content of an instance into store

func (I *Instance) Sync() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := ETC_INST_DIR + "/" + P.Name + "/"

	Gdb.Set(node_name+"Type", I.Type)
	Gdb.Set(node_name+"Masters", fmt.Sprintf("%d", I.Masters))
	Gdb.Set(node_name+"Slaves", fmt.Sprintf("%d", I.Slaves))
	Gdb.Set(node_name+"Status", I.Status)
	Gdb.Set(node_name+"PNames", I.PNames)

	return true
}

func (I *Instance) SyncType(string) bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := ETC_INST_DIR + "/" + P.Name + "/"
	Gdb.Set(node_name+"Type", I.Type)
	return false
}

func (I *Instance) SyncSlaves() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := ETC_INST_DIR + "/" + P.Name + "/"
	Gdb.Set(node_name+"Slaves", fmt.Sprintf("%d", I.Slaves))
	Gdb.Set(node_name+"PNames", I.PNames)
	return true
}

func (I *Instance) SyncMasters() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := ETC_INST_DIR + "/" + P.Name + "/"
	Gdb.Set(node_name+"Masters", fmt.Sprintf("%d", I.Slaves))
	Gdb.Set(node_name+"PNames", I.PNames)
	return true
}
