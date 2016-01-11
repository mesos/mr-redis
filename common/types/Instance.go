package types

import (
	"fmt"
	"log"
	"strconv"

	"../store/etcd"
)

//A structure that will be able to store a tree of data

type Instance struct {
	Name       string           //Name of the instance
	Type       string           //Type of the instance "Single Instance = S; Master-Slave  = MS; Cluster = C"
	Capacity   int              //Capacity of the Instance in MB
	Masters    int              //Number of masters in this Instance
	Slaves     int              //Number of slaves in this Instance
	ExpMasters int              //Expected number of Masters
	ExpSlaves  int              //Expected number of Slaves
	Status     string           //Status of this instance "CREATING/ACTIVE/DELETED/DISABLED"
	Mname      string           //Name / task id of the master redis proc
	Snames     []string         //Name of the slave
	Procs      map[string]*Proc //An array of redis procs to be filled later
}

// Creates a new instance variable
// Fills up the structure and updates the central store
// Returns an instnace pointer
// Returns nil if the instance already exists

func NewInstance(Name string, Type string, Masters int, Slaves int, Cap int) *Instance {

	p := &Instance{Name: Name, Type: Type, ExpMasters: Masters, ExpSlaves: Slaves, Capacity: Cap}
	return p
}

//Load an instnace from the store using Instance Name from the store
// if the instance is unavailable then return nil

func LoadInstance(Name string) *Instance {

	if Gdb.IsSetup() != true {
		return nil
	}

	node_name := etcd.ETC_INST_DIR + "/" + Name

	if ok, _ := Gdb.IsKey(node_name); !ok {
		return nil
	}

	I := &Instance{Name: Name}

	I.Load()

	return I

}

// Loads up the datastructure for the given Service Name to the struture
// If the Instance cannot be loaded the it returns an error

func (I *Instance) Load() bool {

	var err error
	var tmp_str string

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := etcd.ETC_INST_DIR + "/" + I.Name + "/"
	I.Type, err = Gdb.Get(node_name + "Type")
	tmp_str, err = Gdb.Get(node_name + "Capacity")
	I.Capacity, err = strconv.Atoi(tmp_str)
	tmp_str, err = Gdb.Get(node_name + "Masters")
	I.Masters, err = strconv.Atoi(tmp_str)
	tmp_str, err = Gdb.Get(node_name + "Slaves")
	I.Slaves, err = strconv.Atoi(tmp_str)
	tmp_str, err = Gdb.Get(node_name + "ExpMasters")
	I.ExpMasters, err = strconv.Atoi(tmp_str)
	tmp_str, err = Gdb.Get(node_name + "ExpSlaves")
	I.ExpSlaves, err = strconv.Atoi(tmp_str)
	I.Status, err = Gdb.Get(node_name + "Status")
	I.Mname, err = Gdb.Get(node_name + "Mname")

	node_name_slaves := node_name + "Snames/"
	I.Snames, err = Gdb.ListSection(node_name_slaves, false)

	if err != nil {
		log.Printf("The error value is %v", err)
	}

	return true
}

//Writes the entier content of an instance into store

func (I *Instance) Sync() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := etcd.ETC_INST_DIR + "/" + I.Name + "/"

	Gdb.Set(node_name+"Type", I.Type)
	Gdb.Set(node_name+"Masters", fmt.Sprintf("%d", I.Masters))
	Gdb.Set(node_name+"Capacity", fmt.Sprintf("%d", I.Masters))
	Gdb.Set(node_name+"ExpMasters", fmt.Sprintf("%d", I.Masters))
	Gdb.Set(node_name+"ExpSlaves", fmt.Sprintf("%d", I.Slaves))
	Gdb.Set(node_name+"Status", I.Status)
	Gdb.Set(node_name+"Mname", I.Mname)

	//Create Section for Slaves and Procs
	node_name_slaves := node_name + "Snames/"

	Gdb.CreateSection(node_name_slaves)
	for _, sname := range I.Snames {
		Gdb.Set(node_name_slaves+sname, "slave")
	}

	node_name_procs := node_name + "Procs/"
	Gdb.CreateSection(node_name_procs)

	//for _, p := range I.Procs {
	//p.Sync()
	//}

	return true
}

func (I *Instance) SyncType(string) bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := etcd.ETC_INST_DIR + "/" + I.Name + "/"
	Gdb.Set(node_name+"Type", I.Type)
	return false
}

func (I *Instance) SyncSlaves() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := etcd.ETC_INST_DIR + "/" + I.Name + "/"
	Gdb.Set(node_name+"Slaves", fmt.Sprintf("%d", I.Slaves))
	Gdb.Set(node_name+"ExpSlaves", fmt.Sprintf("%d", I.ExpSlaves))
	//Create Section for Slaves and Procs
	node_name_slaves := node_name + "Snames/"

	Gdb.CreateSection(node_name_slaves)
	for _, sname := range I.Snames {
		Gdb.Set(node_name_slaves+sname, "slave")
	}
	return true
}

func (I *Instance) SyncMasters() bool {

	if Gdb.IsSetup() != true {
		return false
	}

	node_name := etcd.ETC_INST_DIR + "/" + I.Name + "/"
	Gdb.Set(node_name+"Masters", fmt.Sprintf("%d", I.Slaves))
	Gdb.Set(node_name+"ExpMasters", fmt.Sprintf("%d", I.ExpSlaves))
	Gdb.Set(node_name+"Mname", I.Mname)
	return true
}
