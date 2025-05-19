package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                  uuid.UUID `db:"id"`
	UserID              uuid.UUID `db:"user_id"`
	Status              string    `db:"status"`
	ShippingName        string    `db:"shipping_name"`
	ShippingAddress     string    `db:"shipping_address"`
	ShippingPhone       string    `db:"shipping_phone"`
	ShippingNumber      string    `db:"shipping_number"`
	ShippingType        string    `db:"shipping_type"`
	TotalQuantity       int       `db:"total_quantity"`
	TotalWeight         float64   `db:"total_weight"`
	VoucherDiscount     int       `db:"voucher_discount"`
	AddressID           int       `db:"address_id"`
	CostName            string    `db:"cost_name"`
	CostService         string    `db:"cost_service"`
	VoucherID           int       `db:"voucher_id"`
	TotalProductAmount  int64     `db:"total_product_amount"`
	TotalShippingCost   int64     `db:"total_shipping_cost"`
	TotalShippingAmount int64     `db:"total_shipping_amount"`
	TotalAmount         int64     `db:"total_amount"`
	Notes               string    `db:"notes"`
	OrderDate           time.Time `db:"order_date"`
	CreatedAt           time.Time `db:"created_at"`
}
