package models

import(
    "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"crypto/md5"
	"io"
	"bytes"
	"fmt"
	"github.com/astaxie/beego/context"
)

const userinfotable = `CREATE TABLE userinfo (
	uid INT(10) AUTO_INCREMENT,
	username VARCHAR(64) NOT NULL UNIQUE,
	password VARCHAR(64) NOT NULL ,
	permission INT(10) NOT NULL DEFAULT 0,
	created timestamp DEFAULT CURRENT_TIMESTAMP,
	index(username),
	PRIMARY KEY (uid ))`
//	update TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE ,
type Userinfo struct {
	Uid int	`orm:"PK"` 
	Username string  
	Password string 
	Permission int 
	Created string `orm:"auto_now_add"`
}
const shippinglinetable = `CREATE TABLE shippingline (
	id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
	fromcity VARCHAR(64)  DEFAULT "--",
	tocity VARCHAR(64)  DEFAULT "--",
	fromport VARCHAR(64) NOT NULL,
	toport VARCHAR(64)  NOT NULL,
	ship VARCHAR(32) NOT NULL,
	price VARCHAR(64) DEFAULT "--",
	mark VARCHAR(64) DEFAULT "--",
	person VARCHAR(16)  DEFAULT "--",
	region VARCHAR(16)  DEFAULT "--",
	note TEXT ,
	updatetime  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	index(fromcity,fromport,tocity,toport),
	PRIMARY KEY (id))`
type Shippingline struct {
	Id uint16  `orm:"PK" xml:"LineNum" `
	Fromcity string `xml:"FromCity"`
	Tocity string	`xml:"ToCity"`
	Fromport string `xml:"FromPort"`
	Toport string 	`xml:"ToPort"`
	Ship string     `xml:"Ship"`
	Price string	`xml:"S-L-XL"`
	Mark string		 `xml:"CloseDate"`
	Person string 	`xml:"By"`
	Region string
	Note string    
	Updatetime string `orm:"auto_now"`
}
//ALTER TABLE shippingline CHANGE `updatetime` `updatetime` TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER note  调整列顺序
type Model struct {
	beego.Controller
	Info Userinfo
}

func Dic(this *Model) {
	//this.Data["WebsiteLogin"] = "http://localhost:8080/login"
	//this.Data["WebsiteLine"] = "localhost:8080/line"
	//this.Data["WebsiteUser"] = "localhost:8080/user"
	//this.Data["WebsiteHome"] = "localhost:8080"
	//this.Data["WebLogOt"] = "http://localhost:8080/login?loginEvent=logout"
	//this.Data["WebAdmin"] = "http://localhost:8080/admin"
	//this.Data["Domain"] = "localhost:8080"

	if this.GetSession("username") != nil {
		usename := this.GetSession("username").(string)
		this.Data["optionUsername"] = usename
		access := this.GetSession("access").(int)

		this.Data["optionAdmub"] = false 
		this.Data["UserManagement"]=false 
		
		if access > 0{
			this.Data["optionAdmub"] = true 
		}
		if access > 1{
			this.Data["UserManagement"]=true 
		}
	}
}

func  PrintoLog(log string)  {
	logs.Info("model:%s",log)
}

func Printlog(err error){
if err != nil {
	logs.Info("%s",err.Error)
}
}
//"mysql","test:1007030237@tcp(songsen.top:3306)/test?charset=utf8"
func init(){

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "username:password@tcp(songsen.top:3306)/dbname?charset=utf8")
	orm.SetMaxIdleConns("default", 30)
	orm.SetMaxOpenConns("default", 30)
	orm.RegisterModel(new(Userinfo),new(Shippingline))
	//orm.RunSyncdb("userinfo", false, true)
	orm.Debug = true
	//orm.RunSyncdb()
}

var FilterMethod = func(ctx *context.Context) {

	if ctx.Input.Query("_method")!="" && ctx.Input.IsPost(){
		  ctx.Request.Method = ctx.Input.Query("_method")
	}
}

func GetMD5(password string ) string {
	md5Password := md5.New()
	io.WriteString(md5Password, password)
	buffer := bytes.NewBuffer(nil)
	fmt.Fprintf(buffer, "%x", md5Password.Sum(nil))
	//newPass := buffer.String()
	return password
}


func CheckSession(this *Model) {
	sess_uid := this.GetSession("uid")

	//sess_username := this.GetSession("username")
	if sess_uid == nil {

		this.Ctx.Redirect(302, "/login")
		this.StopRun()
	}

	this.Info.Username = this.GetSession("username").(string)
	this.Info.Permission = this.GetSession("access").(int)

//this.Data["Username"] = sess_username
}

