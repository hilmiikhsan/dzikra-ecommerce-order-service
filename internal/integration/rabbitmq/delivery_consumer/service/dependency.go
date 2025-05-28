package service

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/delivery_consumer/ports"
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rajaongkir/ports"
	orderPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/ports"
	orderItemHistoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item_history/ports"
	orderPaymentPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
	orderStatusPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/ports"
	"github.com/jmoiron/sqlx"
)

var _ ports.DeliveryConsumerService = &deliveryConsumerService{}

type deliveryConsumerService struct {
	db                           *sqlx.DB
	orderRepository              orderPorts.OrderRepository
	orderPaymentRepository       orderPaymentPorts.OrderPaymentRepository
	orderStatusHistoryRepository orderStatusPorts.OrderStatusHistoryRepository
	externalNotification         externalNotification.ExternalNotification
	rajaOngkirService            rajaongkirPorts.RajaongkirService
	orderItemHistoryRepository   orderItemHistoryPorts.OrderItemHistoryRepository
}

func NewDeliveryConsumerService(
	db *sqlx.DB,
	orderRepository orderPorts.OrderRepository,
	orderPaymentRepository orderPaymentPorts.OrderPaymentRepository,
	orderStatusHistoryRepository orderStatusPorts.OrderStatusHistoryRepository,
	externalNotification externalNotification.ExternalNotification,
	rajaOngkirService rajaongkirPorts.RajaongkirService,
	orderItemHistoryRepository orderItemHistoryPorts.OrderItemHistoryRepository,
) *deliveryConsumerService {
	return &deliveryConsumerService{
		db:                           db,
		orderRepository:              orderRepository,
		orderPaymentRepository:       orderPaymentRepository,
		orderStatusHistoryRepository: orderStatusHistoryRepository,
		externalNotification:         externalNotification,
		rajaOngkirService:            rajaOngkirService,
		orderItemHistoryRepository:   orderItemHistoryRepository,
	}
}
