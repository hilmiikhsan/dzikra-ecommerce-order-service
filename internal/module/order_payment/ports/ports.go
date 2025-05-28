package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/entity"
	"github.com/jmoiron/sqlx"
)

type OrderPaymentRepository interface {
	InsertNewOrderPayment(ctx context.Context, tx *sqlx.Tx, data *entity.OrderPayment) error
	GetLatestByOrderID(ctx context.Context, orderID string) (*entity.OrderPayment, error)
	UpdateOrderPayment(ctx context.Context, tx *sqlx.Tx, data *entity.OrderPayment) error
	UpdateOrderPaymentStatus(ctx context.Context, tx *sqlx.Tx, data *entity.OrderPayment) error
}

type OrderPaymentService interface {
	HandleMidtransCallback(ctx context.Context, req *dto.MidtransCallbackRequest) error
}
