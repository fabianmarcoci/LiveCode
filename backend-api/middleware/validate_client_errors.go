package middleware

import (
	"livecode-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateClientErrorLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.ClientErrorLog

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request format",
			})
			c.Abort()
			return
		}

		if len(payload.ErrorMessage) > 1000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error message too long (max 1000 characters)",
			})
			c.Abort()
			return
		}

		c.Set("validated_payload", payload)
		c.Next()
	}
}
