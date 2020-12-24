package controllers

import (
	"RamdomLearning/models"
	"RamdomLearning/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type info struct {
	UserId    uint
	VideoId   uint
	HistoryId uint
	Username  string
	HisList   []uint
}

type mess struct {
	VideoId   uint
	CommentId uint
	Username  string
	Type      string
	Count     uint
}

type comDetail struct {
	Commenter string
	Content   string
	Type      string
	Origin    uint
	VideoId   uint
}

func FavoriteVideo(c *gin.Context) {
	var temp info
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			// 为了获得视频的路径
			if video, err := models.QueryVideoById(temp.VideoId); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				video := models.Favorites{UserId: temp.UserId, Path: video.Path, VideoId: temp.VideoId}
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

func FindHistory(c *gin.Context) {
	var temp info
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if result, err := models.QueryHistoriesByUserId(temp.UserId); err != nil || len(result) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing at all!"})
			} else {
				c.JSON(http.StatusOK, result)
			}
		}
	}
}

func DeleteOneHistory(c *gin.Context) {
	var temp info
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.DeleteOne(temp.HistoryId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing at all!"})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func DeleteRangeHistory(c *gin.Context) {
	var temp info
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			//println(temp.HisList)
			if err := models.DeleteRange(temp.HisList); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing at all!"})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func LaunchComment(c *gin.Context) {
	var temp comDetail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Commenter) == "" {
			comment := models.Comments{Commenter: temp.Commenter, Content: temp.Content,
				Type: temp.Type, Origin: temp.Origin}
			if temp.Type == "comment" {
				comment.Count = temp.VideoId
				if err := models.CreateComment(&comment, temp.VideoId); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusOK, comment)
				}
			} else {
				if err := models.CreateCommentVideo(&comment, temp.VideoId); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusOK, comment)
				}
			}

		}
	}
}

func LikeComment(c *gin.Context) {
	var temp mess
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.LikeAComment(temp.CommentId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func DisLikeComment(c *gin.Context) {
	var temp mess
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.DisLikeAComment(temp.CommentId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func DeleteComment(c *gin.Context) {
	var temp mess
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if temp.Type == "video" {
				if err := models.DeleteCommentVideo(temp.VideoId, temp.CommentId, temp.Count); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				}
			} else {
				if err := models.DeleteComment(temp.VideoId, temp.CommentId, temp.Count); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				}
			}
		}
	}
}
