package dto

type CreateTransactionRequest struct {
	Name                     string               `json:"name"`
	Status                   string               `json:"status"`
	Email                    string               `json:"email"`
	IsMember                 bool                 `json:"is_member"`
	PhoneNumber              string               `json:"phone_number"`
	TransactionRequest       []TransactionRequest `json:"transaction_requests"`
	CallbackFinish           string               `json:"callback_finish"`
	TableNumber              int                  `json:"table_number"`
	Notes                    string               `json:"notes"`
	PaymentType              string               `json:"payment_type"`
	TotalMoney               int                  `json:"total_money"`
	TotalQuantity            int                  `json:"total_quantity"`
	TotalProductAmount       int                  `json:"total_product_amount"`
	TotalAmount              int                  `json:"total_amount"`
	VPaymentID               string               `json:"v_payment_id"`
	VPaymentRedirectUrl      string               `json:"v_payment_redirect_url"`
	VTransactionID           string               `json:"v_transaction_id"`
	DiscountPercentage       int                  `json:"discount_percentage"`
	ChangeMoney              int                  `json:"change_money"`
	TotalProductCapitalPrice int                  `json:"total_product_capital_price"`
	TaxAmount                int                  `json:"tax_amount"`
	TransactionItems         []TransactionItem    `json:"transaction_items"`
	TransactionID            string               `json:"transaction_id"`
}

type CreateTransactionResponse struct {
	VTransactionID      string            `json:"v_transaction_id"`
	VPaymentID          string            `json:"v_payment_id"`
	VPaymentRedirectUrl string            `json:"v_payment_redirect_url"`
	CreatedAt           string            `json:"created_at"`
	TransactionItems    []TransactionItem `json:"transaction_items"`
}

type TransactionRequest struct {
	Quantity  int `json:"quantity"`
	ProductID int `json:"product_id"`
}

type TransactionItem struct {
	ID                      int    `json:"id"`
	Quantity                string `json:"quantity"`
	TotalAmount             string `json:"total_amount"`
	ProductName             string `json:"product_name"`
	ProductPrice            string `json:"product_price"`
	TransactionID           string `json:"transaction_id"`
	ProductID               int    `json:"product_id"`
	TotalAmountCapitalPrice string `json:"total_amount_capital_price"`
	ProductCapitalPrice     string `json:"product_capital_price"`
}
