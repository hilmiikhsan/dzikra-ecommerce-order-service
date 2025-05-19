package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID                   int       `db:"id"`
	OrderID              uuid.UUID `db:"order_id"`
	ProductID            int       `db:"product_id"`
	ProductName          string    `db:"product_name"`
	ProductVariant       string    `db:"product_variant"`
	ProductDiscount      *int      `db:"product_discount"`
	Quantity             int       `db:"quantity"`
	ProductWeight        float64   `db:"product_weight"`
	ProductPrice         int       `db:"product_price"`
	ProductDiscountPrice int       `db:"product_discount_price"`
	TotalAmount          int       `db:"total_amount"`
	ProductVariantID     int       `db:"product_variant_id"`
}
