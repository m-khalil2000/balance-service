package handlers

import (
	"balance-service/internal/logger"
	"balance-service/pkg/models"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleGetBalance(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		logger.Log.Warn(
			"failed to parse userId - not a valid signed integer",
			zap.String("userIdParam", userIDStr),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId must be a positive integer"})
		return
	}

	balance, err := h.store.GetBalance(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Warn(
			"failed to retrieve user balance",
			zap.Uint64("userId", userID),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	logger.Log.Info(
		"retrieved user balance",
		zap.Uint64("userId", userID),
		zap.String("balance", balance.StringFixed(2)),
	)

	resp := models.BalanceResponse{
		UserID:  userID,
		Balance: balance.StringFixed(2),
	}
	c.JSON(http.StatusOK, resp)
}
