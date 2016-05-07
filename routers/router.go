package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"openhack/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/customer/get", &controllers.UserController{}, "get:FindOrCreateUser")
	beego.Router("/customer/points", &controllers.UserController{}, "get:AddPoints")
	beego.Router("/customer/inventory", &controllers.UserController{}, "get:GetCustomerInventory")
	beego.Router("/products/buy", &controllers.UserController{}, "get:Buy")
	beego.Router("/products/list", &controllers.ProductController{}, "get:GetAllOfferings")
	beego.Router("/route/new", &controllers.RouteController{}, "get:NewRoute")
	beego.Router("/route/new_with_end", &controllers.RouteController{}, "get:NewRouteWithEnd")
	beego.Router("/dashboard", &controllers.DashboardController{})

	var FilterCORS = func(ctx *context.Context) {
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
	}

	beego.InsertFilter("*", beego.BeforeRouter, FilterCORS)
}
