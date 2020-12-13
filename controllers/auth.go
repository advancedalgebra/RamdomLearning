package controllers

import (
	"RamdomLearning/models"
	"RamdomLearning/utils"
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
	Location string
	NewName  string
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
			if err := models.UpdateAuth("token", token, temp.Username); err != nil {
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
	if utils.CheckToken(c, temp.Username) == "" {
		if err := models.UpdateAuth("token", "", temp.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}

func SetTeacher(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if utils.CheckToken(c, temp.Username) == "" {
		if err := models.UpdateUser("identity", "teacher", temp.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}

func SetLocation(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if utils.CheckToken(c, temp.Username) == "" {
		if err := models.UpdateUser("location", temp.Location, temp.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}

func SetUsername(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var result models.Users
	if utils.CheckToken(c, temp.Username) == "" {
		if models.QueryUser(temp.NewName, &result) != nil {
			err := models.UpdateUser("username", temp.NewName, temp.Username)
			errAuth := models.UpdateAuth("username", temp.NewName, temp.Username)
			if err != nil || errAuth != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Username already exist"})
		}
	}
}

func DeleteUser(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		var result models.Auths
		if err := models.QueryAuth(temp.Username, "password", temp.Password, &result); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			errUser := models.DeleteUser(temp.Username)
			err = models.DeleteAuth(temp.Username, temp.Password)
			if err != nil || errUser != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}
