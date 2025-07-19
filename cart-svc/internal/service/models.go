package service

type CartItem struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}

type Cart struct {
	UserId string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}
