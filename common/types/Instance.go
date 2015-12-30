package types

import (
	"fmt"
	"log"
)

type INSTANCE struct {
	Name    string   //Name of the instance
	Type    string   //Type of the instance "Single Instance = S; Master-Slave  = MS; Cluster = C
	Masters int      //Number of masters in this Instance
	Slaves  int      //Number of slaves in this instance
	Status  string   //Status of this instance "CREATING/ACTIVE/DELETED/DISABLED"
	procs   []string //An array of redis-procs in this Instance
}

type PROC struct {
	Name     string //Automatically generated UID
	Instance string //Name of the instance this proc belongs to
	Stats    string //Json string of stats
	CMD      string //Command passing between Maintainer and Executors
	Status   string //Status of the current proc
	Type     string //M=Master; S=Slave;  a single char entry that tells us what type of redis proc this is
	Slaveof  string //If this redis proc is a slave instance
}

func NewInstance(Name string, Type string, Masters int, Slaves int) *INSTANCE {
}

func (P *INSTANCE) Load() bool {
}

func (P *INSTANCE) Sync() bool {
}

func (P *INSTANCE) Type() string {
}

func (P *INSTANCE) SyncType(string) bool {
}
func (P *INSTANCE) SyncSlaves() bool {
}
func (P *INSTANCE) SyncType() bool {
}
