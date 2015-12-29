package serviceproc

import (
	"fmt"
	"log"
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

func NewRedisProc(procofid string, port int) *RedisProc {
	//tbd: should generate and fill in the procid
	//tbd: should find out a mechanism to get the instance id in procofid
	//tbd: get the local system IP and fill in the same; what if there are multiple ips?
	return &RedisProc{Mem: 0, Cpu: 0, Portno: port, IP: "", ID: "", ProcofID: procofid}
}

func (rp *RedisProc) Start(port int) error {

	rp.cmd = exec.Command("redis-server", "--port", fmt.Sprintf("%d", port))
	err := rp.cmd.Start()
	if err != nil {
		fmt.Println("error starting the redis server\n")
		log.Fatal(err)
		return err
	}

	fmt.Println("Waiting for the redis server to finish\n")

	err = rp.cmd.Wait()
	if err != nil {
		fmt.Println("error waiting for redis server to finish\n")
		log.Fatal(err)
		return err
	}

	return nil
}

func (rp *RedisProc) Stop() error {

	//err := nil
	err := rp.cmd.Process.Kill()
	if err != nil {
		log.Printf("Unable to kill the process %v", err)
	}
	return err

}
