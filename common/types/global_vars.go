package types

import (
	"../store"
)

var (
	Gdb   store.DB //Golabal variables related to db connection/instace
	MemDb *InMem   //In memory store
)

//Global db connection pointer, this will be initialized once abe be used everywhere

//global Constants releated to ETCD
const (
	ETC_BASE_DIR = "/MrRedis"
	ETC_INST_DIR = ETC_BASE_DIR + "/Instances"
	ETC_CONF_DIR = ETC_BASE_DIR + "/Config"
)
