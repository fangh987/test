package controllers

import (
	"testblog/models"
	"github.com/astaxie/beego"
)
type ReplyController struct {
	beego.Controller
}
func (this *ReplyController) Add() {
	tid := this.Input().Get("tid")
	err := models.AddReply(tid,this.Input().Get("nickname"),this.Input().Get("content"))
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic/view?tid="+tid, 302)
}
func (this *ReplyController) Delete() {
	if !CheckAccount(this.Ctx) {
		return
	}
	tid := this.Input().Get("tid")
	err := models.DeleteReply(this.Input().Get("rid"))
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic/view?tid="+tid,302)
}