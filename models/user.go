package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"time"
)

type Users struct {
	UserId      uint   `gorm:"primary_key;auto-increment"`
	Identity    string `gorm:"default:'student'"`
	Username    string `gorm:"not null;unique"`
	Follower    uint   `gorm:"default:0"`
	Following   uint   `gorm:"default:0"`
	LaunchCount uint   `gorm:"default:0"`
	LikesCount  uint   `gorm:"default:0"`
	Location    string
	CreatedAt   time.Time
	DeletedAt   *time.Time `gorm:"default:null"`
}

type Auths struct {
	UserId   uint   `gorm:"ForeignKey:UserId;primary_key"`
	Username string `gorm:"ForeignKey:Username"`
	Password string
	Token    string
}

type Follows struct {
	FollowId  uint   `gorm:"primary_key;auto-increment"`
	Follower  string `gorm:"ForeignKey:Username"`
	Followee  string `gorm:"ForeignKey:Username"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"default:null"`
}

//func CreateAuth(auth *Auths) (err error) {
//	err = Db.Create(&auth).Error
//	return
//}
//
//func CreateUser(user *Users) (err error) {
//	err = Db.Create(&user).Error
//	return
//}

func QueryUserAll(username string) (user *Users, err error) {
	user = new(Users)
	if err = Db.Unscoped().Where(&Users{Username: username}).First(&user).Error; err != nil {
		return nil, err
	}
	//print(user.Username)
	return
}

func QueryAuth(username, attribute, value string) (auth *Auths, err error) {
	auth = new(Auths)
	if err = Db.Where(map[string]interface{}{"username": username, attribute: value}).First(&auth).Error; err != nil {
		return nil, err
	}
	return
}

//func QueryFollowItem(username, follower string) (follow *Follows, err error) {
//	follow = new(Follows)
//	if err = Db.Where(map[string]interface{}{
//		"followee": username, "follower": follower}).First(&follow).Error; err != nil {
//		return nil, err
//	}
//	return
//}

func QueryFollower(username string) (followerList []*Follows, err error) {
	if err = Db.Select("followee").Where(&Follows{Follower: username}).Find(&followerList).Error; err != nil {
		return nil, err
	}
	return
}

func QueryFollowing(username string) (followerList []*Follows, err error) {
	if err = Db.Select("follower").Where(&Follows{Followee: username}).Find(&followerList).Error; err != nil {
		return nil, err
	}
	return
}

//func QueryAuthById(id uint, auth *Auths) (err error) {
//	err = Db.Where(map[string]interface{}{"user_id": id}).First(&auth).Error
//	return
//}

//func DeleteUser(username string, db *gorm.DB) (err error) {
//	err = db.Where(&Users{Username: username}).Delete(&Users{}).Error
//	return
//}
//
//func DeleteAuth(username, password string, db *gorm.DB) (err error) {
//	err = db.Where(&Auths{Username: username, Password: password}).Delete(&Auths{}).Error
//	return
//}

func UpdateAuth(attribute, value, username string) (err error) {
	err = Db.Model(Auths{}).Where(&Auths{Username: username}).Update(attribute, value).Error
	return err
}

func UpdateUser(attribute, value, username string) (err error) {
	err = Db.Model(Users{}).Where(&Users{Username: username}).Update(attribute, value).Error
	return err
}

