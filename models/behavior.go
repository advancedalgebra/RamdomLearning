package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Favorites struct {
	VideoId   uint   `gorm:"ForeignKey:VideoId;primary_key;auto_increment:false"`
	UserId    uint   `gorm:"ForeignKey:UserId;primary_key;auto_increment:false"`
	Path      string `gorm:"not null"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"default:null"`
}

func FavoriteTransaction(id uint, favorites *Favorites) error {
	tx := Db.Begin()
	//var video Videos
	//if err := tx.Where(&Videos{VideoId: id}).First(&video).Error; err != nil {
	//	return err
	//}
	if err := tx.Model(Favorites{}).Create(&favorites).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"favorites", gorm.Expr("favorites + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func DisFavoriteTransaction(id, uid uint) error {
	tx := Db.Begin()
	//var favorites Favorites
	//if err := tx.Where(&Favorites{VideoId: id}).First(&favorites).Error; err != nil {
	//	return err
	//}
	if err := tx.Where(&Favorites{VideoId: id, UserId: uid}).Delete(&Favorites{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"favorites", gorm.Expr("favorites - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func QueryPath(id uint) (video *Videos, err error) {
	video = new(Videos)
	if err = Db.Where(&Videos{VideoId: id}).First(&video).Error; err != nil {
		return nil, err
	}
	return
}

func QueryFavoritesByUserId(id uint) (favoritesList []*Favorites, err error) {
	if err = Db.Select("path").Where(&Favorites{UserId: id}).Find(&favoritesList).Error; err != nil {
		return nil, err
	}
	return
}
