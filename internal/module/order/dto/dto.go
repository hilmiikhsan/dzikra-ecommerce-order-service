package dto

import "time"

// type CreateOrderRequest struct {
// 	CostName       string `json:"cost_name" validate:"required"`
// 	CostService    string `json:"cost_service" validate:"required"`
// 	AddressID      string `json:"address_id" validate:"required"`
// 	CallbackFinish string `json:"callback_finish" validate:"required"`
// 	VoucherID      string `json:"voucher_id"`
// 	Notes          string `json:"notes"`
// }

type CreateOrderRequest struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	Email               string     `json:"email"`
	Status              string     `json:"status"`
	ShippingName        string     `json:"shipping_name"`
	ShippingAddress     string     `json:"shipping_address"`
	ShippingPhone       string     `json:"shipping_phone"`
	ShippingNumber      string     `json:"shipping_number"`
	ShippingType        string     `json:"shipping_type"`
	TotalQuantity       int        `json:"total_quantity"`
	TotalWeight         float64    `json:"total_weight"`
	VoucherDiscount     int        `json:"voucher_discount"`
	AddressID           int        `json:"address_id"`
	CostName            string     `json:"cost_name"`
	CostService         string     `json:"cost_service"`
	VoucherID           int        `json:"voucher_id"`
	TotalProductAmount  int64      `json:"total_product_amount"`
	TotalShippingCost   int64      `json:"total_shipping_cost"`
	TotalShippingAmount int64      `json:"total_shipping_amount"`
	TotalAmount         int64      `json:"total_amount"`
	Notes               string     `json:"notes"`
	CartItems           []CartItem `json:"cart_items"`
	OrderDate           time.Time  `json:"order_date"`
	CreatedAt           time.Time  `json:"created_at"`
}

type CreateOrderResponse struct {
	Order               OrderDetail `json:"order"`
	MidtransRedirectUrl string      `json:"midtrans_redirect_url"`
	PaymentID           string      `json:"payment_id"`
}

type OrderDetail struct {
	ID                  string `json:"id"`
	OrderDate           string `json:"order_date"`
	Status              string `json:"status"`
	ShippingName        string `json:"shipping_name"`
	ShippingAddress     string `json:"shipping_address"`
	ShippingPhone       string `json:"shipping_phone"`
	ShippingNumber      string `json:"shipping_number"`
	ShippingType        string `json:"shipping_type"`
	TotalWeight         int    `json:"total_weight"`
	TotalQuantity       int    `json:"total_quantity"`
	TotalShippingCost   string `json:"total_shipping_cost"`
	TotalProductAmount  string `json:"total_product_amount"`
	TotalShippingAmount string `json:"total_shipping_amount"`
	TotalAmount         string `json:"total_amount"`
	VoucherDiscount     int    `json:"voucher_discount"`
	VoucherID           string `json:"voucher_id"`
	CostName            string `json:"cost_name"`
	CostService         string `json:"cost_service"`
	AddressID           int    `json:"address_id"`
	UserID              string `json:"user_id"`
	Notes               string `json:"notes"`
}

type CartItem struct {
	ID                          int              `json:"id"`
	Quantity                    int              `json:"quantity"`
	ProductID                   int              `json:"product_id"`
	ProductVariantID            int              `json:"product_variant_id"`
	ProductName                 string           `json:"product_name"`
	ProductRealPrice            string           `json:"product_real_price"`
	ProductDiscountPrice        string           `json:"product_discount_price"`
	ProductStock                int              `json:"product_stock"`
	ProductWeight               float64          `json:"product_weight"`
	ProductVariantWeight        float64          `json:"product_variant_weight"`
	ProductVariantName          string           `json:"product_variant_name"`
	ProductGrocery              []ProductGrocery `json:"product_grocery"`
	ProductVariantSubName       string           `json:"product_variant_sub_name"`
	ProductVariantRealPrice     string           `json:"product_variant_real_price"`
	ProductVariantDiscountPrice string           `json:"product_variant_discount_price"`
	ProductVariantStock         int              `json:"product_variant_stock"`
	ProductImage                []ProductImage   `json:"product_image"`
}

