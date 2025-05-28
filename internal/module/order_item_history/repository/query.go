package repository

const (
	querySumProductCapital = `
		SELECT
			COALESCE(SUM(i.product_capital_price), 0) AS sum_capital
		FROM order_item_histories AS i
		JOIN order_histories AS h
		ON i.order_id = h.id
		WHERE h.order_date BETWEEN ? AND ?
	`

	queryDupOrderHistory = `
		INSERT INTO order_histories
		(
			id, 
			order_date, 
			status, 
			shipping_name, 
			shipping_address, 
			shipping_phone, 
			shipping_number, 
			shipping_type,
			total_weight, 
			total_quantity, 
			total_shipping_cost, 
			total_product_amount, 
			total_shipping_amount,
			total_amount, 
			voucher_discount, 
			voucher_id, 
			cost_name,
			cost_service, 
			address_id, 
			user_id, 
			created_at
		)
		SELECT 
			id, 
			order_date, 
			status, 
			shipping_name, 
			shipping_address, 
			shipping_phone,
			shipping_number, 
			shipping_type,
			total_weight, 
			total_quantity, 
			total_shipping_cost,
			total_product_amount, 
			total_shipping_amount,
			total_amount, 
			voucher_discount,
			voucher_id, 
			cost_name, 
			cost_service, 
			address_id, 
			user_id, 
			created_at
		FROM orders
		WHERE id = $1
	`

	queryDupItemHistory = `
		INSERT INTO order_item_histories
		(
			order_id, 
			product_id, 
			product_name, 
			product_variant, 
			product_price,
			quantity, 
			total_amount, 
			created_at
		)
		SELECT 
			order_id, 
			product_id, 
			product_name, 
			product_variant, 
			product_price,
			quantity, 
			total_amount,
			created_at
		FROM order_items
		WHERE order_id = $1
	`
)
