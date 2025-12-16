package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := uuid.New().String()
		c.Set("correlation_id", correlationID)

		start := time.Now()

		Logger.Info("incoming_request",
			zap.String("correlation_id", correlationID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		c.Next()

		duration := time.Since(start)

		Logger.Info("request_completed",
			zap.String("correlation_id", correlationID),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int("response_size", c.Writer.Size()),
		)
	}
}

func GetLogger(c *gin.Context) *zap.Logger {
	correlationID, exists := c.Get("correlation_id")
	if !exists {
		return Logger
	}

	return Logger.With(zap.String("correlation_id", correlationID.(string)))
}
