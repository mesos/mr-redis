package types

import (
	"container/list"
	"log"

	"../store/etcd"
)

func Initialize(dbtype string, config string) (bool, error) {

	//Initalize all the communication channels
	OfferList = list.New()
	OfferList.Init()
	Cchan = make(chan *Instance)
	Mchan = make(chan *TaskUpdate) //Channel for Maintainer
	Dchan = make(chan *Proc)       //Channel for Destroyer

	//Initalize the Internal in-memory storage
	MemDb = NewInMem()

	//Initalize the store db
	switch dbtype {
	case "etcd":
		Gdb = etcd.New()
		err := Gdb.Setup(config)
		if err != nil {
			log.Fatalf("Failed to setup etcd database error:%v", err)
		}
		return Gdb.IsSetup(), nil
		break
	}

	return true, nil
}
