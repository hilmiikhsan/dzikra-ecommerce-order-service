package payment_pos

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	payment "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/payment_pos"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) PaymentCallback(ctx context.Context, req *payment.PaymentCallbackRequest) (*payment.PaymentCallbackResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("POS_GRPC_HOST", config.Envs.POS.PosGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::PaymentCallback - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := payment.NewPaymentCallbackServiceClient(conn)

	resp, err := client.PaymentCallback(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::PaymentCallback - Failed to payment callback")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::PaymentCallback - Response error from POS service")
		return nil, fmt.Errorf("get response error from POS service: %s", resp.Message)
	}

	return resp, nil
}
