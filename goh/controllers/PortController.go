package controllers

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/astaxie/beego"
	"goh/models"
	"time"
	"encoding/xml"
	"bytes"
	"fmt"
	"strings"
	"sort"
//	"reflect"
)

type PortController struct {
	models.Model
	models.Shippingline
	access int 
	tableName string 
	feild string 
	msgType string 
	msgValue string 
}

var (
	queryFrom = make([]string , 200)
	queryTo = make([]string , 200)

)

func init(){
	updateSet()
}

func updateSet(){
	tips := getPort("from")
	sort.Strings(*tips)
	queryFrom = Duplicate(*tips)

	tips = getPort("to")
	sort.Strings(*tips)
	queryTo =  Duplicate(*tips)
}

func updateQuerySet(strs *[]string ,str string){
	for _,agr := range *strs {
		if agr == str {
			return 
		}
	}
	*strs = append(*strs,str)
}

// func delDupliFromLine(strings []string , strs ...string) []string{
// 	var newstring []string 
// 	for index := 0; index < len(strings); index++ {
// 		startloop:
// 		for _,str := range strs{
// 			if strings[index] == str {
// 				goto startloop;
// 			}
// 		} 	
// 		newstring = append(newstring,strings[index])
// 	  }
// 	return newstring
// }

// func remove(slice []string, i int) []string {
//     //    copy(slice[i:], slice[i+1:])
//     //    return slice[:len(slice)-1]
//     return append(slice[:i], slice[i+1:]...)
// }


// type Shippingline struct {
// 	Id int  		`xml:"id"`
// 	Fromcity string `xml:"Fromcity"`
// 	Fromport string `xml:"Fromport"`
// 	Tocity string 	`xml:"Tocity"`
// 	Toport string 	`xml:"Toport"`
// 	Line string		`xml:"shippingLine"`
// 	Size int		`xml:"Size"`
// 	Price int 		`xml:"Price"`
// 	Updatetime string `xml:"Updatetime"`
// 	Note string		`xml:"Note"`
// 	Region string	`xml:"Region"`
// 	Person string 	`xml:"Person"`
// }

type Servers struct {
	XMLName xml.Name   `xml:"servers"`
	Version string     `xml:"version,attr"`
	Svs []models.Shippingline `xml:"server"`
}

func (this *PortController) Prepare(){
	//登录验证
	models.CheckSession(&this.Model)
	this.access = this.GetSession("access").(int)
	//提交表单
	this.tableName = "shippingline"
	//line := new(models.Shippingline)
	this.Fromcity = this.GetString("fromcity")
	this.Fromport = this.GetString("fromport")
	this.Tocity = this.GetString("tocity")
	this.Toport = this.GetString("toport")

	this.Price = this.GetString("price")
	// this.Ssize = this.GetString("ssize")
	// this.Xsize = this.GetString("xsize")
	// this.Xlsize = this.GetString("xlsize")
	 this.Ship= this.GetString("ship")
	
	//this.Ship = this.GetString("line")
	this.Mark = this.GetString("mark")

	this.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	this.Person = this.GetSession("username").(string)
	this.Region = this.GetString("region")
	this.Note = this.GetString("note")
	this.feild = this.GetString("feild")

	this.msgType = this.GetString("msgtype")
	this.msgValue = this.GetString("msgValue")

	this.Id,_ = this.GetUint16("uid")
	
}
//向chat控制器发送消息, chat会广播该消息
func (this *PortController) PublicMessage( content string ){

	uname := this.Info.Username

	if len(uname) == 0 || len(content) == 0 {
		return
	}

	publish <- newEvent(models.EVENT_MESSAGE, uname,content)
}

 //update/delete/add/ one line 
