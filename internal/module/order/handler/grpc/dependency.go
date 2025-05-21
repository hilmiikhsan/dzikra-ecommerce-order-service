package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/cmd/proto/order"
	externalAddress "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/address"
	externalCart "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/cart"
	externalProduct "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/product"
	externalProductImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/product_image"
	externalProductVariant "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/product_variant"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/adapter"
	midtransService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/service"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/ports"
	orderRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/service"
	orderItemRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item/repository"
	orderPaymentRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/repository"
	orderStatusHistoryRepository "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/repository"
)

type OrderAPI struct {
	OrderService ports.OrderService
	order.UnimplementedOrderServiceServer
}

func NewOrderAPI() *OrderAPI {
	var handler = new(OrderAPI)

	// external
	externalOrder := &externalCart.External{}
	externalProductImage := &externalProductImage.External{}
	externalAddress := &externalAddress.External{}
	externalProductVariant := &externalProductVariant.External{}
	externalProduct := &externalProduct.External{}

	// integration service
	midtransService := midtransService.NewMidtransService(adapter.Adapters.DzikraMidtrans)

	// repository
	orderRepository := orderRepository.NewOrderRepository(adapter.Adapters.DzikraPostgres)
	orderStatusHistoryRepository := orderStatusHistoryRepository.NewOrderStatusHistoryRepository(adapter.Adapters.DzikraPostgres)
	orderItemRepository := orderItemRepository.NewOrderItemRepository(adapter.Adapters.DzikraPostgres)
	orderPaymentRepository := orderPaymentRepository.NewOrderPaymentRepository(adapter.Adapters.DzikraPostgres)

	// service
	orderService := service.NewOrderService(
		adapter.Adapters.DzikraPostgres,
		orderRepository,
		orderStatusHistoryRepository,
		orderItemRepository,
		externalOrder,
		midtransService,
		orderPaymentRepository,
		externalProductImage,
		externalAddress,
		externalProductVariant,
		externalProduct,
	)

	// handler
	handler.OrderService = orderService

	return handler
}
