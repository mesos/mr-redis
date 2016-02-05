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
	tmp_instance := typ.MemDb.Get(name)
	if tmp_instance == nil {
		tmp_instance = typ.LoadInstance(name)
	}

	//Check the central storage  if the instanc already exist then return

	if tmp_instance != nil {
		typ.MemDb.Add(name, tmp_instance)
		if tmp_instance.Status == typ.INST_STATUS_DELETED {

			this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, but in deleted state re-creating it", name))
		} else {
			this.Ctx.WriteString(fmt.Sprintf("Instance %s already exist, cannot be create", name))
			return
		}
	}

	//create a instance object
	inst_type := typ.INST_TYPE_SINGLE
	if slaves > 0 {
		inst_type = typ.INST_TYPE_MASTER_SLAVE
	}
	tmp_instance = typ.NewInstance(name, inst_type, masters, slaves, capacity)
	tmp_instance.Status = typ.INST_STATUS_CREATING

	tmp_instance.Sync()
	ok, _ := typ.MemDb.Add(name, tmp_instance)
	if !ok {
		//It appears that the element is already there but in deleted state so update it
		typ.MemDb.Update(name, tmp_instance)
	}

	//Send it across to creator's channel
	typ.Cchan <- typ.CreateMaster(tmp_instance)

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
	tmp_inst := typ.MemDb.Get(name)
	if tmp_inst == nil {
		tmp_inst = typ.LoadInstance(name)
	}
	if tmp_inst != nil {
		//get the instance data from central storage

		if tmp_inst.Status == typ.INST_STATUS_DELETED {
			this.Ctx.ResponseWriter.WriteHeader(401)
			this.Ctx.WriteString(fmt.Sprintf("Instance %s is already deleted", name))
			return

		}

		//send info about all procs to be Destroyer
		tmp_proc := tmp_inst.Procs[tmp_inst.Mname]

		log.Printf("Destorying master %v from Instance %v", tmp_proc.ID, tmp_inst.Name)

		typ.Dchan <- tmp_proc

		for _, n := range tmp_inst.Snames {
			tmp_proc = tmp_inst.Procs[n]
			if tmp_proc != nil {
				log.Printf("Destorying slave %v from Instance %v", tmp_proc.ID, tmp_inst.Name)
			} else {
				log.Printf("Destroying Proc of the slave = %v is nil ", n)
			}

			typ.Dchan <- tmp_proc
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

func (this *MainController) StatusAll() {

	//var name string
	var statusAll string

	for _, inst := range typ.MemDb.I {
		if inst.Status == typ.INST_STATUS_RUNNING {
			statusAll = statusAll + inst.ToJson() + "\n"
		}
	}

	//not available in both the retrun error
	this.Ctx.WriteString(statusAll)

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
