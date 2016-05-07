package models

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/postgres"

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost dbname=ecorun user=postgres password=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
