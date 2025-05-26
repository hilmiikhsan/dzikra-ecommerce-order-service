package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/entity"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	InsertNewTransaction(ctx context.Context, tx *sqlx.Tx, data *entity.Transaction) error
	UpdateCashField(ctx context.Context, tx *sqlx.Tx, id, totalMoney, changeMoney string) error
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
}
