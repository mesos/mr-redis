package httplib

import (
	"github.com/astaxie/beego"
	"log"
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

	//Check the in-memory map if the instance already exist then return

	//Check the central storage  if the instanc already exist then return

	//create a instance object

	//Sync the instance object to central store

	//update the in memeory store

	//Send it across to creator's channel
	this.Ctx.WriteString("Request Placed for creating instance")
}

func (this *MainController) DeleteInstance() {

	var name string

	//Parse the input URL

	//Check the in-memory map if the instance does not exisy

	//Check the central storage if the instnace does not exist

	//Send it across to destroyers channel
	this.Ctx.WriteString("Request Placed for destroying")
}

func (this *MainController) Status() {

	var name string

	//parse the input URL

	//Check the in memory map if instnace avaiulable return the status in json

	//Check the central store if yes then return the status

	//not available in both the retrun error
	this.Ctx.WriteString("Status of the instance is ")

}

func (this *MainController) UpdateMemory() {

	var name string

	//parse the input URL

	//Check the instnace in in-memory

	//Check the instance in central storage

	//sedn the instnace to Maintainer channel
	this.Ctx.WriteString("Upgrading the instance")

}

func (this *MainController) UpdateSlaves() {

	var name string

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
