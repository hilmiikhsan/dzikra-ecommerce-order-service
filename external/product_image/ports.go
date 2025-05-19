package product_image

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/product_image"
)

type ExternalProductImage interface {
	GetImagesByProductIds(ctx context.Context, req *product_image.GetImagesRequest) (*product_image.GetImagesResponse, error)
}
