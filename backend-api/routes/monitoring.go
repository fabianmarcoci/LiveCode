package routes

import (
	"livecode-api/middleware"
	"livecode-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogClientError(c *gin.Context) {
	validatedPayload, exists := c.Get("validated_payload")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	payload := validatedPayload.(models.ClientErrorLog)

	middleware.GetLogger(c).Error("client_critical_error",
		zap.String("error_type", payload.ErrorType),
		zap.String("error_message", payload.ErrorMessage),
		zap.String("app_version", payload.AppVersion),
		zap.String("os", payload.OS),
		zap.String("timestamp_client", payload.Timestamp),
	)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
