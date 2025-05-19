package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.OrderStatusHistoryRepository = &orderStatusHistoryRepository{}

type orderStatusHistoryRepository struct {
	db *sqlx.DB
}

func NewOrderStatusHistoryRepository(db *sqlx.DB) *orderStatusHistoryRepository {
	return &orderStatusHistoryRepository{
		db: db,
	}
}

func (r *orderStatusHistoryRepository) InsertNewOrderStatusHistory(ctx context.Context, tx *sqlx.Tx, data *entity.OrderStatusHistory) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewOrderStatusHistory),
		data.OrderID,
		data.Status,
		data.Description,
		data.ChangedBy,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewOrderStatusHistory - Failed to insert new order status history")
		return err
	}

	return nil
}
