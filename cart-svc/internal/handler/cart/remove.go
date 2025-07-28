package cart

import (
	"encoding/json"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/models"
	"net/http"
)

func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.RemoveItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	cart, err := h.svc.RemoveItem(r.Context(), req.UserID, req.SKU)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)

}
