package Typ

type PROC struct {
	Name     string //Automatically generated UID
	Instance string //Name of the instance this proc belongs to
	Stats    string //Json string of stats
	CMD      string //Command passing between Maintainer and Executors
	Status   string //Status of the current proc
	Type     string //M=Master; S=Slave;  a single char entry that tells us what type of redis proc this is
	Slaveof  string //If this redis proc is a slave instance
}
