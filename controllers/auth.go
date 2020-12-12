package controllers

import (
	"RamdomLearning/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type message struct {
	Username string
	Password string
}

func Login(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		var result models.Auths
		if err := models.QueryAuth(temp.Username, temp.Password, &result); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func Register(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var result models.Users
	if err := models.QueryUser(temp.Username, &result); err != nil {
		user := models.Users{Username: temp.Username}
		err = models.CreateUser(&user)
		auth := models.Auths{Username: temp.Username, UserId: user.UserId, Password: temp.Password}
		errAuth := models.CreateAuth(&auth)
		if err != nil || errAuth != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Username already exist"})
	}
}
