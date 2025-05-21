package product

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/product"
)

type ExternalProduct interface {
	GetProductStock(ctx context.Context, req *product.GetProductStockRequest) (*product.GetProductStockResponse, error)
}
