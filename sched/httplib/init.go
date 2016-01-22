package httplib

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/CREATE/:INSTANCENAME/:CAPACITY/:MASTERS/:SLAVES", &MainController{}, "post:CreateInstance")
	beego.Router("/v1/DELETE/:INSTANCENAME", &MainController{}, "delete:DeleteInstance")
	beego.Router("/v1/STATUS/:INSTANCENAME", &MainController{}, "get:Status")
	beego.Router("/v1/UPDATE/:INSTANCENAME/Memory/:CAPACITY", &MainController{}, "put:UpdateMemory")
	beego.Router("/v1/UPDATE/:INSTANCENAME/SLAVES/:SLAVES", &MainController{}, "put:UpdateSlaves")
}
