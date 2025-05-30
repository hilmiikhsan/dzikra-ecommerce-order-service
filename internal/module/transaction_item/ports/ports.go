package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction_item/entity"
	"github.com/jmoiron/sqlx"
)

type TransactionItemRepository interface {
	InsertNewTransactionItem(ctx context.Context, tx *sqlx.Tx, data *entity.TransactionItem) error
}
