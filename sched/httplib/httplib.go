package httplib

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
  "strings"

	"github.com/astaxie/beego"

	typ "github.com/mesos/mr-redis/common/types"
)

type InputData struct {
	Distribution int //Distribution Value
}

//MainController of the HTTP server
type MainController struct {
	beego.Controller
}

//Get Handles a general Get
func (this *MainController) Get() {
	this.Ctx.WriteString("hello world")
}

//CreateInstance Handles a Create Instance
func (this *MainController) CreateInstance() {

	var name string
	var capacity, masters, slaves int
	//Create an Default INPUT data
	IData := InputData{Distribution: 1}

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME")                  //Get the name of the instance
	capacity, _ = strconv.Atoi(this.Ctx.Input.Param(":CAPACITY")) // Get the capacity of the instance in MB
	masters, _ = strconv.Atoi(this.Ctx.Input.Param(":MASTERS"))   // Get the capacity of the instance in MB
	slaves, _ = strconv.Atoi(this.Ctx.Input.Param(":SLAVES"))     // Get the capacity of the instance in MB
	inData := this.Ctx.Input.CopyBody()

  //Check if instance name is valid, e.g.: space and null is not allowed from both front end and rest api
	name = strings.TrimSpace(name)

	if name == "null" {
		this.Ctx.WriteString(fmt.Sprintf("Instance name is null, please provide a valid name"))
		return
	}

	if len(name) == 0 {
		this.Ctx.WriteString(fmt.Sprintf("Instance name consists of spaces, please provide a valid name"))
		return
	}

	if len(inData) > 0 {
		//Some Payload is being supplied for create
		err := json.Unmarshal(inData, &IData)
		if err != nil {
			log.Printf("Invalid JSON format along wtih CREATE call IGNORING")
		}
	}

	log.Printf("Instance Name=%s, Capacity=%d, masters=%d, slaves=%d ConfigJson=%v\n", name, capacity, masters, slaves, IData)

	//Check the in-memory map if the instance already exist then return
	tmpInstance := typ.MemDb.Get(name)
	if tmpInstance == nil {
		tmpInstance = typ.LoadInstance(name)
	}

	//Check the central storage  if the instanc already exist then return

	if tmpInstance != nil {
		typ.MemDb.Add(name, tmpInstance)
		if tmpInstance.Status == typ.INST_STATUS_DELETED {

			this.Ctx.ResponseWriter.WriteHeader(201)
			this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, but in deleted state re-creating it", name))
		} else {
			this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be created", name))
			return
		}
	}

	//create a instance object
	instType := typ.INST_TYPE_SINGLE
	if slaves > 0 {
		instType = typ.INST_TYPE_MASTER_SLAVE
	}
	tmpInstance = typ.NewInstance(name, instType, masters, slaves, capacity)
	tmpInstance.Status = typ.INST_STATUS_CREATING
	tmpInstance.DistributionValue = IData.Distribution
	tmpInstance.Sync()

	ok, _ := typ.MemDb.Add(name, tmpInstance)
	if !ok {
		//It appears that the element is already there but in deleted state so update it
		typ.MemDb.Update(name, tmpInstance)
	}

	//Send it across to creator's channel
	typ.Cchan <- typ.CreateMaster(tmpInstance)

	//this.Ctx.Output.SetStatus(201)
	this.Ctx.ResponseWriter.WriteHeader(201)
	this.Ctx.WriteString("Request Accepted, Instance will be created.")
}

//DeleteInstance handles a delete instance REST call
func (this *MainController) DeleteInstance() {

	//var name string
	var name string

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	//Check the in-memory map if the instance already exists
	tmpInst := typ.MemDb.Get(name)
	if tmpInst == nil {
		tmpInst = typ.LoadInstance(name)
	}

	if tmpInst != nil {
		//get the instance data from central storage

		if tmpInst.Status == typ.INST_STATUS_DELETED {
			this.Ctx.ResponseWriter.WriteHeader(401)
			this.Ctx.WriteString(fmt.Sprintf("Instance %s is already deleted", name))
			return

		}

		//send info about all procs to be Destroyer to kill the master
		var tMsg typ.TaskMsg
		tMsg.P = tmpInst.Procs[tmpInst.Mname]
		tMsg.MSG = typ.TASK_MSG_DESTROY

		log.Printf("Destorying master %v from Instance %v", tMsg.P.ID, tmpInst.Name)

		//Send a message to the Destroyer
		typ.Dchan <- tMsg

		for _, n := range tmpInst.Snames {
			tMsg.P = tmpInst.Procs[n]
			if tMsg.P != nil {
				log.Printf("Destorying slave %v from Instance %v", tMsg.P.ID, tmpInst.Name)
			} else {
				log.Printf("Destroying Proc of the slave = %v is nil ", n)
			}

			//Send a message to the destroyer to kill the slaves
			typ.Dchan <- tMsg
		}
	} else {

		//The instance already exist return cannot create again return error
		this.Ctx.ResponseWriter.WriteHeader(401)
		this.Ctx.WriteString(fmt.Sprintf("Instance %s does not exist, cannot be deleted", name))
		return
	}

	//this.Ctx.Output.SetStatus(201)
	//ToDo: should this be blocking and the return happens when instance successfully deleted
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.WriteString("Request Placed for destroying")
}

//Status handles a STATUS REST call
func (this *MainController) Status() {

	//var name string
	var name string
	var inst *typ.Instance

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	//Check in memory map and store if the instance is available
	inst = typ.MemDb.Get(name)
	if inst == nil {
		inst = typ.LoadInstance(name)
		if inst == nil {
			this.Ctx.ResponseWriter.WriteHeader(501)
			this.Ctx.WriteString(fmt.Sprintf("Instance %s does not exist, error", name))
			return
		}
		typ.MemDb.Add(name, inst)
	}
	//not available in both the retrun error
	this.Ctx.WriteString(inst.ToJson())
}

//StatusAll handles StatusAll REST call
func (this *MainController) StatusAll() {

	var statusAll []typ.Instance_Json

	if len(typ.MemDb.I) == 0 {
		this.Ctx.WriteString("[]")
		return
	}

	for _, inst := range typ.MemDb.I {
		if inst.Status == typ.INST_STATUS_RUNNING {
			//statusAll = statusAll + inst.ToJson() + "\n"
			statusAll = append(statusAll, inst.ToJson_Obj())
		}
	}

	//not available in both the retrun error
	statusBytes, err := json.Marshal(statusAll)
	if err != nil {

		this.Ctx.WriteString("STATUSALL: Json Unmarshalling error")
		return
	}
	this.Ctx.WriteString(string(statusBytes))

}

//UpdateMemory Not yet implemented
func (this *MainController) UpdateMemory() {

	//var name string
	var name string

	//parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	//Check the instance in in-memory
	if !typ.MemDb.IsValid(name) {
		//The instance already exist return cannot create again return error
		this.Ctx.ResponseWriter.WriteHeader(501)
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//Check the instance in central storage

	//send the instance to Maintainer channel
	this.Ctx.WriteString("Upgrading the instance")

}

//UpdateSlaves Not yet implemented
func (this *MainController) UpdateSlaves() {

	//var name string
	var name string

	//parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instance

	//Check the instance in in-memory
	if typ.MemDb.IsValid(name) {
		//The instance already exist return cannot create again return error
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//Check the instance in central storage

	//send the instance to Maintainer channel
	this.Ctx.WriteString("Upgrading the instance slaves")

}

//Run main function that starts the HTTP server
func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Run(":" + config)

}
