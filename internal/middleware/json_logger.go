package middleware

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp string  `json:"timestamp"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	Status    int     `json:"status"`
	Duration  float64 `json:"duration_ms"`
	Slow      bool    `json:"slow,omitempty"`
}

// JSONRequestLogger logs request details in JSON format to a file
func JSONRequestLogger(filePath string, threshold time.Duration) gin.HandlerFunc {
	// Open file in append mode
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		entry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Method:    c.Request.Method,
			Path:      c.FullPath(),
			Status:    c.Writer.Status(),
			Duration:  float64(duration.Milliseconds()),
		}

		if entry.Path == "" {
			entry.Path = c.Request.URL.Path // fallback for 404 or unregistered routes
		}

		if duration > threshold {
			entry.Slow = true
		}

		jsonData, _ := json.Marshal(entry)
		logFile.Write(append(jsonData, '\n')) // write as one line per request
	}
}
