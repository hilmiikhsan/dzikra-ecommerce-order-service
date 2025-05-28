package payment_pos

import (
	"context"

	payment "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/payment_pos"
)

type ExternalPaymentPOS interface {
	PaymentCallback(ctx context.Context, req *payment.PaymentCallbackRequest) (*payment.PaymentCallbackResponse, error)
}
