package Typ

//A structure that will be able to store a tree of data

type INSTANCE struct {
	Name    string   //Name of the instance
	Type    string   //Type of the instance "Single Instance = S; Master-Slave  = MS; Cluster = C
	Masters int      //Number of masters in this Instance
	Slaves  int      //Number of slaves in this instance
	Status  string   //Status of this instance "CREATING/ACTIVE/DELETED/DISABLED"
	procs   []string //An array of redis-procs in this Instance
}

func NewInstance(Name string, Type string, Masters int, Slaves int) *INSTANCE {

	return nil
}

func (P *INSTANCE) Load() bool {

	return false
}

func (P *INSTANCE) Sync() bool {

	return false
}

func (P *INSTANCE) SyncType(string) bool {

	return false
}
func (P *INSTANCE) SyncSlaves() bool {

	return false
}
