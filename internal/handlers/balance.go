package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type balanceResponse struct {
	UserID  uint64 `json:"userId"`
	Balance string `json:"balance"`
}

func (h *Handler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid userId")
		return
	}

	balance, err := h.store.GetBalance(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	resp := balanceResponse{
		UserID:  userID,
		Balance: balance.StringFixed(2),
	}
	respondWithJSON(w, http.StatusOK, resp)
}
