package ports

import (
	"context"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/entity"
	orderItem "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	InsertNewOrder(ctx context.Context, tx *sqlx.Tx, data *entity.Order) (*entity.Order, error)
	CountByFilter(ctx context.Context, userID uuid.UUID, search, status string) (int, error)
	CountListOrderTransactionByFilter(ctx context.Context, search, status string) (int, error)
	FindByFilter(ctx context.Context, userID uuid.UUID, offset, limit int, search, status string) ([]entity.Order, error)
	FindListTransactionOrderByFilter(ctx context.Context, offset, limit int, search, status string) ([]entity.Order, error)
	FindOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	UpdateShippingNumber(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, shippingNumber string) (*entity.Order, error)
	UpdateOrderTransactionStatus(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, status string) error
	CalculateTotalSummary(ctx context.Context, startDate, endDate time.Time) (*entity.OrderHistory, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*dto.GetListOrderResponse, error)
	GetOrderById(ctx context.Context, id string) (*dto.OrderDetail, error)
	GetListOrderTransaction(ctx context.Context, page, limit int, search, status string) (*dto.GetListOrderResponse, error)
	UpdateOrderShippingNumber(ctx context.Context, req *dto.UpdateOrderShippingNumberRequest, id string) (*dto.UpdateOrderShippingNumberResponse, error)
	UpdateOrderStatusTransaction(ctx context.Context, req *dto.UpdateOrderStatusTransactionRequest, id string) error
	GetOrderItemsByOrderID(ctx context.Context, orderID string) ([]*orderItem.OrderItem, error)
	CalculateTotalSummary(ctx context.Context, startDate, endDate string) (*dto.CalculateTotalSummaryResponse, error)
}
