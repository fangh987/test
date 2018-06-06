package controllers

import (
	"github.com/astaxie/beego"
	"testblog/models"
)

type HomeController struct {
	beego.Controller
}
//首页
func (this *HomeController) Get() {
	this.Data["IsHome"] = true
	this.Data["IsLogin"] = CheckAccount(this.Ctx)
	this.TplName = "home.html"
	Topics ,err := models.GetAllTopics(this.Input().Get("cate"),this.Input().Get("label"),true)
	if err != nil {
		beego.Error(err)
	}
	this.Data["Topics"] = Topics

	categories ,err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	this.Data["Categories"] = categories



}
