package dto

type ItemDetails struct {
	ID       string `json:"id,omitempty"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
	Name     string `json:"name,omitempty"`
	Brand    string `json:"brand,omitempty"`
	Category string `json:"category,omitempty"`
	URL      string `json:"url,omitempty"`
}

type CreateTransactionRequest struct {
	OrderID         string        `json:"order_id"`
	TotalAmount     int64         `json:"total_amount"`
	ItemDetails     []ItemDetails `json:"item_details"`
	ShippingName    string        `json:"shipping_name"`
	Email           string        `json:"email"`
	ShippingPhone   string        `json:"shipping_phone"`
	ShippingAddress string        `json:"shipping_address"`
}
