package repository

const (
	queryInsertNewProductItem = `
		INSERT INTO order_items
		(
			order_id,
			product_id,
			product_name,
			product_variant,
			product_discount,
			quantity,
			product_weight,
			product_price,
			product_discount_price,
			total_amount,
			product_variant_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	queryGetOrderItemsByOrderId = `
		SELECT
			id,
			order_id,
			product_name,
			product_variant,
			quantity,
			product_price,
			product_discount_price,
			total_amount,
			product_weight,
			product_id,
			product_variant_id
		FROM order_items
		WHERE order_id = ? AND deleted_at IS NULL
	`
)
