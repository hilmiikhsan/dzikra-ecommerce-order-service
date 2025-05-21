package service

import (
	orderPaymentPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
)

var _ orderPaymentPorts.OrderPaymentService = &orderPaymentService{}

type orderPaymentService struct {
	orderPaymentRepository orderPaymentPorts.OrderPaymentRepository
}

func NewOrderPaymentService(orderPaymentRepository orderPaymentPorts.OrderPaymentRepository) *orderPaymentService {
	return &orderPaymentService{
		orderPaymentRepository: orderPaymentRepository,
	}
}
