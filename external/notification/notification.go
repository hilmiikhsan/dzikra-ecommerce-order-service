package notification

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) SendFcmNotification(ctx context.Context, req *notification.SendFcmNotificationRequest) (*notification.SendFcmNotificationResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("NOTIFICATION_GRPC_HOST", config.Envs.Notification.NotificationGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::SendFcmNotification - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)

	resp, err := client.SendFcmNotification(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::SendFcmNotification - Failed to send FCM notification")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::SendFcmNotification - Response error from notification")
		return nil, fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return resp, nil
}
