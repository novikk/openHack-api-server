package main

import (
	"github.com/astaxie/beego"
	_ "openhack/cron"
	"openhack/models"
	_ "openhack/routers"
)

func main() {
	beego.Run()

	db := models.GetDB()
	db.AutoMigrate(&models.User{})
}