func (this *PortController) Post(){
	
		switch this.msgType{
		case "updatelinebyfeild":
			//models.PrintoLog(this.feild)
			num := UpdateLineData(this.feild,this.Shippingline)
			this.Ctx.WriteString(fmt.Sprintf("影响行数：%d",num))
			context := fmt.Sprintf("update,%d,%s,%s,%s,%s",this.Id,this.Fromport,this.Toport,this.Ship,this.GetString(this.feild))
			go this.PublicMessage(context)
		case "deleteline" :
			o := orm.NewOrm()
			num,_ := o.Delete(&this.Shippingline)
			this.Ctx.WriteString(fmt.Sprintf("影响行数：%d",num))
			go this.PublicMessage(fmt.Sprintf("delet,%d,%s,%s,%s",this.Id,this.Fromport,this.Toport,this.Ship))
			go updateSet()

		case "addoneline" :
			line,err := this.isLineExit(this.Fromport,this.Toport,this.Ship)
			if err != nil {
				this.Ctx.WriteString(fmt.Sprintf("isLineExit error :%s",err.Error()))
				return
			}
			if line.Id != 0 {
				this.Ctx.WriteString("add line error: line已存在")
				return
			}
		
			o := orm.NewOrm()
			_ , err = o.Insert(&this.Shippingline)
		
			if err != nil {
				this.Ctx.WriteString("add line error: 添加错误") 
			}else{
				this.Ctx.WriteString("add line ok") 
				updateQuerySet(&queryFrom,this.Fromcity)
				updateQuerySet(&queryFrom,this.Fromport)
				updateQuerySet(&queryTo,this.Toport)
				updateQuerySet(&queryTo,this.Tocity)
			}
			this.PublicMessage(fmt.Sprintf("add,%s,%s,%s",this.Fromport,this.Toport,this.Ship))
		}
	}

func UpdateLineData(feild string,line models.Shippingline) (int64){
	o := orm.NewOrm()
	num, err := o.Update(&line, feild)
	models.Printlog(err)
	return num
}


func (this *PortController) Get(){ //query line

	ok := this.IsAjax()
	if ok {
	
		switch this.msgType{
		case "fromport":
			// tips := getPort("from")
			// sort.Strings(*tips)
			// *tips = Duplicate(*tips)
			this.Ctx.WriteString(strings.Join(queryFrom, ","))
		case "toport":
			// tips := getPort("to")
			// sort.Strings(*tips)
			// *tips = Duplicate(*tips)
			this.Ctx.WriteString(strings.Join(queryTo, ","))
		// case "region":
		// 	tips := this.getTipByTagName("region")
		// 	sort.Strings(*tips)
		// 	*tips = Duplicate(*tips)
		// 	this.Ctx.WriteString(strings.Join(*tips, ","))
		case "person" ,"ship","fromcity","tocity","region": //tip 
			tips := this.getTipByTagName(this.msgType)
			sort.Strings(*tips)
			*tips = Duplicate(*tips)
			this.Ctx.WriteString(strings.Join(*tips, ","))
		case "getlinebyfeild"://

			lineinfo,err := this.getLineByField(this.feild,this.msgValue)
			models.Printlog(err)
		
			v := &Servers{Version: "1"}
			v.Svs = *lineinfo
			output, err := xml.Marshal(v)
			//output, err := xml.MarshalIndent(v, " ", " ")
			models.Printlog(err)
			
			bytesArry := new(bytes.Buffer)
			bytesArry.Write([]byte(xml.Header))
			bytesArry.Write(output)
			
			this.Ctx.ResponseWriter.Write(bytesArry.Bytes())

		case "getfromtoline"://query line 
					//models.PrintoLog("他是AJAX")
			from := this.GetString("from")
			to := this.GetString("to")
		
			lineinfo,err := this.getAllLine(from, to )
			models.Printlog(err)
		
			v := &Servers{Version: "1"}
			v.Svs = *lineinfo
			output, err := xml.Marshal(v)
			//output, err := xml.MarshalIndent(v, " ", " ")
			models.Printlog(err)
			
			bytesArry := new(bytes.Buffer)
			bytesArry.Write([]byte(xml.Header))
			bytesArry.Write(output)
			
			this.Ctx.ResponseWriter.Write(bytesArry.Bytes())
		}


	}else{
		//models.PrintoLog("他bu是AJAX")
		models.Dic(&this.Model)
		//this.Data["navline"] = "active"
		this.TplName = "line.html"
	}

}

	
	// func Duplicate(a interface{}) (ret []interface{}) {
	// 	va := reflect.ValueOf(a)
	// 	for i := 0; i < va.Len(); i++ {
	// 	   if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
	// 		  continue
	// 	   }
	// 	   ret = append(ret, va.Index(i).Interface())
	// 	}
	// 	return ret
	//  }
	
	 func Duplicate(a []string) []string {
		 //onesig := new([]string)
		 onesig := []string{"..."}
	
		//onesig[0] = a[0]
		i:=0
		 for _,s:= range a {
			if onesig[i] != s{
				onesig = append(onesig,s)
				i++
			}
		 }
	// now,pre,x:=0,0,0;
		return onesig
	 }

