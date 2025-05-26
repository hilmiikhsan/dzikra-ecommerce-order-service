package service

import (
	midtransPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/ports"
	transactionPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/ports"
	transactionItemPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction_item/ports"
	"github.com/jmoiron/sqlx"
)

var _ transactionPorts.TransactionService = &transactionService{}

type transactionService struct {
	db                        *sqlx.DB
	transactionRepository     transactionPorts.TransactionRepository
	transactionItemRepository transactionItemPorts.TransactionItemRepository
	midtransService           midtransPorts.MidtransService
}

func NewTransactionService(
	db *sqlx.DB,
	transactionRepository transactionPorts.TransactionRepository,
	transactionItemRepository transactionItemPorts.TransactionItemRepository,
	midtransService midtransPorts.MidtransService,
) *transactionService {
	return &transactionService{
		db:                        db,
		transactionRepository:     transactionRepository,
		transactionItemRepository: transactionItemRepository,
		midtransService:           midtransService,
	}
}
