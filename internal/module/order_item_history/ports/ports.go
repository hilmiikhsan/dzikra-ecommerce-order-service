package ports

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderItemHistoryRepository interface {
	SumProductCapital(ctx context.Context, startDate, endDate time.Time) (float64, error)
	DuplicateOrderItemHistory(ctx context.Context, tx *sqlx.Tx, orderID string) error
}
