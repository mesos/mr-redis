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

func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Router("/", &MainController{})
	beego.Run()

}
