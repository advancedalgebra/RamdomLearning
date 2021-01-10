package utils

import (
	"RamdomLearning/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context, username string) (err string) {
	if token := c.GetHeader("token"); token != "" {
		if _, err := models.QueryAuth(username, "token", token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "error": "Wrong Token"})
			return "wrong"
		} else {
			return ""
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "error": "Token Missing"})
		return "missing"
	}
}
