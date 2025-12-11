package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated.",
		})
		return
	}

	username, _ := c.Get("username")
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully.",
		"user": gin.H{
			"id":       userID,
			"username": username,
			"email":    email,
		},
	})
}
