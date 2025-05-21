package ports

import (
	"context"
	"time"
)

type OrderItemHistoryRepository interface {
	SumProductCapital(ctx context.Context, startDate, endDate time.Time) (float64, error)
}
