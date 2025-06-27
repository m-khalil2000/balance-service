package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"balance-service/internal/storage"
	"balance-service/pkg/models"
)

type Handler struct {
	store storage.Storage
}

func NewHandler(s storage.Storage) *Handler {
	return &Handler{store: s}
}

// HandleTransaction processes a transaction for a user.
func (h *Handler) HandleTransaction(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	sourceType := c.GetHeader("Source-Type")
	if sourceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing Source-Type header"})
		return
	}

	var payload models.TransactionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	if payload.TransactionID == "" || (payload.State != "win" && payload.State != "lose") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state or missing transactionId"})
		return
	}

	if _, err := uuid.Parse(payload.TransactionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transactionId: must be a valid UUID"})
		return
	}

	amount, err := decimal.NewFromString(payload.Amount)
	if err != nil || amount.LessThan(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	oldBalance, newBalance, err := h.store.ProcessTransaction(c.Request.Context(), userID, payload.TransactionID, payload.State, sourceType, amount)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "transaction already processed":
			status = http.StatusConflict
		case "invalid state":
			status = http.StatusUnprocessableEntity
		case "insufficient balance":
			status = http.StatusUnprocessableEntity
		case "user not found":
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	resp := models.TransactionResponse{
		Message:    "transaction processed successfully",
		OldBalance: oldBalance.StringFixed(2),
		NewBalance: newBalance.StringFixed(2),
	}
	c.JSON(http.StatusOK, resp)

}
