package rest

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/notification"
	externalPaymentCallback "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/payment_pos"
	externalUserFcmTOken "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/user_fcm_token"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/adapter"
	orderRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
	orderPaymentRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/repository"
	orderPaymentService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/service"
	orderStatusHistory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/repository"
)

type orderPaymentHandler struct {
	service   ports.OrderPaymentService
	validator adapter.Validator
}

func NewOrderPaymentHandler() *orderPaymentHandler {
	var handler = new(orderPaymentHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalNotification := &externalNotification.External{}
	externalPaymentCallback := &externalPaymentCallback.External{}
	externalUserFcmTOken := &externalUserFcmTOken.External{}

	// repository
	orderPaymentRepository := orderPaymentRepository.NewOrderPaymentRepository(adapter.Adapters.DzikraPostgres)
	orderRepository := orderRepository.NewOrderRepository(adapter.Adapters.DzikraPostgres)
	orderStatusHistoryRepository := orderStatusHistory.NewOrderStatusHistoryRepository(adapter.Adapters.DzikraPostgres)

	// order payment service
	orderPaymentService := orderPaymentService.NewOrderPaymentService(
		adapter.Adapters.DzikraPostgres,
		orderPaymentRepository,
		orderRepository,
		orderStatusHistoryRepository,
		externalNotification,
		externalPaymentCallback,
		externalUserFcmTOken,
	)

	// handler
	handler.service = orderPaymentService
	handler.validator = validator

	return handler
}
