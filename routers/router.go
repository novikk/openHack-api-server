package routers

import (
	"github.com/astaxie/beego"
	"openhack/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/customer/get", &controllers.UserController{}, "get:FindOrCreateUser")
	beego.Router("/customer/points", &controllers.UserController{}, "get:AddPoints")
}
