package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"balance-service/internal/storage"
)

type Handler struct {
	store storage.Storage
}

func NewHandler(s storage.Storage) *Handler {
	return &Handler{store: s}
}

type transactionPayload struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionID string `json:"transactionId"`
}

func (h *Handler) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid userId")
		return
	}

	sourceType := r.Header.Get("Source-Type")
	if sourceType == "" {
		respondWithError(w, http.StatusBadRequest, "missing Source-Type header")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	var payload transactionPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	if payload.TransactionID == "" || (payload.State != "win" && payload.State != "lose") {
		respondWithError(w, http.StatusBadRequest, "invalid state or missing transactionId")
		return
	}

	// ✅ Validate transactionId is UUID
	if _, err := uuid.Parse(payload.TransactionID); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid transactionId: must be a valid UUID")
		return
	}

	amount, err := decimal.NewFromString(payload.Amount)
	if err != nil || amount.LessThan(decimal.Zero) {
		respondWithError(w, http.StatusBadRequest, "invalid amount")
		return
	}

	err = h.store.ProcessTransaction(r.Context(), userID, payload.TransactionID, payload.State, sourceType, amount)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "transaction processed successfully"})
}
