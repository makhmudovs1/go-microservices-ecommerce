package models

// AddItemRequest — тело POST /cart
type AddItemRequest struct {
	UserID   string `json:"user_id"`  // user's UUID
	SKU      string `json:"sku"`      //  article of product
	Quantity int    `json:"quantity"` // number of pieces (>0)
}

// RemoveItemRequest - DELETE /cart
type RemoveItemRequest struct {
	UserID string `json:"user_id"`
	SKU    string `json:"sku"`
}

// CartResponse — ответ GET /cart
type CartItem struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}

// CartItemInfo — элемент в ответе
type Cart struct {
	UserId string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}
