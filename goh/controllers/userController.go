package controllers

import (
	//"github.com/astaxie/beego"
	"goh/models"
	"time"
	"fmt"
	"errors"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"encoding/xml"
	"bytes"
)

type UserController struct {
	//beego.Controller
	models.Model
	models.Userinfo
	msgtype string 
	feid string 
}

type userinfo struct {
	XMLName xml.Name   `xml:"Name"`
	Version string     `xml:"version,attr"`
	Svs []models.Userinfo `xml:"userinfo"`
}


func (this *UserController) Prepare(){ //用户控制器需要admin权限
	models.CheckSession(&this.Model)

	access := this.GetSession("access").(int)
	if access<2 {
		//permission denied
		this.StopRun()
	}

	this.Username = this.GetString("username")//this.Ctx.Request.Form.Get("username")
	this.Permission,_ = this.GetInt("permission")

	//models.PrintoLog(fmt.Sprintf("this perrmiss %d",this.Permission))

	this.Password = models.GetMD5(this.GetString("password")) 
	this.Uid,_ = this.GetInt("uid")
	this.Created = time.Now().Format("2006-01-02 15:04:05")
	this.msgtype = this.GetString("msgtype")
	this.feid = this.GetString("feid")
}

func (this *UserController) Post(){
	switch this.msgtype {
	case "modify":
		//models.PrintoLog(this.feid)
		num := UpdateUserData(this.feid,this.Userinfo)
		this.Ctx.WriteString(fmt.Sprintf("影响行数：%d",num))
	case "deleteuser":
		o := orm.NewOrm()
		num,_ := o.Delete(&this.Userinfo)
		this.Ctx.WriteString(fmt.Sprintf("影响行数：%d",num))
	}
}

func (this *UserController) Get() {//getall user data; 
	// this.Data["Website"] = "localhost:8080/user"
	// this.Data["WebsiteReg"] = "localhost:8080/reg"
	// this.Data["WebsiteReg"] = "localhost:8080/reg"
	// this.TplName = "login.html"


	users,err := this.getAllUser()
	models.Printlog(err)

	v := &userinfo{Version: "1"}
	v.Svs = *users
	output, err := xml.Marshal(v)
	models.Printlog(err)

	bytesArry := new(bytes.Buffer)
	bytesArry.Write([]byte(xml.Header))
	bytesArry.Write(output)
	
	this.Ctx.ResponseWriter.Write(bytesArry.Bytes())
}

func (this *UserController) Delete(){ //删除用户 byname
	o := orm.NewOrm()
	_,err := o.QueryTable("userinfo").Filter("username",this.Username).Delete()
	if err != nil{
		this.Ctx.WriteString(fmt.Sprintf("删除%s的资料出错",this.Username))
	}
	this.Ctx.WriteString("资料已经删除") 
}

func (this *UserController) getAllUser() (*[]models.Userinfo,error){

	o := orm.NewOrm()

	users := new([]models.Userinfo)

	_, err := o.QueryTable("userinfo").All(users)
	//fmt.Printf("Returned Rows Num: %s, %s", num, err)
	return users,err
}

func (this *UserController) Put(){ //注册用户

	if this.Username == "" {
		this.Ctx.WriteString("用户名不能为空") 
		return 
	}
	if this.Password == "" {
		this.Ctx.WriteString("密码不能为空") 
		return 
	}

	name := this.getUserInfoByName(this.Username)
	if name.Username == this.Username {
		this.Ctx.WriteString("用户名已经存在") 
		return 
	}

	o := orm.NewOrm()
	_ , err := o.Insert(&this.Userinfo)
	if err != nil {
		//logs.Info("%s",err.Error()) //
		this.Ctx.WriteString(fmt.Sprintf("注册失败,错误代码:%s",err.Error())) 
	}
	this.Ctx.WriteString(fmt.Sprintf("%s注册用户",this.Username)) 
}

func (this *UserController) Patch() {//修改用户资料
	user := this.getUserInfoByName(this.Userinfo.Username)
	if this.Userinfo.Password != "" {
		user.Password = models.GetMD5(this.Userinfo.Password)
	}
	if this.Userinfo.Password != "" {
		user.Permission = this.Userinfo.Permission 
	}
	o := orm.NewOrm()
	_,err := o.Update(user)
	if err != nil{
		this.Ctx.WriteString(fmt.Sprintf("修改%s的资料出错,错误代码:%s",this.Username,err.Error())) 
	}else{
		this.Ctx.WriteString(fmt.Sprintf("修改%s的资料成功",this.Username))
	}
}

func UpdateUserData(feid string,user models.Userinfo) (int64){
	o := orm.NewOrm()
	num, err := o.Update(&user, feid)
	models.Printlog(err)
	return num
}

func  (this *UserController)UpdatePwd(name,pwd string) (error ) {//用户更改密码
	if name == "" {
		err := errors.New("用户名不能为空")
		return  err
	}
	if pwd == "" {
		err := errors.New("密码不能为空")
		return  err
	}
	user := this.getUserInfoByName(name)
	user.Password = models.GetMD5(pwd)
	o := orm.NewOrm()
	_ , err := o.Update(user, "password")
	return err 
}

func (this *UserController) getUserInfoByName(name string) (*models.Userinfo){
	user := new(models.Userinfo)
	o := orm.NewOrm()
	o.QueryTable("userinfo").Filter("username",name).One(user)
	return user
}