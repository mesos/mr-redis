package types

//A structure that will be able to store a tree of data

type ServiceInstance struct {
	Name    string   //Name of the instance
	Type    string   //Type of the instance "Single Instance = S; Master-Slave  = MS; Cluster = C
	Masters int      //Number of masters in this Instance
	Slaves  int      //Number of slaves in this instance
	Status  string   //Status of this instance "CREATING/ACTIVE/DELETED/DISABLED"
	procs   []string //An array of redis-procs in this Instance
}

func NewServiceInstance(Name string, Type string, Masters int, Slaves int) *ServiceInstance {

	return nil
}

func (P *ServiceInstance) Load() bool {

	return false
}

func (P *ServiceInstance) Sync() bool {

	return false
}

func (P *ServiceInstance) SyncType(string) bool {

	return false
}
func (P *ServiceInstance) SyncSlaves() bool {

	return false
}