//该方法已经被弃用
func (this *PortController) Patch(){ //更新一条line's

	line,err := this.isLineExit(this.Fromport,this.Toport,this.Ship)
	if err != nil {
		this.Ctx.WriteString(fmt.Sprintf("isLineExit error :%s",err.Error()))
		return
	}

	if line.Id == 0 {
		this.Ctx.WriteString("add line error: line不存在")
		return
	}

	// if this.Ssize != "" {
	// 	line.Ssize = this.Ssize
	// }
	// if this.Xsize != "" {
	// 	line.Xsize = this.Xsize
	// }
	// if this.Xlsize != "" {
	// 	line.Xlsize = this.Xlsize
	// }

	if this.Price != "" {
		line.Price = this.Price
	}
	if this.Mark != ""{
		line.Mark = this.Mark
	}
	if this.Note != "" {
		line.Mark = this.Note
	}

	o := orm.NewOrm() 
	_ , err = o.Update(line)
	if err != nil {
		this.Ctx.WriteString(fmt.Sprintf("ORM.updte err: %s",err.Error())) 
	}else{ 
		this.Ctx.WriteString("update ok") 
	}
}

func getPort(flag string) *[]string{
	var desc,dest string 

	if flag == "from" {
		desc = "fromcity";
		dest = "fromport";
	}else if flag == "to" {
		desc = "tocity";
		dest = "toport";
	}else {
		return nil
	}

	o := orm.NewOrm()

	var lists []orm.ParamsList
	qs := o.QueryTable("shippingline")
	
	qs.ValuesList(&lists, desc, dest)
	qs.Distinct()

	for _, row := range lists {

		fmt.Printf("Name: %s, Age: %s", row[0], row[1])
	}

	paramSlice := new([]string)
    for _, param := range lists {
        *paramSlice = append(*paramSlice, param[0].(string),param[1].(string))
	}
	return paramSlice
}

func (this *PortController) getTipByTagName(tag string ) *[]string{

		o := orm.NewOrm()
		
		var list orm.ParamsList
		_, err := o.QueryTable("shippingline").ValuesFlat(&list, tag)
		if err != nil {
			models.Printlog(err)
			return nil
		}
		
		paramSlice := new([]string)
		for _, param := range list {
			*paramSlice = append(*paramSlice, param.(string))
		}

		return paramSlice
}

func (this *PortController) getLineByField(feild ,msgvalue string ) (*[]models.Shippingline,error){

	line := new([]models.Shippingline)

	o := orm.NewOrm()
	qs := o.QueryTable(this.tableName)

	qs = qs.Distinct()
	qs = qs.Filter(feild, msgvalue)
	
	_,err := qs.All(line)
	
	return line,err 
}

func (this *PortController) getAllLine(from,to string) (*[]models.Shippingline,error){
	var cond2 *orm.Condition 
	cond := orm.NewCondition()
	cond1 := orm.NewCondition()

	if from != "" {
		cond2 = cond1.AndCond(cond.Or("fromcity",from).Or("fromport",from))
	}

	if to != "" {
		cond2 = cond1.AndCond(cond.Or("tocity",to).Or("toport",to))
	}

	// if from =="" && to == "" {

	// }

	o := orm.NewOrm()
	qs := o.QueryTable(this.tableName)
    line := new([]models.Shippingline)
	_,err := qs.SetCond(cond2).Distinct().All(line)

	return line,err
}


func (this *PortController) isLineExit(fromport,toport string,ship string ) (*models.Shippingline,error){
	
		cond := orm.NewCondition()
		cond2 := cond.And("fromport",fromport).And("toport",toport).And("ship",ship)
	
		o := orm.NewOrm()
		qs := o.QueryTable(this.tableName)
	
		line := new(models.Shippingline)
		line.Id = 0
		_,err := qs.SetCond(cond2).Distinct().All(line)
		return line, err
}