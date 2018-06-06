package main

import (
	"testblog/controllers"
	_ "testblog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"testblog/models"
	"os"
)
func init() {
	models.RegisterDB()
}
func main() {
	orm.Debug = true
	orm.RunSyncdb("default",false,true)

	//创建附件目录
	os.Mkdir("attachment",os.ModePerm)
	/*
	方法一：附件的处理
	//作为静态文件
	beego.SetStaticPath("/attachment","attachment")
	*/
	//方法二：
	//作为单独一个控制器来处理
	beego.Router("/attachment/:all",&controllers.AttachController{})
	beego.Run()
}

