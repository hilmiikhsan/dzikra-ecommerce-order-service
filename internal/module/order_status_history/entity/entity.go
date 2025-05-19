package entity

import (
	"github.com/google/uuid"
)

type OrderStatusHistory struct {
	ID          int       `db:"id"`
	OrderID     uuid.UUID `db:"order_id"`
	Status      string    `db:"status"`
	Description string    `db:"description"`
	ChangedBy   uuid.UUID `db:"changed_by"`
}
