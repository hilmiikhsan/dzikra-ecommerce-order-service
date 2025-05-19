package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/entity"
	"github.com/jmoiron/sqlx"
)

type OrderStatusHistoryRepository interface {
	InsertNewOrderStatusHistory(ctx context.Context, tx *sqlx.Tx, data *entity.OrderStatusHistory) error
}
