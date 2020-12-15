package models

import (
	"RamdomLearning/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open(conf.DBType, conf.DbUrl)
	if err != nil {
		log.Println(err)
	}
	Db.SingularTable(true)
	Db.AutoMigrate(&Auths{}, &Users{}, &Follows{}, &Videos{})
}
