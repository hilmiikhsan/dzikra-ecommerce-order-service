package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	InsertNewOrder(ctx context.Context, tx *sqlx.Tx, data *entity.Order) (*entity.Order, error)
	CountByFilter(ctx context.Context, userID uuid.UUID, search, status string) (int, error)
	FindByFilter(ctx context.Context, userID uuid.UUID, offset, limit int, search, status string) ([]entity.Order, error)
	FindOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*dto.GetListOrderResponse, error)
	GetOrderById(ctx context.Context, id string) (*dto.OrderDetail, error)
}
