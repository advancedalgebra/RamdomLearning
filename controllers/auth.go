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
	models.DeleteUser("a1")
}

func Register(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var result models.Users
	err := models.QueryUser(temp.Username, &result)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exist"})
	}
}
