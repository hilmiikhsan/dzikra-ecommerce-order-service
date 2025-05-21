package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item/ports"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.OrderItemRepository = &orderItemHistoryRepository{}

type orderItemHistoryRepository struct {
	db *sqlx.DB
}

func NewOrderItemRepository(db *sqlx.DB) *orderItemHistoryRepository {
	return &orderItemHistoryRepository{
		db: db,
	}
}

func (r *orderItemHistoryRepository) InsertNewOrderItem(ctx context.Context, tx *sqlx.Tx, data *entity.OrderItem) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewProductItem),
		data.OrderID,
		data.ProductID,
		data.ProductName,
		data.ProductVariant,
		data.ProductDiscount,
		data.Quantity,
		data.ProductWeight,
		data.ProductPrice,
		data.ProductDiscountPrice,
		data.TotalAmount,
		data.ProductVariantID,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewProductItem - Failed to insert new order item")
		return err
	}

	return nil
}

func (r *orderItemHistoryRepository) GetOrderItemsByOrderID(ctx context.Context, tx *sqlx.Tx, orderID string) ([]*entity.OrderItem, error) {
	var rows []*entity.OrderItem

	if err := tx.SelectContext(ctx, &rows, r.db.Rebind(queryGetOrderItemsByOrderId), orderID); err != nil {
		log.Error().Err(err).Msgf("repository::GetOrderItemsByOrderId - failed for order_id=%s", orderID)
		return nil, err
	}

	return rows, nil
}

func (r *orderItemHistoryRepository) GetByOrderIDs(ctx context.Context, orderIDs []uuid.UUID) ([]entity.OrderItem, error) {
	if len(orderIDs) == 0 {
		log.Warn().Msg("repository::GetByOrderIDs - orderIDs is empty")
		return nil, nil
	}

	query, args, err := sqlx.In(`
        SELECT
            id,
            order_id,
            product_id,
            product_name,
            product_variant,
            quantity,
            product_price,
            product_discount_price,
            total_amount
        FROM order_items
        WHERE order_id IN (?)
    `, orderIDs)
	if err != nil {
		log.Error().Err(err).Msgf("repository::GetByOrderIDs - failed to build query for order IDs: %v", orderIDs)
		return nil, err
	}

	query = r.db.Rebind(query)

	var items []entity.OrderItem
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		log.Error().Err(err).Msgf("repository::GetByOrderIDs - failed to get order items by order IDs: %v", orderIDs)
		return nil, err
	}

	return items, nil
}

func (r *orderItemHistoryRepository) FindOrderItemsByOrderID(ctx context.Context, orderID string) ([]*entity.OrderItem, error) {
	var rows []*entity.OrderItem

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(queryGetOrderItemsByOrderId), orderID); err != nil {
		log.Error().Err(err).Msgf("repository::FindOrderItemsByOrderID - failed for order_id=%s", orderID)
		return nil, err
	}

	return rows, nil
}
