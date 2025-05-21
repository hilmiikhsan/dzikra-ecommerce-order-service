package repository

import (
	"context"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item_history/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.OrderItemHistoryRepository = &orderItemHistoryRepository{}

type orderItemHistoryRepository struct {
	db *sqlx.DB
}

func NewOrderItemHistoryRepository(db *sqlx.DB) *orderItemHistoryRepository {
	return &orderItemHistoryRepository{
		db: db,
	}
}

func (r *orderItemHistoryRepository) SumProductCapital(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	var sum float64

	if err := r.db.GetContext(ctx, &sum, r.db.Rebind(querySumProductCapital), startDate, endDate); err != nil {
		log.Error().Err(err).Msg("repository::SumProductCapital failed")
		return 0, err
	}

	return sum, nil
}
