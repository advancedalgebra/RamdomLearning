package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Favorites struct {
	VideoId   uint   `gorm:"ForeignKey:VideoId;primary_key;auto_increment:false"`
	UserId    uint   `gorm:"ForeignKey:UserId;primary_key;auto_increment:false"`
	Path      string `gorm:"not null;ForeignKey:Path"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"default:null"`
}

type Histories struct {
	HisId     uint   `gorm:"primary_key;auto-increment"`
	VideoId   uint   `gorm:"ForeignKey:VideoId"`
	VideoName string `gorm:"ForeignKey:Name"`
	Path      string `gorm:"ForeignKey:Path"`
	UserId    uint   `gorm:"ForeignKey:UserId"`
	Count     uint
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"default:null"`
}

type Comments struct {
	CommentId uint `gorm:"primary_key;auto-increment"`
	Commenter string
	Likes     uint
	Content   string
	Type      string
	Origin    uint
	Count     uint
	UpdatedAt time.Time
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

func QueryVideoById(id uint) (video *Videos, err error) {
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

func QueryHistoriesByUserId(id uint) (historyList []*Histories, err error) {
	if err = Db.Where(&Histories{UserId: id}).Find(&historyList).Error; err != nil {
		return nil, err
	}
	return
}

func DeleteOne(id uint) (err error) {
	err = Db.Where(&Histories{HisId: id}).Delete(&Histories{}).Error
	return
}

func DeleteRange(histList []uint) (err error) {
	err = Db.Where("his_id in (?)", histList).Delete(&Histories{}).Error
	return
}

func CreateCommentVideo(comments *Comments, id uint) (err error) {
	tx := Db.Begin()
	if err := tx.Create(&comments).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"comments", gorm.Expr("comments + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func CreateComment(comments *Comments, id uint) (err error) {
	tx := Db.Begin()
	if err := tx.Create(&comments).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: id}).Update(
		"comments", gorm.Expr("comments + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Comments{}).Where(&Comments{CommentId: comments.Origin}).Update(
		"count", gorm.Expr("count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func LikeAComment(commentId uint) error {
	if err := Db.Model(Comments{}).Where(&Comments{CommentId: commentId}).Update(
		"likes", gorm.Expr("likes + 1")).Error; err != nil {
		return err
	}
	return nil
}

func DisLikeAComment(commentId uint) error {
	if err := Db.Model(Comments{}).Where(&Comments{CommentId: commentId}).Update(
		"likes", gorm.Expr("likes - 1")).Error; err != nil {
		return err
	}
	return nil
}

func DeleteComment(videoId, commentId, origin uint) (err error) {
	tx := Db.Begin()
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: videoId}).Update(
		"comments", gorm.Expr("comments - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(Comments{}).Where(&Comments{CommentId: origin}).Update(
		"count", gorm.Expr("count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := Db.Where(&Comments{CommentId: commentId}).Delete(&Comments{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func DeleteCommentVideo(videoId, commentId, count uint) (err error) {
	tx := Db.Begin()
	if err := tx.Model(Videos{}).Where(&Videos{VideoId: videoId}).Update(
		"comments", gorm.Expr("comments - ?", count+1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := Db.Where(&Comments{CommentId: commentId}).Delete(&Comments{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := Db.Where(&Comments{Origin: commentId}).Delete(&Comments{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
