package service

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/notification"
	externalPaymentCallback "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/payment_pos"
	externalUserFcmToken "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/user_fcm_token"
	orderPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/ports"
	orderPaymentPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
	orderStatusHistoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/ports"
	"github.com/jmoiron/sqlx"
)

var _ orderPaymentPorts.OrderPaymentService = &orderPaymentService{}

type orderPaymentService struct {
	db                           *sqlx.DB
	orderPaymentRepository       orderPaymentPorts.OrderPaymentRepository
	orderRepository              orderPorts.OrderRepository
	orderStatusHistoryRepository orderStatusHistoryPorts.OrderStatusHistoryRepository
	externalNotification         externalNotification.ExternalNotification
	externalPaymentCallback      externalPaymentCallback.ExternalPaymentPOS
	externalUserFcmToken         externalUserFcmToken.ExternalUserFcmToken
}

func NewOrderPaymentService(
	db *sqlx.DB,
	orderPaymentRepository orderPaymentPorts.OrderPaymentRepository,
	orderRepository orderPorts.OrderRepository,
	orderStatusHistoryRepository orderStatusHistoryPorts.OrderStatusHistoryRepository,
	externalNotification externalNotification.ExternalNotification,
	externalPaymentCallback externalPaymentCallback.ExternalPaymentPOS,
	externalUserFcmToken externalUserFcmToken.ExternalUserFcmToken,
) *orderPaymentService {
	return &orderPaymentService{
		db:                           db,
		orderPaymentRepository:       orderPaymentRepository,
		orderRepository:              orderRepository,
		orderStatusHistoryRepository: orderStatusHistoryRepository,
		externalNotification:         externalNotification,
		externalPaymentCallback:      externalPaymentCallback,
		externalUserFcmToken:         externalUserFcmToken,
	}
}
