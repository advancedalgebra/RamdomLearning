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
	count     uint
}

func LaunchVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			var CategoryDetail models.Categories
			if err := c.ShouldBind(&CategoryDetail); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				video := models.Videos{Owner: temp.Username, Path: temp.Path, Name: temp.VideoName}
				category := models.Categories{VideoId: video.VideoId, Category: CategoryDetail.Category,
					Path: video.Path, Dad: CategoryDetail.Dad, Level: CategoryDetail.Level}
				if err := models.LaunchTransaction(temp.Username, &video, &category); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusOK, video)
				}
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

func UnLaunchVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UnLaunchTransaction(temp.Username, temp.VideoId, temp.count); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func DisLikeVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.DisLikeTransaction(temp.VideoId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func FindVideosByOwner(c *gin.Context) {
	if result, err := models.QueryVideosByOwner(c.Query("username")); err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing at all!"})
	} else {
		var PathSLice []string
		for _, v := range result {
			PathSLice = append(PathSLice, v.Path)
		}
		c.JSON(http.StatusOK, PathSLice)
	}
}

func FindByCategory(c *gin.Context) {
	if result, err := models.QueryByCategory(c.Query("category")); err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing at all!"})
	} else {
		var CategorySLice []string
		for _, v := range result {
			CategorySLice = append(CategorySLice, v.Path)
		}
		c.JSON(http.StatusOK, CategorySLice)
	}
}
