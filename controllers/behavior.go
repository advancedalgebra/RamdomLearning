package controllers

import (
	"RamdomLearning/models"
	"RamdomLearning/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type info struct {
	UserId   uint
	VideoId  uint
	Username string
}

func FavoriteVideo(c *gin.Context) {
	var temp info
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if path, err := models.QueryPath(temp.VideoId); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Video"})
			} else {
				video := models.Favorites{UserId: temp.UserId, Path: path.Path, VideoId: temp.VideoId}
				if err := models.FavoriteTransaction(temp.VideoId, &video); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusOK, video)
				}
			}
		}
	}
}

func DisFavoriteVideo(c *gin.Context) {
	var temp info
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.DisFavoriteTransaction(temp.VideoId, temp.UserId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func FindFavoritesByUserId(c *gin.Context) {
	if id, err := strconv.Atoi(c.Query("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if result, err := models.QueryFavoritesByUserId(uint(id)); err != nil || len(result) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing at all!"})
		} else {
			var PathSLice []string
			for _, v := range result {
				PathSLice = append(PathSLice, v.Path)
			}
			c.JSON(http.StatusOK, PathSLice)
		}
	}
}