type ProductGrocery struct {
	ID        int `json:"id"`
	MinBuy    int `json:"min_buy"`
	Discount  int `json:"discount"`
	ProductID int `json:"product_id"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ImageURL  string `json:"image_url"`
	Position  int    `json:"position"`
	ProductID int    `json:"product_id"`
}

type GetListOrderResponse struct {
	Orders      []GetListOrder `json:"orders"`
	TotalPages  int            `json:"total_page"`
	CurrentPage int            `json:"current_page"`
	PageSize    int            `json:"page_size"`
	TotalData   int            `json:"total_data"`
}

type GetListOrder struct {
	ID                  string      `json:"id"`
	OrderDate           string      `json:"order_date"`
	Status              string      `json:"status"`
	TotalQuantity       int         `json:"total_quantity"`
	TotalAmount         int         `json:"total_amount"`
	ShippingNumber      string      `json:"shipping_number"`
	TotalShippingAmount int         `json:"total_shipping_amount"`
	CostName            string      `json:"cost_name"`
	CostService         string      `json:"cost_service"`
	VoucherID           int         `json:"voucher_id"`
	VoucherDiscount     int         `json:"voucher_disc"`
	UserID              string      `json:"user_id"`
	Notes               string      `json:"notes"`
	SubTotal            int         `json:"sub_total"`
	Address             Address     `json:"address"`
	OrderItems          []OrderItem `json:"order_item"`
	Payment             Payment     `json:"redirect_url"`
}

type OrderItem struct {
	ProductID             int            `json:"product_id"`
	ProductName           string         `json:"product_name"`
	ProductVariantSubName string         `json:"product_variant_sub_name"`
	ProductVariant        string         `json:"product_variant"`
	TotalAmount           int            `json:"total_amount"`
	ProductDisc           *int           `json:"product_disc"`
	Quantity              int            `json:"quantity"`
	FixPricePerItem       int            `json:"fix_price_per_item"`
	ProductImages         []ProductImage `json:"product_image"`
}

type Payment struct {
	RedirectURL string `json:"redirect_url"`
	Status      string `json:"status"`
}

type Address struct {
	ID                  int    `json:"id"`
	Province            string `json:"province"`
	City                string `json:"city"`
	District            string `json:"district"`
	SubDistrict         string `json:"subdistrict"`
	PostalCode          string `json:"postal_code"`
	Address             string `json:"address"`
	ReceivedName        string `json:"received_name"`
	UserID              string `json:"user_id"`
	CityVendorID        string `json:"city_vendor_id"`
	ProvinceVendorID    string `json:"province_vendor_id"`
	SubDistrictVendorID string `json:"subdistrict_vendor_id"`
}

type UpdateOrderShippingNumberRequest struct {
	ShippingNumber string `json:"shipping_number" validate:"required,min=3,max=50,xss_safe"`
}

type UpdateOrderShippingNumberResponse struct {
	ID                  string  `json:"id"`
	OrderDate           string  `json:"order_date"`
	Status              string  `json:"status"`
	ShippingName        string  `json:"shipping_name"`
	ShippingAddress     string  `json:"shipping_address"`
	ShippingPhone       string  `json:"shipping_phone"`
	ShippingNumber      string  `json:"shipping_number"`
	ShippingType        string  `json:"shipping_type"`
	TotalWeight         int     `json:"total_weight"`
	TotalQuantity       int     `json:"total_quantity"`
	TotalShippingCost   string  `json:"total_shipping_cost"`
	TotalProductAmount  string  `json:"total_product_amount"`
	TotalShippingAmount string  `json:"total_shipping_amount"`
	TotalAmount         string  `json:"total_amount"`
	VoucherDiscount     int     `json:"voucher_disc"`
	VoucherID           *string `json:"voucher_id"`
	CostName            string  `json:"cost_name"`
	CostService         string  `json:"cost_service"`
	AddressID           int     `json:"address_id"`
	UserID              string  `json:"user_id"`
	Notes               string  `json:"notes"`
}

type UpdateOrderStatusTransactionRequest struct {
	Status string `json:"status" validate:"required,min=3,max=50,xss_safe"`
}

type CalculateTotalSummaryResponse struct {
	TotalAmount         float64 `json:"total_amount"`
	TotalTransaction    float64 `json:"total_transaction"`
	TotalSellingProduct int     `json:"total_selling_product"`
	TotalCapital        float64 `json:"total_capital"`
	Netsales            float64 `json:"net_sales"`
}
