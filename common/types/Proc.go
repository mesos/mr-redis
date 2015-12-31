package types

import (
	"os/exec"
)

type RedisProc struct {
	cmd      *exec.Cmd
	Mem      int
	Cpu      int
	Portno   int
	IP       string //this machines ip
	ID       string //to be filled as unique id
	ProcofID string //the service insts id which this proc is part of
	State    string
}
