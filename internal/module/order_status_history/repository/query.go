package repository

const (
	queryInsertNewOrderStatusHistory = `
		INSERT INTO order_status_histories
		(
			order_id,
			status,
			description,
			changed_by
		) VALUES (?, ?, ?, ?)
	`
)
