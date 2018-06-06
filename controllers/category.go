package controllers
import (
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
	"testblog/models"
)

type CategoryController struct {
	beego.Controller
}

func (this *CategoryController) Get() {
	op := this.Input().Get("op")
	var err error
	switch op {
	case "del":
		id := this.Input().Get("id")
		if len(id) == 0 {
			break
		}
		err = models.DelCategory(id)
		if err != nil {
			beego.Error(err)
		}
		this.Redirect("/category",301)
		return
		
	}
	this.Data["Categories"],err = models.GetAllCategories()

	this.Data["IsCategory"] = true
	this.Data["IsLogin"] = CheckAccount(this.Ctx)
	this.TplName = "category.html"
}

func (this *CategoryController) Post() {
	op := this.Input().Get("op")
	var err error
	switch op {
	case "add":
		name := this.Input().Get("name")
		if len(name) == 0 {
			break
		}
		err = models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		this.Redirect("/category",301)
		return
		
	}
	this.Data["Categories"],err = models.GetAllCategories()

	this.TplName = "category.html"
}

