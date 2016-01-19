package mesoslib

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"

	typ "../../common/types"
)

func serveExecutorArtifact(path string, IP, Port string) (*string, string) {
	serveFile := func(pattern string, filename string) {
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filename)
		})
	}

	// Create base path (http://foobar:5000/<base>)
	pathSplit := strings.Split(path, "/")
	var base string
	if len(pathSplit) > 0 {
		base = pathSplit[len(pathSplit)-1]
	} else {
		base = path
	}
	serveFile("/"+base, path)

	hostURI := fmt.Sprintf("http://%s:%s/%s", IP, Port, base)
	log.Printf("Hosting artifact '%s' at '%s'", path, hostURI)

	return &hostURI, base
}

func prepareExecutorInfo(IP, Port, executorPath, DbType, DbEndPoint string) *mesos.ExecutorInfo {
	executorUris := []*mesos.CommandInfo_URI{}
	uri, executorCmd := serveExecutorArtifact(executorPath, IP, Port)
	executorUris = append(executorUris, &mesos.CommandInfo_URI{Value: uri, Executable: proto.Bool(true)})

	executorCommand := fmt.Sprintf("./%s -logtostderr=true -DbType=%s -DbEndPoint=%s", executorCmd, DbType, DbEndPoint)

	go http.ListenAndServe(fmt.Sprintf("%s:%s", IP, Port), nil)
	log.Printf("Serving executor artifacts...")

	// Create mesos scheduler driver.
	return &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID("default"),
		Name:       proto.String("MrRedis Executor"),
		Source:     proto.String("MrRedis"),
		Command: &mesos.CommandInfo{
			Value: proto.String(executorCommand),
			Uris:  executorUris,
		},
	}
}

// Mesos library will recive a string comman seperated with values that it needs to run with
// this function should parse those comma seperated values and supply it to mesos-library
// format config = "MasterIP","currentServerIP","MasterPort","currentServerPort"
// MasterIP/Port = Mesos Master ip or port
// Curre3ntServerIP = the ip address of the server at which framework/scheduler will run
// CurrentServerPort = The port at which we will distribute the executor to slaves
// Master Port and Current server port has default falues

func parseConfig(config string) (string, string, string, string) {

	splitconfig := strings.Split(config, ",")

	mIP := ""
	mP := "5050"
	sIP := ""
	sP := "5544"

	for i := 0; len(splitconfig) > 0; i++ {
		switch i {
		case 0:
			mIP = splitconfig[i] //Extract the master IP
			break
		case 1:
			sIP = splitconfig[i] //Extract the current server ip
			break
		case 2:
			mP = splitconfig[i] //Extract the master Port
			break
		case 3:
			sP = splitconfig[i] //Extract the current servers port at whichwe will distribute the executor
			break
		}
	}

	return mIP, mP, sIP, sP

}

func parseIP(address string) net.IP {
	addr, err := net.LookupIP(address)
	if err != nil {
		log.Fatal(err)
	}
	if len(addr) < 1 {
		log.Fatalf("failed to parse IP from address '%v'", address)
	}
	return addr[0]
}

const FailoverTime = 60 //Frameowkr and its task will be terminated if the framework is not started in 60 secons
const TimeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

//If the frameowkr was regiestered before the Failover tiemout value then regiester as a new framework
func GetFrameWorkID() (string, float64) {

	fTimout := float64(FailoverTime)
	fwTStamp, terr := typ.Gdb.Get(typ.ETC_CONF_DIR + "/RegesteredAt")
	t, tperr := time.Parse(TimeFormat, fwTStamp)
	fwID, err := typ.Gdb.Get(typ.ETC_CONF_DIR + "/FrameworkID")

	if err != nil || terr != nil || tperr != nil {
		log.Printf("Not registered previously")
		return "", fTimout
	}

	delta_t := time.Now().Sub(t)

	if (delta_t / time.Second) < FailoverTime {
		return fwID, fTimout
	}

	return "", fTimout

}

func Run(MasterIP, MasterPort, ServerIP, ServerPort, executorPath, DbType, DbEndPoint string) {

	//Split the configuration string

	//MasterIP, MasterPort, ServerIP, ServerPort = parseConfig(config)

	//Get executor information
	exec := prepareExecutorInfo(ServerIP, ServerPort, executorPath, DbType, DbEndPoint)

	fwID, fTimout := GetFrameWorkID()

	// the framework
	fwinfo := &mesos.FrameworkInfo{
		User:            proto.String(""), // Mesos-go will fill in user.
		Name:            proto.String("MrRedis"),
		Id:              &mesos.FrameworkID{Value: &fwID},
		FailoverTimeout: &fTimout,
	}

	//Add mesos authentication code
	//TODO

	//create the scheduler dirver object
	sched_config := sched.DriverConfig{
		Scheduler:      NewMrRedisScheduler(exec),
		Framework:      fwinfo,
		Master:         MasterIP + ":" + MasterPort,
		Credential:     nil,
		BindingAddress: parseIP(ServerIP),
	}

	driver, err := sched.NewMesosSchedulerDriver(sched_config)

	if err != nil {
		log.Fatalf("Framework is not created error %v", err)
	}

	log.Printf("The Framework ID is %v and %v", fwinfo.Id, sched_config.Framework.Id)

	status, err := driver.Run()

	if err != nil {
		log.Fatalf("Framework Run() error %v", err)
	}

	log.Printf("Framework Terminated with status %v", status.String())

}
