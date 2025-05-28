package service

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/consumer/ports"
	orderPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/ports"
	orderPaymentPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
	orderStatusPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/ports"
	"github.com/jmoiron/sqlx"
)

var _ ports.ConsumerService = &consumerService{}

type consumerService struct {
	db                     *sqlx.DB
	orderRepository        orderPorts.OrderRepository
	orderPaymentRepository orderPaymentPorts.OrderPaymentRepository
	orderStatusRepository  orderStatusPorts.OrderStatusHistoryRepository
	externalNotification   externalNotification.ExternalNotification
}

func NewConsumerService(
	db *sqlx.DB,
	orderRepository orderPorts.OrderRepository,
	orderPaymentRepository orderPaymentPorts.OrderPaymentRepository,
	orderStatusRepository orderStatusPorts.OrderStatusHistoryRepository,
	externalNotification externalNotification.ExternalNotification,
) *consumerService {
	return &consumerService{
		db:                     db,
		orderRepository:        orderRepository,
		orderPaymentRepository: orderPaymentRepository,
		orderStatusRepository:  orderStatusRepository,
		externalNotification:   externalNotification,
	}
}
