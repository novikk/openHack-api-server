package main

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "openhack/cron"
	"openhack/models"
	_ "openhack/routers"
)

func main() {
	fmt.Println("Welcome to EcoRun")
	db := models.GetDB()
	db.AutoMigrate(&models.User{}, &models.PollutionPoint{})

	beego.Run()

	/*if !db.Debug().HasTable(&models.User{}) {
		db.Debug().CreateTable(&models.User{})
	}

	if !db.Debug().HasTable(&models.PollutionPoint{}) {
		db.Debug().CreateTable(&models.PollutionPoint{})
	}*/
}
