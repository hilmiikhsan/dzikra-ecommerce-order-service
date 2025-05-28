package notification

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/notification"
)

type ExternalNotification interface {
	SendFcmNotification(ctx context.Context, req *notification.SendFcmNotificationRequest) (*notification.SendFcmNotificationResponse, error)
}
