package cart

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/cart"
)

type ExternalCart interface {
	DeleteCartByUserID(ctx context.Context, req *cart.DeleteCartByUserIdRequest) (*cart.DeleteCartByUserIdResponse, error)
}