func GetUserInfo(name string) (Userinfo){
	var user Userinfo
	o := orm.NewOrm()
	//o.Using("userinfo")
	o.QueryTable("userinfo").Filter("username",name).One(&user)
	return user
}

// func Session(this *beego.Controller) {
// 	username := this.GetString("username");//this.Ctx.Request.Form.Get("username")
// 	password := this.GetString("password");//this.Ctx.Request.Form.Get("password")

// 	userinfo := GetUserInfo(username)
// 	newPass := GetMD5(password)
	
// 	if userinfo.Password == newPass {
// 		///models.LastLogin[username] = time.Now()
// 		this.SetSession("uid",userinfo.Uid)
// 		this.SetSession("username",userinfo.Username)
// 		//this.Ctx.WriteString(username +" password ok") 
// 	}else{
// 		this.Ctx.Redirect(302, "/login")
// 		//this.Ctx.WriteString(username +" password error") 
// 		this.StopRun()
// 	}
// }

// UPDATE：POST
// READ：GET
// CREATE ：PUT   //add
// DELETE：DELETE
// Patch()  //update更新数据,若没有就增添

// Prepare()
// Get() 	//： 请求指定的页面信息，并返回实体主体
// Post() 
// Delete() //请求服务器删除指定的页面。
// Put() 
// Head()   //HEAD： 只请求页面的首部。
// Patch()  //update更新数据,若没有就增添
// Options() //允许客户端查看服务器的性能。
// Finish()

// func  UpdateLinePirce(line *Shippingline) error{
// 	o := orm.NewOrm() //.Update(orm.Params("price": line.Price,))
// 	//o.QueryTable("shippingline").Filter("fromport",line.Fromport).Filter("toport",line.Toport).Filter("size",line.Size).Filter("line",line.Line).Update("price",line.Price)
// 	// _,err :=o.QueryTable("shippingline").Filter("fromport",line.Fromport).Filter("toport",line.Toport).Filter("size",line.Size).Filter("line",line.Line).Update(orm.Params{
// 	// 	"price": line.Price,
// 	// })
// 	_ , err := o.Update(line, "price")
// 		// qs := o.QueryTable("shippingline")
// 		// qs.SetCond().UpdateUpdatePermission
// 	//_, err :=o.Update(line)
// 	return err
// }

// func InsetLine(line *Shippingline) error{
// 	o := orm.NewOrm()
// 	_ , err := o.Insert(line)
// 	Printlog(err)
// 	return err
// }



// func QueryLine(from,to string) []Shippingline{
// 	cond := orm.NewCondition()
// 	condfrom := cond.Or("fromcity",from).Or("fromport",from)

// 	condto := cond.Or("tocity",to).Or("toport",to)
	
// 	cond2 := cond.AndCond(condfrom).AndCond(condto)
// 	o := orm.NewOrm()
// 	qs := o.QueryTable("shippingline")
// 	var line []Shippingline
// 	_,err := qs.SetCond(cond2).Distinct().All(&line)
// 	Printlog(err )
// 	return line
// }



// var (
// 	LastLogin map[string]time.Time
// )



// func DelteUser(name string ) {
// 	o := orm.NewOrm()
// 	_,err := o.QueryTable("userinfo").Filter("username",name).Delete()
// 	Printlog(err)
// }

// func Adduserinfo(info *Userinfo ) ( error ) {

// 	if info.Username == "" {
// 		err := errors.New("用户名不能为空")
// 		return  err
// 	}
// 	if info.Password == "" {
// 		err := errors.New("密码不能为空")
// 		return  err
// 	}
// 	name := GetUserInfo(info.Username)
// 	if name.Username == info.Username {
// 		err := errors.New("用户名已经存在")
// 		return  err
// 	}

// 	o := orm.NewOrm()
// 	//err := o.Using("userinfo")
// 	// if err != nil {
// 	// 	logs.Info("%s ",err.Error)
// 	// }
// 	_ , err := o.Insert(info)
// 	Printlog(err)
// 	return err
// }

// func UpdatePwd(info *Userinfo) (error ) {

// 	o := orm.NewOrm()
// 	//_,err := o.Update(info, "password")
// 	//Printlog(err)
// 	//o.Using("userinfo")
// 	_ , err := o.Update(info, "password")
// 	//_, err := o.QueryTable("userinfo").Update(info, "password")

// 	// _, err := o.QueryTable("userinfo").Filter("name", "slene").Update(orm.Params{
// 	// 	"name": "astaxie",
// 	// })

// 	Printlog(err)
// 	return err

// }

// func UpdatePermission(info *Userinfo)  (error ){
// 	//info := Userinfo{Permission:i}
// 	o := orm.NewOrm()
// 	//o.Using("userinfo")
// 	_ , err := o.Update(info, "permission")
// 	Printlog(err)
// 	return err
// }



