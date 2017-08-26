package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
)

type MainController struct {
	beego.Controller
}

type UserController struct {
	beego.Controller
}

type PostTestController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplName = "index.tpl"
}

func (this *UserController) Get() {
	this.Ctx.WriteString("this is in UserController get\n")
}

func (this *UserController) Post() {
	this.Ctx.WriteString("this is in UserController post\n")
}

func (this *PostTestController) Post() {
	key := this.GetString("key")
	id := this.GetString("id")
	fmt.Printf("key:%s\n", key)
	fmt.Printf("id:%s\n", id)
	this.Ctx.WriteString("this is in test post\n")
}

