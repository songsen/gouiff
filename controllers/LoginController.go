package controllers

import (
	//"time"
	//"github.com/astaxie/beego"
	"goh/models"
)


type LoginController struct {
	models.Model
}

func (this *LoginController) Get() {

	//sess_username := this.GetSession("username")
	sess_uid := this.GetSession("uid")
	if sess_uid == nil {
		this.Data["WebsiteLogin"] = "localhost:8080/login"
		this.TplName = "login.html"
		return 
	}

		
	event := this.GetString("loginEvent")
	switch event{
	case "logout":
		//this.DelSession("uid")
		this.DestroySession()
		this.Ctx.Redirect(302, "/login")
	default:
		this.Ctx.Redirect(302, "/")
		this.StopRun()
	}
}

func (this *LoginController) Post() {

	username := this.GetString("username");//this.Ctx.Request.Form.Get("username")
	password := this.GetString("password");//this.Ctx.Request.Form.Get("password")

	userinfo := models.GetUserInfo(username)
	newPass := models.GetMD5(password)

	if userinfo.Password == newPass {
		///models.LastLogin[username] = time.Now()
		this.SetSession("uid",userinfo.Uid)
		this.SetSession("username",userinfo.Username)
		this.SetSession("access",userinfo.Permission)
		//this.Ctx.WriteString(fmt.Sprintf("%s登录成功",username)) 
		// models.Dic(&this.Model)
		// this.TplName = "home.html"
		this.Ctx.Redirect(302, "/")
		//this.Ctx.Redirect(302, "/home")
	}else{
		this.Ctx.Redirect(302, "/login")
	}
}

