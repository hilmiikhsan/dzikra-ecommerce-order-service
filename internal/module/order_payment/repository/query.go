package repository

const (
	queryInsertNewOrderPayment = `
		INSERT INTO order_payments
		(
			id,
			order_id,
			payment_method,
			payment_status,
			payment_type,
			transaction_id,
			gross_amount,
			transaction_status,
			payment_code,
			signature_key,
			midtrans_response,
			callback_response,
			transaction_request,
			transaction_time,
			expired_at,
			application_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
)
