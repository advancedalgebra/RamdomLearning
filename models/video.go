package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"time"
)

type Videos struct {
	VideoId   uint   `gorm:"primary_key;auto-increment"`
	Name      string `gorm:"not null;unique"`
	Owner     string `gorm:"not null;ForeignKey:Username"`
	Likes     uint   `gorm:"default:0"`
	Favorites uint   `gorm:"default:0"`
	Path      string `gorm:"not null;unique"`
	View      uint   `gorm:"default:0"`
	Comments  uint   `gorm:"default:0"`
	Forward   uint   `gorm:"default:0"`
	Size      uint   `gorm:"default:0"`
	Duration  time.Duration
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"default:null"`
}

type Categories struct {
	VideoId  uint `gorm:"ForeignKey:VideoId;primary_key"`
	Category string
	Dad      string
	Path     string `gorm:"ForeignKey:Path"`
	Level    uint   `gorm:"default:0"`
}

func QueryVideosByOwner(username string) (VideoList []*Videos, err error) {
	if err = Db.Select("path").Where(&Videos{Owner: username}).Find(&VideoList).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateForward(id uint) (err error) {
	err = Db.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"forward", gorm.Expr("forward + 1")).Error
	return err
}

func UpdateView(id uint) (err error) {
	err = Db.Model(Videos{}).Where(&Videos{VideoId: id}).Update("view", gorm.Expr("view + 1")).Error
	return err
}

func QueryHistory(userid, videoId uint) (history *Histories, err error) {
	history = new(Histories)
	if err = Db.Unscoped().Order("updated_at desc").Where(&Histories{VideoId: videoId, UserId: userid}).First(&history).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateNewViewToken(id uint, histories *Histories) (err error) {
	tx := Db.Begin()
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"view", gorm.Expr("view + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Histories{}).Create(&histories).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func UpdateDelViewToken(histories *Histories) (err error) {
	err = Db.Create(&histories).Error
	return
}

func UpdateOldViewToken(id, userid uint) (err error) {
	err = Db.Model(Histories{}).Where(&Histories{VideoId: id, UserId: userid}).Update(
		"count", gorm.Expr("count + 1")).Error
	return err
}

func LaunchTransaction(username string, video *Videos, categories *Categories) error {
	tx := Db.Begin()
	if err := tx.Model(Videos{}).Create(&video).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Categories{}).Create(&categories).Error; err != nil {
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

func UnLaunchTransaction(username string, id, count uint) error {
	tx := Db.Begin()
	//var video Videos
	//if err := tx.Where(&Videos{VideoId: id}).First(&video).Error; err != nil {
	//	return err
	//}
	if err := tx.Where(&Videos{VideoId: id}).Delete(&Videos{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Users{}).Where(&Users{Username: username}).Update(
		"launch_count", gorm.Expr("launch_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Users{}).Where(&Users{Username: username}).Update(
		"likes_count", gorm.Expr("likes_count - ?", count)).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&Favorites{VideoId: id}).Delete(&Favorites{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&Categories{VideoId: id}).Delete(&Categories{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&Histories{VideoId: id}).Delete(&Histories{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func LikeTransaction(id uint) error {
	tx := Db.Begin()
	var video Videos
	// 查出视频的作者
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

func DisLikeTransaction(id uint) error {
	tx := Db.Begin()
	var video Videos
	if err := tx.Where(&Videos{VideoId: id}).First(&video).Error; err != nil {
		return err
	}
	if err := tx.Model(Users{}).Where(&Users{Username: video.Owner}).Update(
		"Likes_count", gorm.Expr("Likes_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"likes", gorm.Expr("likes - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func QueryByCategory(category string) (videoList []*Categories, err error) {
	if err = Db.Select("path").Where(&Categories{Category: category}).Find(&videoList).Error; err != nil {
		return nil, err
	}
	return
}

func QueryTagById(id uint) (video *Categories, err error) {
	video = new(Categories)
	if err = Db.Where(&Categories{VideoId: id}).First(&video).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateVideoName(newName string, videoId uint) error {
	if err := Db.Model(Videos{}).Where(&Videos{VideoId: videoId}).Update("name", newName).Error; err != nil {
		return err
	}
	return nil
}

func QueryVideo(videoName string) (video *Videos, err error) {
	video = new(Videos)
	if err = Db.Unscoped().Where(&Videos{Name: videoName}).First(&video).Error; err != nil {
		return nil, err
	}
	return
}
