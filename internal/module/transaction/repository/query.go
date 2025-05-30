package repository

const (
	queryInsertNewTransaction = `
		INSERT INTO transactions
		(
			id,
			status,
			phone_number,
			name,
			email,
			is_member,
			total_quantity,
			total_product_amount,
			total_amount,
			v_payment_id,
			v_payment_redirect_url,
			v_transaction_id,
			discount_percentage,
			change_money,
			payment_type,
			total_money,
			table_number,
			total_product_capital_price,
			tax_amount,
			notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	queryUpdateCashField = `
		UPDATE transactions
		SET total_money = $1, change_money = $2
		WHERE id = $3
	`
)