func Commit(args ...interface{}) error {
	tx := Db.Begin()
	for _, v := range args {
		if err := tx.Create(v).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

//func UpdateTransaction(username, NewName string) error {
//	tx := Db.Begin()
//	if err := tx.Model(Auths{}).Where(&Auths{Username: username}).Update("username", NewName).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	if errAuth := tx.Model(Users{}).Where(&Users{Username: username}).Update(
//		"username", NewName).Error; errAuth != nil {
//		tx.Rollback()
//		return errAuth
//	}
//	tx.Commit()
//	return nil
//}

func UpdateTransaction(username, NewName string) error {
	if err := Db.Model(Users{}).Where(&Users{Username: username}).Update("username", NewName).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTransaction(username string, id uint) error {
	tx := Db.Begin()
	if err := tx.Where(&Users{Username: username}).Delete(&Users{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if errAuth := tx.Where(&Auths{Username: username}).Delete(&Auths{}).Error; errAuth != nil {
		tx.Rollback()
		return errAuth
	}
	var followerList []*Follows
	if err := Db.Select("followee").Where(&Follows{Follower: username}).Find(&followerList).Error; err != nil {
		return err
	}
	for _, v := range followerList {
		if err := tx.Model(Users{}).Where(&Users{Username: v.Followee}).Update(
			"following", gorm.Expr("following - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where(&Histories{UserId: id}).Delete(&Histories{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := Db.Select("follower").Where(&Follows{Followee: username}).Find(&followerList).Error; err != nil {
		return err
	}
	for _, v := range followerList {
		if err := tx.Model(Users{}).Where(&Users{Username: v.Follower}).Update(
			"follower", gorm.Expr("follower - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where(&Follows{Follower: username}).Delete(&Follows{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&Follows{Followee: username}).Delete(&Follows{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&Favorites{UserId: id}).Delete(&Favorites{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	var VideoList []*Videos
	if err := tx.Select("video_id").Where(&Videos{Owner: username}).Find(&VideoList).Error; err != nil {
		return err
	}
	for _, v := range VideoList {
		//println(v.VideoId)
		if err := tx.Where(&Favorites{VideoId: v.VideoId}).Delete(&Favorites{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where(&Categories{VideoId: v.VideoId}).Delete(&Categories{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where(&Histories{VideoId: v.VideoId}).Delete(&Histories{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where(&Videos{Owner: username}).Delete(&Videos{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	var CommentList []*Comments
	if err := Db.Where(&Comments{Commenter: username, Type: "video"}).Find(&CommentList).Error; err != nil {
		return err
	}
	for _, v := range CommentList {
		if err := DeleteCommentVideo(v.Origin, v.CommentId, v.Count); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := Db.Where(&Comments{Commenter: username, Type: "comment"}).Find(&CommentList).Error; err != nil {
		return err
	}
	for _, v := range CommentList {
		if err := DeleteComment(v.Count, v.CommentId, v.Origin); err != nil {
			tx.Rollback()
			return err
		}
	}
	//tx.Rollback()
	tx.Commit()
	return nil
}

//func DeleteTransaction(username string, id uint) error {
//	tx := Db.Begin()
//	if err := tx.Where(&Users{Username: username}).Delete(&Users{}).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	if errAuth := tx.Where(&Auths{Username: username}).Delete(&Auths{}).Error; errAuth != nil {
//		tx.Rollback()
//		return errAuth
//	}
//	if err := tx.Where(&Follows{Follower: username}).Delete(&Follows{}).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	if err := tx.Where(&Follows{Followee: username}).Delete(&Follows{}).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	if err := tx.Where(&Favorites{UserId: id}).Delete(&Favorites{}).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	if err := tx.Where(&Videos{Owner: username}).Delete(&Videos{}).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	tx.Commit()
//	return nil
//}

//func DeleteTransaction(username string) error {
//	if err := Db.Where(&Users{Username: username}).Delete(&Users{}).Error; err != nil {
//		return err
//	}
//	return nil
//}

func FollowTransaction(username, follower string) error {
	tx := Db.Begin()
	if err := tx.Model(Users{}).Where(&Users{Username: username}).Update(
		"following", gorm.Expr("following + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Users{}).Where(&Users{Username: follower}).Update(
		"follower", gorm.Expr("follower + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&Follows{Follower: follower, Followee: username}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func UnFollowTransaction(username, follower string) error {
	tx := Db.Begin()
	if err := tx.Model(Users{}).Where(&Users{Username: username}).Update(
		"following", gorm.Expr("following - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Users{}).Where(&Users{Username: follower}).Update(
		"follower", gorm.Expr("follower - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&Follows{Followee: username, Follower: follower}).Delete(&Follows{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
