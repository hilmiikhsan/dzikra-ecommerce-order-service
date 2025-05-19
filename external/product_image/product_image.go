package product_image

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/product_image"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) GetImagesByProductIds(ctx context.Context, req *product_image.GetImagesRequest) (*product_image.GetImagesResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("AUTH_GRPC_HOST", config.Envs.Auth.AuthGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetImagesByProductIds - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := product_image.NewProductImageServiceClient(conn)

	resp, err := client.GetImagesByProductIds(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetImagesByProductIds - Failed to get images by product ids")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::GetImagesByProductIds - Response error from auth")
		return nil, fmt.Errorf("get response error from auth: %s", resp.Message)
	}

	return resp, nil
}
