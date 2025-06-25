package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type balanceResponse struct {
	UserID  uint64 `json:"userId"`
	Balance string `json:"balance"`
}

func (h *Handler) HandleGetBalance(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	balance, err := h.store.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	resp := balanceResponse{
		UserID:  userID,
		Balance: balance.StringFixed(2),
	}
	c.JSON(http.StatusOK, resp)
}
