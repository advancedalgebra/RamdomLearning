package controllers

import (
	"RamdomLearning/dao"
	"RamdomLearning/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	c.JSON(200, gin.H{"error": "nyanb"})
}

func Register(c *gin.Context) {
	var auth models.Auth
	c.BindJSON(&auth)
	err := dao.Db.Create(&auth).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, auth)
	}
}
