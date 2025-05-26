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

type CreateTransactionEcommerceRequest struct {
	OrderID         string        `json:"order_id"`
	TotalAmount     int64         `json:"total_amount"`
	ItemDetails     []ItemDetails `json:"item_details"`
	ShippingName    string        `json:"shipping_name"`
	Email           string        `json:"email"`
	ShippingPhone   string        `json:"shipping_phone"`
	ShippingAddress string        `json:"shipping_address"`
}

type CreateTransactionPOSRequest struct {
	OrderID         string                   `json:"order_id"`
	Amount          string                   `json:"amount"`
	ItemDetails     []map[string]interface{} `json:"item_details"`
	CustomerDetails map[string]interface{}   `json:"customer_details"`
	EnabledPayments []string                 `json:"enabled_payments"`
	Expiry          map[string]interface{}   `json:"expiry"`
	CallbackFinish  string                   `json:"callback_finish"`
	ApplicationID   string                   `json:"application_id"`
}

type CreateTransactionPOSResponse struct {
	Transaction   map[string]interface{} `json:"transaction"`
	PaymentRecord map[string]interface{} `json:"paymentRecord"`
}
