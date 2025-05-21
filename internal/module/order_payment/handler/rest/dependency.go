package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
	orderPaymentRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/repository"
	orderPaymentService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/service"
)

type orderPaymentHandler struct {
	service   ports.OrderPaymentService
	validator adapter.Validator
}

func NewOrderPaymentHandler() *orderPaymentHandler {
	var handler = new(orderPaymentHandler)

	// validator
	validator := adapter.Adapters.Validator

	// repository
	orderPaymentRepository := orderPaymentRepository.NewOrderPaymentRepository(adapter.Adapters.DzikraPostgres)

	// order payment service
	orderPaymentService := orderPaymentService.NewOrderPaymentService(orderPaymentRepository)

	// handler
	handler.service = orderPaymentService
	handler.validator = validator

	return handler
}
