package routers

import (
	"goh/controllers"
	"github.com/astaxie/beego"
	"goh/models"
)


func init() {
	//beego.InsertFilter("*", beego.BeforeRouter, models.FilterMethod)
	beego.InsertFilter("*", beego.BeforeRouter, models.FilterMethod)
	
	//controllers.UserController.Adduserinfo
	
//	beego.Router("/reg",&controllers.RegiterController{})
//	beego.Router("/del",&controllers.DelController{})
//	beego.Router("/chg",&controllers.ChgController{})
	// Long polling.
	// beego.Router("/lp", &controllers.LongPollingController{}, "get:Join")
	// beego.Router("/lp/fetch", &controllers.LongPollingController{}, "get:Fetch")

	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")

	beego.Router("/", &controllers.MainController{})
	beego.Router("/*.html", &controllers.MainController{})
	beego.Router("/login",&controllers.LoginController{})
	 beego.Router("/user",&controllers.UserController{})
	 beego.Router("/line",&controllers.PortController{})
	// beego.Router("/admin",&controllers.AdminCountroller{})
	//models.InitDB()
}
