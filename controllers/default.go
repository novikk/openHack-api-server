package controllers

import (
	"github.com/astaxie/beego"
	"os"
)

var TMF_CUSTOMER_KEY string
var TMF_CATALOG_KEY string
var TMF_INVENTORY_KEY string

type MainController struct {
	beego.Controller
}

func init() {
	TMF_CUSTOMER_KEY = os.Getenv("TMF_CUSTOMER_KEY")
	TMF_CATALOG_KEY = os.Getenv("TMF_CATALOG_KEY")
	TMF_INVENTORY_KEY = os.Getenv("TMF_INVENTORY_KEY")
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
