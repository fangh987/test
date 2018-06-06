package controllers

import(
	"github.com/astaxie/beego"
	"testblog/models"
	"strings"
	"path"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get() {


	this.Data["IsTopic"] = true
	this.Data["IsLogin"] = CheckAccount(this.Ctx)
	this.TplName = "topic.html"
	Topics ,err := models.GetAllTopics("","",false)
	if err != nil {
		beego.Error(err)
	}
	this.Data["Topics"] = Topics
	

}
func (this *TopicController) Post() {
	if !CheckAccount(this.Ctx) {
		this.Redirect("/login",302)
		return
	}
	//解析表单
	title := this.Input().Get("title")
	content := this.Input().Get("content")
	tid := this.Input().Get("tid")
	category := this.Input().Get("category")
	label := this.Input().Get("label")
	//获取附件
	_,fh,err := this.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string

	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("attachment",path.Join("attachment",attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	
	
	//判断是否是添加文章还是修改文章
	if len(tid)==0 {
		err = models.AddTopic(title,content,category,label,attachment)
	} else {
		err = models.ModifyTopic(tid,title,category,content,label,attachment)
	}
	
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic",302)
	
}
func (this *TopicController) Add() {

	this.TplName = "topic_add.html"
}
func (this *TopicController) View() {
	this.TplName = "topic_view.html"
	tid := this.Input().Get("tid")
	topic,err := models.GetTopic(tid)
	//topic,err := models.GetTopic(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		this.Redirect("/",302)
		return
	}
	this.Data["Topic"] = topic
	this.Data["Tid"] = tid
	this.Data["Labels"] = strings.Split(topic.Labels," ")
	replies, err := models.GetAllReplies(tid)
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["Replies"] = replies
	this.Data["IsLogin"] = CheckAccount(this.Ctx)
	//this.Data["Tid"] = this.Ctx.Input.Param("0")
}
func (this *TopicController) Modify() {
	this.TplName = "topic_modify.html"
	 tid := this.Input().Get("tid")
	 topic,err := models.GetTopic(tid)
	 if err != nil {
		 beego.Error(err)
		 this.Redirect("/",302)
		 return
	 }
	 this.Data["Topic"] = topic
	 this.Data["Tid"] = tid
}
func (this *TopicController) Delete() {
	if !CheckAccount(this.Ctx) {
		this.Redirect("/login",302)
		return
	}
	err := models.DeleteTopic(this.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/",302)
}