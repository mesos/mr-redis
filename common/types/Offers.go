package types

type Offer struct {
	Taskname string //Name of the redis proc
	Cpu      int    //CPU default is one
	Mem      int    //Memory in MB
	IsMaster bool   //Is this instance a master
}

func NewOffer(name string, cpu int, mem int, ismaster bool) Offer {
	return Offer{Taskname: name, Cpu: cpu, Mem: mem, IsMaster: ismaster}
}
