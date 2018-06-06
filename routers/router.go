package routers

import (
	"testblog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//首页
	beego.Router("/",&controllers.HomeController{})
	//登录
	beego.Router("/login",&controllers.LoginController{})
	//分类
	beego.Router("/category",&controllers.CategoryController{})
	//评论
	beego.Router("/reply",&controllers.ReplyController{})
	beego.Router("/reply/add",&controllers.ReplyController{},"post:Add")
	beego.Router("/reply/delete",&controllers.ReplyController{},"get:Delete")
	//文章
	beego.Router("/topic",&controllers.TopicController{})
	//自动路由
	beego.AutoRouter(&controllers.TopicController{})

	
	
}
