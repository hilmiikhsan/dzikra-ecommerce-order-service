package product_variant

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/product_variant"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) GetProductVariantStock(ctx context.Context, req *product_variant.GetProductVariantStockRequest) (*product_variant.GetProductVariantStockResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("AUTH_GRPC_HOST", config.Envs.Auth.AuthGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetProductVariantStock - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := product_variant.NewProductVariantServiceClient(conn)

	resp, err := client.GetProductVariantStock(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetProductVariantStock - Failed to get product variant stock")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::GetProductVariantStock - Response error from auth")
		return nil, fmt.Errorf("get response error from auth: %s", resp.Message)
	}

	return resp, nil
}
