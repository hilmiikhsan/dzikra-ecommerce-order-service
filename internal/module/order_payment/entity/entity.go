package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OrderPayment struct {
	ID                 uuid.UUID       `db:"id"`
	OrderID            uuid.UUID       `db:"order_id"`
	PaymentMethod      string          `db:"payment_method"`
	PaymentStatus      string          `db:"payment_status"`
	PaymentType        string          `db:"payment_type"`
	TransactionID      uuid.UUID       `db:"transaction_id"`
	GrossAmount        int64           `db:"gross_amount"`
	TransactionStatus  string          `db:"transaction_status"`
	PaymentCode        string          `db:"payment_code"`
	SignatureKey       string          `db:"signature_key"`
	MidtransResponse   json.RawMessage `db:"midtrans_response"`
	CallbackResponse   json.RawMessage `db:"callback_response"`
	TransactionRequest json.RawMessage `db:"transaction_request"`
	TransactionTime    time.Time       `db:"transaction_time"`
	ExpiredAt          time.Time       `db:"expired_at"`
	ApplicationID      string          `db:"application_id"`
}
