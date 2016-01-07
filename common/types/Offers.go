package types

type Offer struct {
	taskname string //Name of the redis proc
	cpu      int    //CPU default is one
	mem      int    //Memory in MB
}

func NewOffer(name string, cpu int, mem int) *Offer {
	return &Offer{taskname: name, cpu: cpu, mem: mem}
}
