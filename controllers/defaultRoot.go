package controllers

import (
	//"github.com/astaxie/beego"
	"goh/models"
	//"fmt"
	"strings"
)

type MainController struct {
	models.Model
}

func (this *MainController) Prepare(){
	//登录验证
	models.CheckSession(&this.Model)
	//models.PrintoLog(this.Ctx.Input.Host())
	//this.Controller
}

func (this *MainController) Get() {
	// c.Data["Website"] = "beego.me"
	// c.Data["Email"] = "astaxie@gmail.com"
	// c.TplName = "index.tpl"
	//c.Ctx.ResponseWriter.ResponseWriter("xxx")
	//fmt.Fprintf(*c , )
	//c.Ctx.WriteString(beego.AppConfig.String("mysqluser"))
	// this.Data["WebsiteLogin"] = "localhost:8080/login"
	// this.Data["WebsiteLine"] = "localhost:8080/line"
	// this.Data["WebsiteUser"] = "localhost:8080/user"
	//ctx.Input.URL() 
	pathfile := this.Ctx.Input.URL()
	pathfile = strings.TrimPrefix(pathfile,"/")
	//models.PrintoLog(pathfile)
	models.Dic(&this.Model)
	this.Data["navhome"] = "active"
	if pathfile == "" {
		pathfile = "home.html"
	}
	this.TplName = pathfile
}
