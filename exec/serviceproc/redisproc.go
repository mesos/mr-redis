package serviceproc

import (

	"../../common/store"

	"fmt"
	"log"
	"os/exec"
	"strconv"
)

//Golabal variables related to db connection/instace
var Store store.DB //Global db connection pointer, this will be initialized once abe be used everywhere

//global Constants releated to ETCD
const (
	ETC_BASE_DIR = "/MrRedis"
	ETC_INST_DIR = ETC_BASE_DIR + "/Instances"
	ETC_CONF_DIR = ETC_BASE_DIR + "/Config"
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

func NewRedisProc(procofid string, port int, procid string) *RedisProc {

	//tbd: should find out a mechanism to get the instance id in procofid
	//tbd: get the local system IP and fill in the same; what if there are multiple ips?

	return &RedisProc{Mem: 0, Cpu: 0, Portno: port, IP: "", ID: procid, ProcofID: procofid}
}

func (rp *RedisProc) SettoStore() error {
	if Store.IsSetup() != true {
		return fmt.Errorf("the etcd store isn't setup")
	}

	node_name := ETC_INST_DIR + "/" + rp.ProcofID + "/" + rp.ID + "/"

	Store.Set(node_name+"Portno", fmt.Sprintf("%d", rp.Portno))

	return nil
}

func (rp *RedisProc) GetfromStore() error {

	if Store.IsSetup() != true {
		return fmt.Errorf("the etcd store isn't setup")
	}

	node_name := ETC_INST_DIR + "/" + rp.ProcofID + "/" + rp.ID + "/"
	tmp_str, err := Store.Get(node_name + "Portno")
	rp.Portno, err = strconv.Atoi(tmp_str)

	fmt.Println("STORE TEST: The retrieved value of port is %s", tmp_str)

	return err
}

func (rp *RedisProc) Start(port int) error {

	rp.cmd = exec.Command("redis-server", "--port", fmt.Sprintf("%d", port))
	err := rp.cmd.Start()
	if err != nil {
		fmt.Println("error starting the redis server\n")
		log.Println(err)
		return err
	}

	fmt.Println("Waiting for the redis server to finish\n")

	//store the data of this proc to etcd
	err = rp.SettoStore()
	if err != nil {
		fmt.Println("error writing to store\n")
		log.Println(err)
		return err
	}

	err = rp.cmd.Wait()
	if err != nil {
		fmt.Println("error waiting for redis server to finish\n")
		log.Println(err)
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
