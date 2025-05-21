package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderItemRepository interface {
	InsertNewOrderItem(ctx context.Context, tx *sqlx.Tx, data *entity.OrderItem) error
	GetOrderItemsByOrderID(ctx context.Context, tx *sqlx.Tx, orderID string) ([]*entity.OrderItem, error)
	GetByOrderIDs(ctx context.Context, orderIDs []uuid.UUID) ([]entity.OrderItem, error)
	FindOrderItemsByOrderID(ctx context.Context, orderID string) ([]*entity.OrderItem, error)
}
