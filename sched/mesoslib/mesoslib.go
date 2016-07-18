package mesoslib

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"

	typ "github.com/mesos/mr-redis/common/types"
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

func prepareExecutorInfo(IP, Port, executorPath, redisPath, DbType, DbEndPoint string) *mesos.ExecutorInfo {
	executorUris := []*mesos.CommandInfo_URI{}
	uri, executorCmd := serveExecutorArtifact(executorPath, IP, Port)
	executorUris = append(executorUris, &mesos.CommandInfo_URI{Value: uri, Executable: proto.Bool(true)})
	uri, _ = serveExecutorArtifact(redisPath, IP, Port)
	executorUris = append(executorUris, &mesos.CommandInfo_URI{Value: uri, Executable: proto.Bool(true)})

	executorCommand := fmt.Sprintf("./%s -logtostderr=true -DbType=%s -DbEndPoint=%s", executorCmd, DbType, DbEndPoint)

	/* If possible override the artifact hosting IP to below env variable */

	go func(IP, Port string) {

		hostIP := os.Getenv("HOST")

		if hostIP == "" {
			hostIP = IP
		}

		log.Printf("hostIP = %s going to listen and serve", hostIP)

		err := http.ListenAndServe(fmt.Sprintf("%s:%s", hostIP, Port), nil)
		log.Printf("Serving executor artifacts... error = %v", err)
	}(IP, Port)

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

// Mesos library will recive a string comman separated with values that it needs to run with
// this function should parse those comma separated values and supply it to mesos-library
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
	hostIP := os.Getenv("HOST")

	if hostIP == "" {
		hostIP = address
	}
	addr, err := net.LookupIP(hostIP)
	if err != nil {
		log.Fatal(err)
	}
	if len(addr) < 1 {
		log.Fatalf("failed to parse IP from address '%v'", address)
	}
	return addr[0]
}

//FailoverTime Frameowkr and its task will be terminated if the framework is not started in 60 secons
const FailoverTime = 60

//TimeFormat we need to parse the Timestamp
const TimeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

//GetFrameWorkID If the framework was regiestered before the Failover timeout value then regiester as a new framework
func GetFrameWorkID() (string, float64) {

	fTimout := float64(FailoverTime)
	fwTStamp, terr := typ.Gdb.Get(typ.ETC_CONF_DIR + "/RegisteredAt")
	t, tperr := time.Parse(TimeFormat, fwTStamp)
	fwID, err := typ.Gdb.Get(typ.ETC_CONF_DIR + "/FrameworkID")

	if err != nil || terr != nil || tperr != nil {
		log.Printf("Not registered previously")
		return "", fTimout
	}

	deltaT := time.Now().Sub(t)
	log.Printf("Delta of the previously registered framework is = %v", deltaT)

	if (deltaT / time.Second) < FailoverTime {
		return fwID, fTimout
	}

	return "", fTimout

}

//Run primary function that starts the Mesos Scheduler
func Run(MasterEndPoint, ServerIP, ServerPort, executorPath, redisPath, DbType, DbEndPoint string) {

	//Split the configuration string

	//MasterIP, MasterPort, ServerIP, ServerPort = parseConfig(config)

	//Get executor information
	exec := prepareExecutorInfo(ServerIP, ServerPort, executorPath, redisPath, DbType, DbEndPoint)

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
	schedConfig := sched.DriverConfig{
		Scheduler:      NewMrRedisScheduler(exec),
		Framework:      fwinfo,
		Master:         MasterEndPoint,
		Credential:     (*mesos.Credential)(nil),
		BindingAddress: parseIP(ServerIP),
	}

	driver, err := sched.NewMesosSchedulerDriver(schedConfig)

	if err != nil {
		log.Fatalf("Framework is not created error %v", err)
	}

	log.Printf("The Framework ID is %v and %v", fwinfo.Id, schedConfig.Framework.Id)

	status, err := driver.Run()

	if err != nil {
		log.Fatalf("Framework Run() error %v", err)
	}

	log.Printf("Framework Terminated with status %v", status.String())

}
