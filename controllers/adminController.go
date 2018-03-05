package controllers

import (
	//"github.com/astaxie/beego"
	"goh/models"
	//"fmt"
)

type AdminCountroller struct {
	models.Model
}

func (this *AdminCountroller) Prepare(){
	//登录验证
	models.CheckSession(&this.Model)
	//this.Controller
}

func (this *AdminCountroller) Get() {
	// c.Data["Website"] = "beego.me"
	// c.Data["Email"] = "astaxie@gmail.com"
	// c.TplName = "index.tpl"
	//c.Ctx.ResponseWriter.ResponseWriter("xxx")
	//fmt.Fprintf(*c , )
	//c.Ctx.WriteString(beego.AppConfig.String("mysqluser"))
	// this.Data["WebsiteLogin"] = "localhost:8080/login"
	// this.Data["WebsiteLine"] = "localhost:8080/line"
	// this.Data["WebsiteUser"] = "localhost:8080/user"
	//
	models.Dic(&this.Model)
	this.Data["navhome"] = "active"

	this.TplName = "admin.html"
}