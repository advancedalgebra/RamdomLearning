package controllers

import (
	"RamdomLearning/models"
	"RamdomLearning/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type detail struct {
	Path      string
	Username  string
	VideoName string
	VideoId   uint
	Count     uint
	UserId    uint
	NewName   string
}

type information struct {
	Path      string
	Username  string
	VideoName string
	Category  string
	Dad       string
	Level     uint
}

func LaunchVideo(c *gin.Context) {
	var temp information
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			video := models.Videos{Owner: temp.Username, Path: temp.Path, Name: temp.VideoName}
			category := models.Categories{VideoId: video.VideoId, Category: temp.Category, Path: video.Path,
				Dad: temp.Dad, Level: temp.Level}
			if err := models.LaunchTransaction(temp.Username, &video, &category); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success", "content": video})
			}
		}
	}
}

func LikeVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.LikeTransaction(temp.VideoId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func ForwardVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UpdateForward(temp.VideoId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

//func ViewVideo(c *gin.Context) {
//	var temp detail
//	if err := c.ShouldBind(&temp); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//	} else {
//		if utils.CheckToken(c, temp.Username) == "" {
//			if err := models.UpdateView(temp.VideoId); err != nil {
//				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//			} else {
//				c.JSON(http.StatusOK, gin.H{"message": "success"})
//			}
//		}
//	}
//}

// 更新不存在的视频不会带来安全性问题
func ViewVideo(c *gin.Context) {
	if id, err := strconv.Atoi(c.Query("video_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if err := models.UpdateView(uint(id)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "Nothing at all!"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}

func ViewVideoToken(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if history, err := models.QueryHistory(temp.UserId, temp.VideoId); err != nil {
				//println("New")
				if video, err := models.QueryVideoById(temp.VideoId); err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "error": err.Error()})
				} else {
					history := models.Histories{VideoId: temp.VideoId, UserId: temp.UserId,
						Path: video.Path, VideoName: video.Name, Count: 1}
					if err := models.UpdateNewViewToken(temp.VideoId, &history); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
					} else {
						c.JSON(http.StatusOK, gin.H{"message": "success"})
					}
				}
			} else if history.DeletedAt == nil {
				//println("no_del")
				//println(history.DeletedAt)
				if err := models.UpdateOldViewToken(temp.VideoId, temp.UserId); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				}
			} else {
				//println("del")
				history := models.Histories{VideoId: history.VideoId, UserId: history.UserId, Path: history.Path,
					VideoName: history.VideoName, Count: history.Count + 1}
				if err := models.UpdateDelViewToken(&history); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				}
			}
		}
	}
}

// 之后要加入不能删别人视频的检查（video的owner和username匹配）
func UnLaunchVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UnLaunchTransaction(temp.Username, temp.VideoId, temp.Count); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func DisLikeVideo(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.DisLikeTransaction(temp.VideoId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func FindVideosByOwner(c *gin.Context) {
	if result, err := models.QueryVideosByOwner(c.Query("username")); err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "Nothing at all!"})
	} else {
		var PathSLice []string
		for _, v := range result {
			PathSLice = append(PathSLice, v.Path)
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "content": PathSLice})
	}
}

func FindByCategory(c *gin.Context) {
	if result, err := models.QueryByCategory(c.Query("category")); err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "Nothing at all!"})
	} else {
		var CategorySLice []string
		for _, v := range result {
			CategorySLice = append(CategorySLice, v.Path)
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "content": CategorySLice})
	}
}

func FindById(c *gin.Context) {
	//for i := 0; i <= 500000; i++ {
	//	if id, err := strconv.Atoi(c.Query("video_id")); err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	//	} else {
	//		if result, err := models.QueryTagById(uint(id)); err != nil {
	//			c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "Nothing at all!"})
	//		} else {
	//			c.JSON(http.StatusOK, gin.H{"message": result})
	//		}
	//	}
	//}
	if id, err := strconv.Atoi(c.Query("video_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if result, err := models.QueryTagById(uint(id)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "Nothing at all!"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": result})
		}
	}
}

func SetVideoName(c *gin.Context) {
	var temp detail
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if _, err := models.QueryVideo(temp.NewName); err != nil {
				if err := models.UpdateVideoName(temp.NewName, temp.VideoId); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				}
			} else {
				c.JSON(http.StatusForbidden, gin.H{"message": "error", "error": "Username already exist"})
			}
		}
	}
}
