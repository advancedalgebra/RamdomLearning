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
	auth := router.Group("/auth")
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
