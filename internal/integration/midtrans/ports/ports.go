package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/dto"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*snap.Response, error)
}
