package controllers

import (
	"github.com/astaxie/beego"
	"openhack/models"
)

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Get() {
	db := models.GetDB()
	var data []models.PollutionPoint

	db.Find(&data)

	c.Data["data"] = data
	c.TplName = "dashboard.tpl"
}
