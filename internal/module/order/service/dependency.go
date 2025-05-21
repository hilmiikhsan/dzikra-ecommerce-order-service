package service

import (
	externalAddress "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/address"
	externalCart "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/cart"
	externalProduct "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/product"
	externalProductImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/product_image"
	externalProductVariant "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/product_variant"
	midtransPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/ports"
	orderPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/ports"
	orderItemPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item/ports"
	orderPaymentPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/ports"
	orderStatusHistoryPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/ports"
	"github.com/jmoiron/sqlx"
)

var _ orderPorts.OrderService = &orderService{}

type orderService struct {
	db                           *sqlx.DB
	orderRepository              orderPorts.OrderRepository
	orderStatusHistoryRepository orderStatusHistoryPorts.OrderStatusHistoryRepository
	orderItemRepository          orderItemPorts.OrderItemRepository
	externalCart                 externalCart.ExternalCart
	midtransService              midtransPorts.MidtransService
	orderPaymentRepository       orderPaymentPorts.OrderPaymentRepository
	externalProductImage         externalProductImage.ExternalProductImage
	externalAddress              externalAddress.ExternalAddress
	externalProductVariant       externalProductVariant.ExternalProductVariant
	externalProduct              externalProduct.ExternalProduct
}

func NewOrderService(
	db *sqlx.DB,
	orderRepository orderPorts.OrderRepository,
	orderStatusHistoryRepository orderStatusHistoryPorts.OrderStatusHistoryRepository,
	orderItemRepository orderItemPorts.OrderItemRepository,
	externalCart externalCart.ExternalCart,
	midtransService midtransPorts.MidtransService,
	orderPaymentRepository orderPaymentPorts.OrderPaymentRepository,
	externalProductImage externalProductImage.ExternalProductImage,
	externalAddress externalAddress.ExternalAddress,
	externalProductVariant externalProductVariant.ExternalProductVariant,
	externalProduct externalProduct.ExternalProduct,
) *orderService {
	return &orderService{
		db:                           db,
		orderRepository:              orderRepository,
		orderStatusHistoryRepository: orderStatusHistoryRepository,
		orderItemRepository:          orderItemRepository,
		externalCart:                 externalCart,
		midtransService:              midtransService,
		orderPaymentRepository:       orderPaymentRepository,
		externalProductImage:         externalProductImage,
		externalAddress:              externalAddress,
		externalProductVariant:       externalProductVariant,
		externalProduct:              externalProduct,
	}
}
