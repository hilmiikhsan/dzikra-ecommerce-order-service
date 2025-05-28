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

	queryGetLatestByOrderID = `
		SELECT
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
        FROM order_payments
        WHERE order_id = $1
        ORDER BY transaction_time DESC
        LIMIT 1
	`

	queryUpdateOrderPayment = `
		UPDATE order_payments SET
			transaction_id     = ?,
			transaction_status = ?,
			payment_type       = ?,
			signature_key      = ?,
			transaction_time   = ?,
			payment_status     = ?,
			updated_at         = NOW()
		WHERE id = ?
	`

	queryUpdateOrderPaymentStatus = `
		UPDATE order_payments SET
			transaction_status = ?,
			payment_status = ?,
			updated_at     = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
)
