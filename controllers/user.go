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
	Follower string
	UserId   uint
}

func Login(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		// 检查用户名密码是否正确
		if _, err := models.QueryAuth(temp.Username, "password", temp.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "error": err.Error()})
		} else {
			h := md5.New()
			_, _ = io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
			token := fmt.Sprintf("%x", h.Sum(nil))
			if err := models.UpdateAuth("token", token, temp.Username); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
			}
		}
	}
}

func Register(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		// 检查用户名是否重复
		if _, err := models.QueryUserAll(temp.Username); err != nil {
			user := models.Users{Username: temp.Username}
			auth := models.Auths{Username: temp.Username, UserId: user.UserId, Password: temp.Password}
			if err := models.Commit(&user, &auth); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success", "content": user})
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "error", "error": "Username already exist"})
		}
	}
}

func Logout(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UpdateAuth("token", "", temp.Username); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func SetTeacher(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UpdateUser("identity", "teacher", temp.Username); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func SetLocation(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UpdateUser("location", temp.Location, temp.Username); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

//func SetUsername(c *gin.Context) {
//	var temp message
//	if err := c.ShouldBind(&temp); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//	}
//	if utils.CheckToken(c, temp.Username) == "" {
//		if user, err := models.QueryUser(temp.NewName); err != nil {
//			if user, err =models.QueryUser(temp.Username); err != nil {
//				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//			} else {
//				user.Username = temp.NewName
//				auth := models.Auths{Username: temp.NewName, UserId: user.UserId}
//				if err := models.Commit(&user, &auth); err != nil {
//					c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "RollBack! Try again!"})
//				} else {
//					c.JSON(http.StatusOK, gin.H{"message": "success"})
//				}
//			}
//		} else {
//			c.JSON(http.StatusForbidden, gin.H{"message": "error", "error": "Username already exist"})
//		}
//	}
//}

func SetUsername(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			// 检查用户名是否存在
			if _, err := models.QueryUserAll(temp.NewName); err != nil {
				if err := models.UpdateTransaction(temp.Username, temp.NewName); err != nil {
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

func DeleteUser(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		// 检查用户名密码是否正确
		if _, err := models.QueryAuth(temp.Username, "password", temp.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "error": err.Error()})
		} else {
			if err := models.DeleteTransaction(temp.Username, temp.UserId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

//之后要检查不能关注自己的检查（username和follower不能重复）
//func Follow(c *gin.Context) {
//	var temp message
//	if err := c.ShouldBind(&temp); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//	} else {
//		if utils.CheckToken(c, temp.Username) == "" {
//			if _, err := models.QueryUser(temp.Follower); err != nil || temp.Follower == temp.Username {
//				c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "error": "Wrong User"})
//			} else {
//				if err := models.FollowTransaction(temp.Username, temp.Follower); err != nil {
//					c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//				} else {
//					c.JSON(http.StatusOK, gin.H{"message": "success"})
//				}
//			}
//		}
//	}
//}

func Follow(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.FollowTransaction(temp.Username, temp.Follower); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

//func UnFollow(c *gin.Context) {
//	var temp message
//	if err := c.ShouldBind(&temp); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//	} else {
//		if utils.CheckToken(c, temp.Username) == "" {
//			if _, err := models.QueryFollowItem(temp.Username, temp.Follower); err != nil {
//				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//			} else {
//				if err := models.UnFollowTransaction(temp.Username, temp.Follower); err != nil {
//					c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
//				} else {
//					c.JSON(http.StatusOK, gin.H{"message": "success"})
//				}
//			}
//		}
//	}
//}

func UnFollow(c *gin.Context) {
	var temp message
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
	} else {
		if utils.CheckToken(c, temp.Username) == "" {
			if err := models.UnFollowTransaction(temp.Username, temp.Follower); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
		}
	}
}

func FindFollower(c *gin.Context) {
	if result, err := models.QueryFollower(c.Query("username")); err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "Nothing at all!"})
	} else {
		var FolloweeSLice []string
		for _, v := range result {
			FolloweeSLice = append(FolloweeSLice, v.Followee)
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "content": FolloweeSLice})
	}
}

func FindFollowing(c *gin.Context) {
	if result, err := models.QueryFollowing(c.Query("username")); err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "Nothing at all!"})
	} else {
		var FollowingSLice []string
		for _, v := range result {
			FollowingSLice = append(FollowingSLice, v.Follower)
		}
		c.JSON(http.StatusOK, gin.H{"message": "success", "content": FollowingSLice})
	}
}
