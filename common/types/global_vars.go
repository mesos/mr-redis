package types

import (
	"container/list"

	"../store"
)

var (
	Gdb   store.DB //Golabal variables related to db connection/instace
	MemDb *InMem   //In memory store

	OfferList *list.List       //list for having offer
	Cchan     chan *Instance   //Channel for Creator
	Mchan     chan *TaskUpdate //Channel for Maintainer
	Dchan     chan *Proc       //Channel for Destroyer
)

//Global db connection pointer, this will be initialized once abe be used everywhere

//global Constants releated to ETCD
const (
	ETC_BASE_DIR = "/MrRedis"
	ETC_INST_DIR = ETC_BASE_DIR + "/Instances"
	ETC_CONF_DIR = ETC_BASE_DIR + "/Config"
)

//Global constancts for Instnace Status
//CREATING/ACTIVE/DELETED/DISABLED
const (
	INST_STATUS_CREATING = "CREATING"
	INST_STATUS_RUNNING  = "RUNNING"
	INST_STATUS_DISABLED = "DISABLED"
)
