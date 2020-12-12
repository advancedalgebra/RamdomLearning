package controllers

import (
	"RamdomLearning/models"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
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
		if err := models.QueryAuth(temp.Username, "password", temp.Password, &result); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			h := md5.New()
			_, _ = io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
			token := fmt.Sprintf("%x", h.Sum(nil))
			if err := models.UpdateAuth("token", token, temp.Username, &result); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
			}
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

func Logout(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var result models.Auths
	if token := c.GetHeader("token"); token != "" {
		if err := models.QueryAuth(temp.Username, "token", token, &result); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Token"})
		} else {
			if err := models.UpdateAuth("token", "", temp.Username, &result); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Missing"})
	}
}
