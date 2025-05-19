package address

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/address"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) GetAddressesByIds(ctx context.Context, req *address.GetAddressesByIdsRequest) (*address.GetAddressesResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("AUTH_GRPC_HOST", config.Envs.Auth.AuthGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetAddressesByIds - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := address.NewAddressServiceClient(conn)

	resp, err := client.GetAddressesByIds(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetAddressesByIds - Failed to get addresses by ids")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::GetAddressesByIds - Response error from auth")
		return nil, fmt.Errorf("get response error from auth: %s", resp.Message)
	}

	return resp, nil
}
