package main

import (
	_ "goh/routers"
	"github.com/astaxie/beego"
)

func main() {
	//beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}

