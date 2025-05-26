package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/dto"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreateTransactionEcommerce(ctx context.Context, req *dto.CreateTransactionEcommerceRequest) (*snap.Response, error)
	CreateTransactionPOS(ctx context.Context, req *dto.CreateTransactionPOSRequest) (*dto.CreateTransactionPOSResponse, error)
}
