package main

import (
	"RamdomLearning/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	// 解决跨域问题
	router.Use(cors())
	auth := router.Group("/user")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		auth.POST("/logout", controllers.Logout)
		auth.POST("/teacher", controllers.SetTeacher)
		auth.POST("/location", controllers.SetLocation)
		auth.POST("/rename", controllers.SetUsername)
		auth.POST("/delete_user", controllers.DeleteUser)
		auth.POST("/follow", controllers.Follow)
		auth.POST("/unfollow", controllers.UnFollow)
		auth.GET("/find_follower", controllers.FindFollower)
		auth.GET("/find_following", controllers.FindFollowing)
	}
	video := router.Group("/video")
	{
		video.POST("/launch", controllers.LaunchVideo)
		video.POST("/like", controllers.LikeVideo)
		video.POST("/forward", controllers.ForwardVideo)
		video.GET("/view", controllers.ViewVideo)
		video.POST("/view_token", controllers.ViewVideoToken)
		video.POST("/unlaunch", controllers.UnLaunchVideo)
		video.POST("/dislike", controllers.DisLikeVideo)
		video.POST("/rename", controllers.SetVideoName)
		video.GET("/find_video_by_owner", controllers.FindVideosByOwner)
		video.GET("/find_video_by_category", controllers.FindByCategory)
	}
	UserBehavior := router.Group("/behavior")
	{
		UserBehavior.POST("/favorite_video", controllers.FavoriteVideo)
		UserBehavior.POST("/disfavorite_video", controllers.DisFavoriteVideo)
		UserBehavior.GET("/find_favorite", controllers.FindFavoritesByUserId)
	}
	return router
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
