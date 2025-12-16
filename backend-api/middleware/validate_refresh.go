package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ValidateRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required,max=500"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			GetLogger(c).Warn("refresh_token_invalid_payload",
				zap.String("error", err.Error()),
				zap.String("ip", c.ClientIP()),
			)

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid refresh token format.",
			})
			c.Abort()
			return
		}

		c.Set("validated_payload", req)
		c.Next()
	}
}
