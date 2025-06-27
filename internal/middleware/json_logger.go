package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger logs HTTP requests to Zap.
func ZapRequestLogger(log *zap.Logger, threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Float64("duration_ms", float64(duration.Milliseconds())),
		}

		// Fallback for unregistered routes (e.g. 404)
		if c.FullPath() == "" {
			fields = append(fields,
				zap.String("unregistered_path", c.Request.URL.Path),
			)
		}

		if duration > threshold {
			fields = append(fields, zap.Bool("slow", true))
			log.Warn("slow HTTP request", fields...)
		} else {
			log.Info("HTTP request", fields...)
		}
	}
}
