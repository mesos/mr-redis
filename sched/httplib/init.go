package httplib

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	beego.Router("/v1/CREATE/:INSTANCENAME/:CAPACITY/:MASTERS/:SLAVES", &MainController{}, "post:CreateInstance")
	beego.Router("/v1/DELETE/:INSTANCENAME", &MainController{}, "delete:DeleteInstance")
	beego.Router("/v1/STATUS/:INSTANCENAME", &MainController{}, "get:Status")
	beego.Router("/v1/STATUS/", &MainController{}, "get:StatusAll")
	beego.Router("/v1/UPDATE/:INSTANCENAME/Memory/:CAPACITY", &MainController{}, "put:UpdateMemory")
	beego.Router("/v1/UPDATE/:INSTANCENAME/SLAVES/:SLAVES", &MainController{}, "put:UpdateSlaves")
}
