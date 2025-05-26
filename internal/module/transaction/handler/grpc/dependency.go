package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/cmd/proto/transaction"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/adapter"
	midtransService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/service"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/ports"
	transactionRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/service"
	transactionItemRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction_item/repository"
)

type TransactionAPI struct {
	TransactionService ports.TransactionService
	transaction.UnimplementedTransactionServiceServer
}

func NewTransactionAPI() *TransactionAPI {
	var handler = new(TransactionAPI)

	// integration service
	midtransService := midtransService.NewMidtransService(adapter.Adapters.DzikraMidtrans)

	// repository
	transactionRepository := transactionRepository.NewTransactionRepository(adapter.Adapters.DzikraPostgres)
	transactionItemRepository := transactionItemRepository.NewTransactionItemRepository(adapter.Adapters.DzikraPostgres)

	// service
	transactionService := service.NewTransactionService(
		adapter.Adapters.DzikraPostgres,
		transactionRepository,
		transactionItemRepository,
		midtransService,
	)

	// handler
	handler.TransactionService = transactionService

	return handler
}
