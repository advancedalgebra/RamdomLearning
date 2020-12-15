package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"time"
)

type Videos struct {
	VideoId   uint   `gorm:"primary_key;auto-increment"`
	Name      string `gorm:"not null"`
	Owner     string `gorm:"not null"`
	Likes     uint   `gorm:"default:0"`
	Path      string `gorm:"not null"`
	View      uint   `gorm:"default:0"`
	Forward   uint   `gorm:"default:0"`
	Size      uint   `gorm:"default:0"`
	Duration  time.Duration
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"default:null"`
}

type Tags struct {
	VideoId uint `gorm:"ForeignKey:UserId;primary_key"`
	Tag     string
	Dad     string
	Level   uint `gorm:"default:0"`
}

func LaunchTransaction(username string, video *Videos) error {
	tx := Db.Begin()
	if err := tx.Model(Videos{}).Create(&video).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Users{}).Where(&Users{Username: username}).Update(
		"launch_count", gorm.Expr("launch_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func LikeTransaction(id uint) error {
	tx := Db.Begin()
	var video Videos
	if err := tx.Where(&Videos{VideoId: id}).First(&video).Error; err != nil {
		return err
	}
	if err := tx.Model(Users{}).Where(&Users{Username: video.Owner}).Update(
		"Likes_count", gorm.Expr("Likes_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"likes", gorm.Expr("likes + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func UpdateForward(id uint) (err error) {
	err = Db.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"forward", gorm.Expr("forward + 1")).Error
	return err
}

func UpdateView(id uint) (err error) {
	err = Db.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"view", gorm.Expr("view + 1")).Error
	return err
}

//func QueryByTag(username string) (followerList []*Follows, err error) {
//	if err = Db.Select("followee").Where(&Follows{Follower: username}).Find(&followerList).Error; err != nil {
//		return nil, err
//	}
//	return
//}
