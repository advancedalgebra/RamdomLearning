package models

import (
	_ "github.com/jinzhu/gorm"
)

type Auth struct {
	AuthId   uint `gorm:"primary_key;auto-increment"`
	Name     string
	Password string
}

type LO struct {
	LOid      string
	Likes     int
	Announcer string
}
