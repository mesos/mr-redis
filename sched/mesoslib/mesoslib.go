package mesoslib

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/auth"
	"github.com/mesos/mesos-go/auth/sasl"
	"github.com/mesos/mesos-go/auth/sasl/mech"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"golang.org/x/net/context"

	"../../common/types"
)

func Run(config string) {

	var ServerIP string

	// the framework
	fwinfo := &mesos.FrameworkInfo{
		User: proto.String(""), // Mesos-go will fill in user.
		Name: proto.String("MrRedis"),
	}

	//Add mesos authentication code

	//create the scheduler dirver object
	sched_config := sched.DriverConfig{
		Scheduler:      newExampleScheduler(exec),
		Framework:      fwinfo,
		Master:         *master,
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
