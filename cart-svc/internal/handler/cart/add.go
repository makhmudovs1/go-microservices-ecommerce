package cart

import (
	"encoding/json"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/models"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/service"
	"net/http"
)

type Handler struct {
	svc service.CartService
}

func NewHandler(svc service.CartService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req models.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		http.Error(w, "quantity must be > 0", http.StatusBadRequest)
		return
	}
	if req.SKU == "" {
		http.Error(w, "sku is required", http.StatusBadRequest)
		return
	}
	cart, err := h.svc.AddItem(r.Context(), req.UserID, req.SKU, req.Quantity)
	if err != nil {
		http.Error(w, "service error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
