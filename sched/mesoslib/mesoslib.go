package mesoslib

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
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

	hostURI := fmt.Sprintf("http://%s:%d/%s", IP, Port, base)
	log.Printf("Hosting artifact '%s' at '%s'", path, hostURI)

	return &hostURI, base
}

func prepareExecutorInfo(IP, Port string) *mesos.ExecutorInfo {
	executorUris := []*mesos.CommandInfo_URI{}
	uri, executorCmd := serveExecutorArtifact(executorPath)
	executorUris = append(executorUris, &mesos.CommandInfo_URI{Value: uri, Executable: proto.Bool(true)})

	executorCommand := fmt.Sprintf("./%s -logtostderr=true ", executorCmd)

	go http.ListenAndServe(fmt.Sprintf("%s:%d", *address, *artifactPort), nil)
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

func parseConfig(config string) {

	splitconfig := strings.Split(config, ",")

	mIP := ""
	mP := "5050"
	sIP := ""
	sP := "5544"

	for i := 0; len(splitconfig); i++ {
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

func Run(config string) {

	var MasterIP string
	var MasterPort string
	var ServerIP string
	var ServerPort string

	//Split the configuration string

	MasterIP, MasterPort, ServerIP, ServerPort = parseConfig(config)

	//Get executor information
	exec := prepareExecutorInfo(ServerIP, ServerPort)

	// the framework
	fwinfo := &mesos.FrameworkInfo{
		User: proto.String(""), // Mesos-go will fill in user.
		Name: proto.String("MrRedis"),
	}

	//Add mesos authentication code
	//TODO

	//create the scheduler dirver object
	sched_config := sched.DriverConfig{
		Scheduler:      NewMrRedisScheduler(exec),
		Framework:      fwinfo,
		Master:         MasterIP + ":" + MasterPort,
		Credential:     nil,
		BindingAddress: ServerIP,
	}

	driver, err := sched.NewMesosSchedulerDriver(sched_config)

	if err != nil {
		log.Fatalf("Framework is not created error %v", err)
	}

	status, err := driver.Run()

	if err != nil {
		log.Fatalf("Framework Run() error %v", err)
	}

	log.Printf("Framework Terminated with status %v", status.String())

}
