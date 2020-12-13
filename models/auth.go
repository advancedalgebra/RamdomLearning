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

type Follows struct {
	Follower  string `gorm:"primary_key"`
	Followee  string `gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt time.Time `gorm:"default:null"`
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

func QueryAuth(username, attribute, value string, auth *Auths) (err error) {
	err = Db.Where(map[string]interface{}{"username": username, attribute: value}).First(&auth).Error
	return
}

func DeleteUser(username string) (err error) {
	err = Db.Where(&Users{Username: username}).Delete(&Users{}).Error
	return
}

func DeleteAuth(username, password string) (err error) {
	err = Db.Where(&Auths{Username: username, Password: password}).Delete(&Auths{}).Error
	return
}

func UpdateAuth(attribute, value, username string) (err error) {
	err = Db.Model(Auths{}).Where(&Auths{Username: username}).Update(attribute, value).Error
	return err
}

func UpdateUser(attribute, value, username string) (err error) {
	err = Db.Model(Users{}).Where(&Users{Username: username}).Update(attribute, value).Error
	return err
}

//func (token *Token) Query() bool{
//	var user Auth
//	err := db.Where("token = ?",token.Token).First(&user).Error
//	return gorm.IsRecordNotFoundError(err)
//}
