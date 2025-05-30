package repository

const (
	queryInsertNewTransactionItem = `
		INSERT INTO transaction_items
		(
			id,
			quantity,
			total_amount,
			product_name,
			product_price,
			transaction_id,
			product_id,
			product_capital_price,
			total_amount_capital_price
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
)
