package cart

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/cart"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) DeleteCartByUserID(ctx context.Context, req *cart.DeleteCartByUserIdRequest) (*cart.DeleteCartByUserIdResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("AUTH_GRPC_HOST", config.Envs.Auth.AuthGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::DeleteCartByUserID - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := cart.NewCartServiceClient(conn)

	resp, err := client.DeleteCartByUserID(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::DeleteCartByUserID - Failed to delete cart by user id")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::DeleteCartByUserID - Response error from auth")
		return nil, fmt.Errorf("get response error from auth: %s", resp.Message)
	}

	return resp, nil
}
