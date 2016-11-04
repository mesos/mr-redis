package types

//Offer Structure that is used between creator and Mesos Scheduler
type Offer struct {
	Name         string //Name of the instance
	Taskname     string //Name of the redis proc
	Cpu          int    //CPU default is one
	Mem          int    //Memory in MB
	DValue       int    //Distribution Value
	IsMaster     bool   //Is this instance a master
	MasterIpPort string //If this is slave then send the masters IP and prot number
}

//NewOffer Returns a new offer which will be interpreted by the scheduler
func NewOffer(name string, tname string, cpu int, mem int, ismaster bool, masterIPPort string, dvalue int) Offer {
	return Offer{Name: name, Taskname: tname, Cpu: cpu, Mem: mem, IsMaster: ismaster, MasterIpPort: masterIPPort, DValue: dvalue}
}
