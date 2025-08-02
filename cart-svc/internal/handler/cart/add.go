package cart

import (
	"encoding/json"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/models"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req models.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	cart, err := h.svc.AddItem(r.Context(), req.UserID, req.SKU, req.Quantity)
	if err != nil {
		h.logger.Warn("Failed to add item",
			zap.Error(err),
			zap.String("user_id", req.UserID),
			zap.String("sku", req.SKU),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
