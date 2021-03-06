package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)
type LoginController struct {
	beego.Controller
}
func (this *LoginController) Get() {
	IsExit := this.Input().Get("exit") == "true"
	if IsExit {
		this.Ctx.SetCookie("uname","",-1,"/")
		this.Ctx.SetCookie("pwd","",-1,"/")
		this.Redirect("/",301)
		return
	}
	this.TplName = "login.html"
}

func (this *LoginController) Post() {
	uname := this.Input().Get("uname")
	pwd := this.Input().Get("pwd")
	autoLogin := this.Input().Get("autoLogin") == "on"
	if uname == beego.AppConfig.String("account") && pwd == beego.AppConfig.String("password") {
		maxAge := 0
		//自动登录
		if autoLogin {
			maxAge = 1<<31 - 1

		}
		this.Ctx.SetCookie("uname",uname,maxAge,"/")
		this.Ctx.SetCookie("pwd",pwd,maxAge,"/")
		
	}
	this.Redirect("/",301)
	return
}

func CheckAccount(c  *context.Context) bool {
	ck,err := c.Request.Cookie("uname")
	if err != nil {
		return false
	}
	uname := ck.Value
	ck,err = c.Request.Cookie("pwd")
	if err != nil {
		return false
	}
	pwd := ck.Value
	return uname == beego.AppConfig.String("account") && pwd == beego.AppConfig.String("password")
}