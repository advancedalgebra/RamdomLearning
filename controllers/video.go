package controllers

import (
	"RamdomLearning/models"
	"RamdomLearning/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type detail struct {
	Path      string
	Username  string
	VideoName string
	VideoId   uint
}

func LaunchVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			video := models.Videos{Owner: temp.Username, Path: temp.Path, Name: temp.VideoName}
			if err := models.LaunchTransaction(temp.Username, &video); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, video)
			}
		}
	}
}

func LikeVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.LikeTransaction(temp.VideoId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func ForwardVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UpdateForward(temp.VideoId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func ViewVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UpdateView(temp.VideoId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}
