package product_variant

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/product_variant"
)

type ExternalProductVariant interface {
	GetProductVariantStock(ctx context.Context, req *product_variant.GetProductVariantStockRequest) (*product_variant.GetProductVariantStockResponse, error)
}
