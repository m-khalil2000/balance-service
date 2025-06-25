// internal/server/router.go
package server

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"balance-service/internal/handlers"
	"balance-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.New()

	r.Use(middleware.JSONRequestLogger("requests.log", 500*time.Millisecond))
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("PANIC recovered in handler: %v\n%s", recovered, debug.Stack())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "balance-service",
		})
	})

	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().Unix(),
		})
	})

	r.POST("/user/:userId/transaction", h.HandleTransaction)
	r.GET("/user/:userId/balance", h.HandleGetBalance)

	return r
}
