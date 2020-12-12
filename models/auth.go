package models

import (
	_ "github.com/jinzhu/gorm"
	"time"
)

type Users struct {
	UserId      uint   `gorm:"primary_key;auto-increment"`
	Identity    string `gorm:"default:'student'"`
	Username    string `gorm:"not null"`
	Follower    uint   `gorm:"default:0"`
	Followee    uint   `gorm:"default:0"`
	LaunchCount uint   `gorm:"default:0"`
	LikesCount  uint   `gorm:"default:0"`
	Location    string
	CreatedAt   time.Time
	DeletedAt   time.Time `gorm:"default:null"`
}

type Auths struct {
	UserId   uint `gorm:"ForeignKey:UserId"`
	Username string
	Password string
	Token    string
}

func CreateAuth(auth *Auths) (err error) {
	err = Db.Create(&auth).Error
	return
}

func CreateUser(user *Users) (err error) {
	err = Db.Create(&user).Error
	return
}

func QueryUser(username string, user *Users) (err error) {
	err = Db.Where(&Users{Username: username}).First(&user).Error
	return
}

func QueryAuth(username string, password string, auth *Auths) (err error) {
	err = Db.Where(&Auths{Username: username, Password: password}).First(&auth).Error
	return
}

func DeleteUser(username string) (err error) {
	err = Db.Where(&Users{Username: username}).Delete(&Users{}).Error
	return
}

//func (user *Auths) Update(attribute,value string) {
//	Db.Model(user).Update(attribute,value)
//}

//func (token *Token) Query() bool{
//	var user Auth
//	err := db.Where("token = ?",token.Token).First(&user).Error
//	return gorm.IsRecordNotFoundError(err)
//}
