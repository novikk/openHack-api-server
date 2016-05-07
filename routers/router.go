package routers

import (
	"github.com/astaxie/beego"
	"openhack/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/customer/get", &controllers.UserController{}, "get:FindOrCreateUser")
	beego.Router("/customer/points", &controllers.UserController{}, "get:AddPoints")
	beego.Router("/product/buy", &controllers.UserController{}, "get:Buy")
	beego.Router("/products/list", &controllers.ProductController{}, "get:GetAllOfferings")
}
