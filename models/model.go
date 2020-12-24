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
	Db.AutoMigrate(&Users{}, &Auths{}, &Follows{}, &Videos{}, &Favorites{}, &Categories{}, &Histories{})
	Db.Model(&Auths{}).AddForeignKey("username", "users(username)", "CASCADE", "CASCADE")
	Db.Model(&Follows{}).AddForeignKey("follower", "users(username)", "CASCADE", "CASCADE")
	Db.Model(&Follows{}).AddForeignKey("followee", "users(username)", "CASCADE", "CASCADE")
	Db.Model(&Videos{}).AddForeignKey("owner", "users(username)", "CASCADE", "CASCADE")
	Db.Model(&Categories{}).AddForeignKey("path", "videos(path)", "CASCADE", "CASCADE")
	Db.Model(&Favorites{}).AddForeignKey("path", "videos(path)", "CASCADE", "CASCADE")
	Db.Model(&Histories{}).AddForeignKey("path", "videos(path)", "CASCADE", "CASCADE")
	Db.Model(&Histories{}).AddForeignKey("video_name", "videos(name)", "CASCADE", "CASCADE")
}
