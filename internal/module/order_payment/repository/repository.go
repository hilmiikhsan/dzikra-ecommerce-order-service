package repository

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.OrderPaymentRepository = &orderPaymentRepository{}

type orderPaymentRepository struct {
	db *sqlx.DB
}

func NewOrderPaymentRepository(db *sqlx.DB) *orderPaymentRepository {
	return &orderPaymentRepository{
		db: db,
	}
}

func (r *orderPaymentRepository) InsertNewOrderPayment(ctx context.Context, tx *sqlx.Tx, data *entity.OrderPayment) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewOrderPayment),
		data.ID,
		data.OrderID,
		data.PaymentMethod,
		data.PaymentStatus,
		data.PaymentType,
		data.TransactionID,
		data.GrossAmount,
		data.TransactionStatus,
		data.PaymentCode,
		data.SignatureKey,
		data.MidtransResponse,
		data.CallbackResponse,
		data.TransactionRequest,
		data.TransactionTime,
		data.ExpiredAt,
		data.ApplicationID,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewOrderPayment - Failed to insert new order payment")
		return err
	}

	return nil
}

func (r *orderPaymentRepository) GetLatestByOrderID(ctx context.Context, orderID string) (*entity.OrderPayment, error) {
	const query = `
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

	var p entity.OrderPayment
	err := r.db.GetContext(ctx, &p, query, orderID)
	if err == sql.ErrNoRows {
		log.Error().Err(err).Msgf("repository::GetLatestByOrderID - no order payment found for order_id=%s", orderID)
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err).Msgf("repository::GetLatestByOrderID - failed for order_id=%s", orderID)
		return nil, err
	}

	return &p, nil
}
