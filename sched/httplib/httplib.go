package httplib

import (
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego"

	typ "../../common/types"
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
	if typ.MemDb.IsValid(name) {
		//The instance already exist return cannot create again return error
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//Check the central storage  if the instanc already exist then return
	tmp_instance := typ.LoadInstance(name)

	if tmp_instance != nil {
		typ.MemDb.Add(name, tmp_instance)
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//create a instance object
	inst_type := typ.INST_TYPE_SINGLE
	if slaves > 0 {
		inst_type = typ.INST_TYPE_MASTER_SLAVE
	}
	tmp_instance = typ.NewInstance(name, inst_type, masters, slaves, capacity)
	tmp_instance.Sync()
	tmp_instance.Status = typ.INST_STATUS_CREATING

	//Send it across to creator's channel
	typ.Cchan <- tmp_instance

	//this.Ctx.Output.SetStatus(201)
	this.Ctx.ResponseWriter.WriteHeader(201)
	this.Ctx.WriteString("Request Accepted, Instance will be created.")
}

func (this *MainController) DeleteInstance() {

	//var name string
	var name string

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instnace

	//Check the in-memory map if the instance already exists
	if typ.MemDb.IsValid(name) {
		//get the instance data from central storage
		tmp_inst := typ.LoadInstance(name)
		tmp_inst.LoadProcs()

		//send info about all procs to be Destroyer
		tmp_proc := tmp_inst.Procs[tmp_inst.Mname]

		log.Printf("SENT Destorying proc %v from Instance %v", tmp_proc.ID, tmp_proc.Instance)

		typ.Dchan <- tmp_proc

		for _, n := range tmp_inst.Snames {
			//ToDo is it needed to load the proc info from store also??
			//ToDo is there a delay needed between sending multiple values on this channel
			typ.Dchan <- tmp_inst.Procs[n]
		}

	} else {

		//The instance already exist return cannot create again return error
		this.Ctx.ResponseWriter.WriteHeader(401)
		this.Ctx.WriteString(fmt.Sprintf("Instance %s does not exist", name))
		return
	}

	//this.Ctx.Output.SetStatus(201)
	//ToDo: should this be blocking and the return happens when instance successfully deleted
	this.Ctx.ResponseWriter.WriteHeader(201)
	this.Ctx.WriteString("Request Placed for destroying")
}

func (this *MainController) Status() {

	//var name string
	var name string
	var inst *typ.Instance

	//Parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instnace

	//Check in memory map and store if the instance is available
	inst = typ.MemDb.Get(name)
	if inst == nil {
		inst = typ.LoadInstance(name)
		if inst == nil {
			this.Ctx.ResponseWriter.WriteHeader(501)
			this.Ctx.WriteString(fmt.Sprintf("Instance %s does not exist, error", name))
			return
		} else {
			typ.MemDb.Add(name, inst)
		}
	}

	//not available in both the retrun error
	this.Ctx.WriteString(inst.ToJson())

}

func (this *MainController) UpdateMemory() {

	//var name string
	var name string

	//parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instnace

	//Check the instnace in in-memory
	if !typ.MemDb.IsValid(name) {
		//The instance already exist return cannot create again return error
		this.Ctx.ResponseWriter.WriteHeader(501)
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//Check the instance in central storage

	//sedn the instnace to Maintainer channel
	this.Ctx.WriteString("Upgrading the instance")

}

func (this *MainController) UpdateSlaves() {

	//var name string
	var name string

	//parse the input URL
	name = this.Ctx.Input.Param(":INSTANCENAME") //Get the name of the instnace

	//Check the instnace in in-memory
	if typ.MemDb.IsValid(name) {
		//The instance already exist return cannot create again return error
		this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
		return
	}

	//Check the instance in central storage

	//sedn the instnace to Maintainer channel
	this.Ctx.WriteString("Upgrading the instance slaves")

}
func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Run()

}
