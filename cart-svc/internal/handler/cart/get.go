package cart

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.logger.Warn("No user_id in request")
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	cartData, err := h.svc.GetCartByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Warn("Failed to get cart",
			zap.Error(err),
			zap.String("user_id", userID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartData)
}
