package utils

import (
	"RamdomLearning/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context, username string) (err string) {
	var result models.Auths
	if token := c.GetHeader("token"); token != "" {
		if err := models.QueryAuth(username, "token", token, &result); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Token"})
			return "wrong"
		} else {
			return ""
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Missing"})
		return "missing"
	}
}
