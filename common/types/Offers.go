package types

//Offer Structure that is used between creator and Mesos Scheduler
type Offer struct {
	Taskname     string //Name of the redis proc
	Cpu          int    //CPU default is one
	Mem          int    //Memory in MB
	IsMaster     bool   //Is this instance a master
	MasterIpPort string //If this is slave then send the masters IP and prot number
}

//NewOffer Returns a new offer which will be interpreted by the scheduler
func NewOffer(name string, cpu int, mem int, ismaster bool, masterIPPort string) Offer {
	return Offer{Taskname: name, Cpu: cpu, Mem: mem, IsMaster: ismaster, MasterIpPort: masterIPPort}
}
