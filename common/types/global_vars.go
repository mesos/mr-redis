package types

import (
	"../store"
)

//Golabal variables related to db connection/instace
var Gdb store.DB

//Global db connection pointer, this will be initialized once abe be used everywhere

//global Constants releated to ETCD
const (
	ETC_BASE_DIR = "/MrRedis"
	ETC_INST_DIR = ETC_BASE_DIR + "/Instances"
	ETC_CONF_DIR = ETC_BASE_DIR + "/Config"
)
