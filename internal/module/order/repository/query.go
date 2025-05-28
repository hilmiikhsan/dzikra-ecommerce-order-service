package repository

const (
	queryInsertNewOrder = `
		INSERT INTO orders (
			id,
			user_id,
			status,
			shipping_name,
			shipping_address,
			shipping_phone,
			shipping_number,
			shipping_type,
			total_quantity,
			total_weight,
			voucher_discount,
			address_id,
			cost_name,
			cost_service,
			voucher_id,
			total_product_amount,
			total_shipping_cost,
			total_shipping_amount,
			total_amount,
			notes,
			order_date,
			created_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
		RETURNING
			id,
			user_id,
			status,
			shipping_name,
			shipping_address,
			shipping_phone,
			shipping_number,
			shipping_type,
			total_quantity,
			total_weight,
			voucher_discount,
			address_id,
			cost_name,
			cost_service,
			voucher_id,
			total_product_amount,
			total_shipping_cost,
			total_shipping_amount,
			total_amount,
			notes,
			order_date,
			created_at
	`

	queryFindOrderByID = `
		SELECT
			id,
			user_id,
			status,
			shipping_name,
			shipping_address,
			shipping_phone,
			shipping_number,
			shipping_type,
			total_quantity,
			total_weight,
			voucher_discount,
			address_id,
			cost_name,
			cost_service,
			voucher_id,
			total_product_amount,
			total_shipping_cost,
			total_shipping_amount,
			total_amount,
			notes,
			order_date
		FROM orders
		WHERE id = ? AND deleted_at IS NULL
	`

	queryUpdateShippingNumber = `
		UPDATE orders
		SET
			shipping_number = ?,
			updated_at = NOW()
		WHERE 
			id = ? 
			AND deleted_at IS NULL
		RETURNING
			id,
			user_id,
			status,
			shipping_name,
			shipping_address,
			shipping_phone,
			shipping_number,
			shipping_type,
			total_quantity,
			total_weight,
			voucher_discount,
			address_id,
			cost_name,
			cost_service,
			voucher_id,
			total_product_amount,
			total_shipping_cost,
			total_shipping_amount,
			total_amount,
			notes,
			order_date
	`

	queryUpdateOrderTransactionStatus = `
		UPDATE orders
		SET
			status = ?,
			updated_at = NOW()
		WHERE
			id = ?
			AND deleted_at IS NULL
	`

	queryCalculateTotalSummary = `
		SELECT
			COALESCE(SUM(total_shipping_amount),0)  AS sum_shipping,
			COALESCE(SUM(total_amount),0)           AS sum_amount,
			COUNT(*)                                AS count_tx,
			COALESCE(SUM(total_quantity),0)         AS sum_quantity
		FROM order_histories
		WHERE order_date BETWEEN ? AND ?
	`

	queryUpdateOrderStatus = `
		UPDATE orders
		SET 
			status = ?, 
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
)
