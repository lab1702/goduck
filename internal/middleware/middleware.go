package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			"status_code":   param.StatusCode,
			"latency":       param.Latency,
			"client_ip":     param.ClientIP,
			"method":        param.Method,
			"path":          param.Path,
			"user_agent":    param.Request.UserAgent(),
			"error_message": param.ErrorMessage,
		}).Info("HTTP Request")
		return ""
	})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithField("panic", recovered).Error("Panic recovered")
		c.JSON(500, gin.H{
			"error":     "Internal server error",
			"timestamp": time.Now(),
		})
	})
}
