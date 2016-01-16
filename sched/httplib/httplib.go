package httplib

import (
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego"

	"../../common/types"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Ctx.WriteString("hello world")
}

func (this *MainController) CreateInstance() {

	var name string
	var capacity, masters, slaves int

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME")                  //Get the name of the instnace
	capacity, _ = strconv.Atoi(this.Ctx.Input.Param(":CAPACITY")) // Get the capacity of the instance in MB
	masters, _ = strconv.Atoi(this.Ctx.Input.Param(":MASTERS"))   // Get the capacity of the instance in MB
	slaves, _ = strconv.Atoi(this.Ctx.Input.Param(":SLAVES"))     // Get the capacity of the instance in MB

	log.Printf("Instance Name=%s, Capacity=%d, masters=%d, slaves=%d\n", name, capacity, masters, slaves)

	//Check the in-memory map if the instance already exist then return
	if types.MemDb.IsValid(name) {
		//The instance already exist return cannot create again return error
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//Check the central storage  if the instanc already exist then return
	tmp_instance := types.LoadInstance(name)

	if tmp_instance != nil {
		types.MemDb.Add(name, tmp_instance)
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//create a instance object
	tmp_instance = types.NewInstance(name, "S", masters, slaves, capacity)
	tmp_instance.Sync()
	tmp_instance.Status = "STARTING"

	//Send it across to creator's channel
	types.Cchan <- tmp_instance

	//this.Ctx.Output.SetStatus(201)
	this.Ctx.ResponseWriter.WriteHeader(201)
	this.Ctx.WriteString("Request Accepted, Instance will be created.")
}

func (this *MainController) DeleteInstance() {

	//var name string

	//Parse the input URL

	//Check the in-memory map if the instance does not exisy

	//Check the central storage if the instnace does not exist

	//Send it across to destroyers channel
	this.Ctx.WriteString("Request Placed for destroying")
}

func (this *MainController) Status() {

	//var name string

	//parse the input URL

	//Check the in memory map if instnace avaiulable return the status in json

	//Check the central store if yes then return the status

	//not available in both the retrun error
	this.Ctx.WriteString("Status of the instance is ")

}

func (this *MainController) UpdateMemory() {

	//var name string

	//parse the input URL

	//Check the instnace in in-memory

	//Check the instance in central storage

	//sedn the instnace to Maintainer channel
	this.Ctx.WriteString("Upgrading the instance")

}

func (this *MainController) UpdateSlaves() {

	//var name string

	//parse the input URL

	//Check the instnace in in-memory

	//Check the instance in central storage

	//sedn the instnace to Maintainer channel
	this.Ctx.WriteString("Upgrading the instance slaves")

}
func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Run()

}
