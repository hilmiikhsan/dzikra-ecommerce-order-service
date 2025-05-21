package grpc

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/cmd/proto/order"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/dto"
)

func (api *OrderAPI) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var cartItems []dto.CartItem
	if len(req.CartItems) > 0 {
		cartItems = make([]dto.CartItem, len(req.CartItems))
		for i, ci := range req.CartItems {
			var pg []dto.ProductGrocery
			for _, g := range ci.ProductGroceries {
				pg = append(pg, dto.ProductGrocery{
					ID:        int(g.Id),
					MinBuy:    int(g.MinBuy),
					Discount:  int(g.Discount),
					ProductID: int(g.ProductId),
				})
			}
			var pi []dto.ProductImage
			for _, img := range ci.ProductImages {
				pi = append(pi, dto.ProductImage{
					ID:        int(img.Id),
					ImageURL:  img.ImageUrl,
					Position:  int(img.Position),
					ProductID: int(img.ProductId),
				})
			}

			cartItems[i] = dto.CartItem{
				ID:                          int(ci.Id),
				Quantity:                    int(ci.Quantity),
				ProductID:                   int(ci.ProductId),
				ProductVariantID:            int(ci.ProductVariantId),
				ProductName:                 ci.ProductName,
				ProductRealPrice:            ci.ProductRealPrice,
				ProductDiscountPrice:        ci.ProductDiscountPrice,
				ProductStock:                int(ci.ProductStock),
				ProductWeight:               ci.ProductWeight,
				ProductVariantWeight:        ci.ProductVariantWeight,
				ProductVariantName:          ci.ProductVariantName,
				ProductGrocery:              pg,
				ProductVariantSubName:       ci.ProductVariantSubName,
				ProductVariantRealPrice:     ci.ProductVariantRealPrice,
				ProductVariantDiscountPrice: ci.ProductVariantDiscountPrice,
				ProductVariantStock:         int(ci.ProductVariantStock),
				ProductImage:                pi,
			}
		}
	}

	res, err := api.OrderService.CreateOrder(ctx, &dto.CreateOrderRequest{
		ID:                  req.Id,
		UserID:              req.UserId,
		Email:               req.Email,
		Status:              req.Status,
		ShippingName:        req.ShippingName,
		ShippingAddress:     req.ShippingAddress,
		ShippingPhone:       req.ShippingPhone,
		ShippingNumber:      "",
		ShippingType:        req.ShippingType,
		TotalQuantity:       int(req.TotalQuantity),
		TotalWeight:         req.TotalWeight,
		VoucherDiscount:     int(req.VoucherDiscount),
		AddressID:           int(req.AddressId),
		CostName:            req.CostName,
		CostService:         req.CostService,
		VoucherID:           int(req.VoucherId),
		TotalProductAmount:  req.TotalProductAmount,
		TotalShippingCost:   req.TotalShippingCost,
		TotalShippingAmount: req.TotalShippingAmount,
		TotalAmount:         req.TotalAmount,
		Notes:               req.Notes,
		CartItems:           cartItems,
		OrderDate:           req.CreatedAt.AsTime(),
		CreatedAt:           req.CreatedAt.AsTime(),
	})
	if err != nil {
		log.Err(err).Msg("order::CreateOrder - Failed to create order")
		return &order.CreateOrderResponse{
			Message:             "failed to create order",
			Order:               nil,
			MidtransRedirectUrl: "",
		}, nil
	}

	return &order.CreateOrderResponse{
		Message: "success",
		Order: &order.OrderDetail{
			Id:                  res.Order.ID,
			OrderDate:           res.Order.OrderDate,
			Status:              res.Order.Status,
			ShippingName:        res.Order.ShippingName,
			ShippingAddress:     res.Order.ShippingAddress,
			ShippingPhone:       res.Order.ShippingPhone,
			ShippingNumber:      res.Order.ShippingNumber,
			ShippingType:        res.Order.ShippingType,
			TotalWeight:         int64(res.Order.TotalWeight),
			TotalQuantity:       int64(res.Order.TotalQuantity),
			TotalShippingCost:   res.Order.TotalShippingCost,
			TotalProductAmount:  res.Order.TotalProductAmount,
			TotalShippingAmount: res.Order.TotalShippingAmount,
			TotalAmount:         res.Order.TotalAmount,
			VoucherDiscount:     int64(res.Order.VoucherDiscount),
			VoucherId:           res.Order.VoucherID,
			CostName:            res.Order.CostName,
			CostService:         res.Order.CostService,
			AddressId:           int64(res.Order.AddressID),
			UserId:              res.Order.UserID,
			Notes:               res.Order.Notes,
		},
		MidtransRedirectUrl: res.MidtransRedirectUrl,
	}, nil
}

func (api *OrderAPI) GetListOrder(ctx context.Context, req *order.GetListOrderRequest) (*order.GetListOrderResponse, error) {
	res, err := api.OrderService.GetListOrder(ctx, int(req.Page), int(req.Limit), req.Search, req.Status, req.UserId)
	if err != nil {
		log.Err(err).Msg("order::GetListOrder - Failed to get list order")
		return &order.GetListOrderResponse{
			Message: "failed to get list order",
			Orders:  nil,
		}, nil
	}

	var orders []*order.GetListOrder
	for _, orderDetail := range res.Orders {
		var items []*order.OrderItem
		for _, it := range orderDetail.OrderItems {
			var imgs []*order.ProductImage
			for _, pi := range it.ProductImages {
				imgs = append(imgs, &order.ProductImage{
					Id:        int64(pi.ID),
					ImageUrl:  pi.ImageURL,
					Position:  int64(pi.Position),
					ProductId: int64(pi.ProductID),
				})
			}

			items = append(items, &order.OrderItem{
				ProductId:             int64(it.ProductID),
				ProductName:           it.ProductName,
				ProductVariantSubName: it.ProductVariantSubName,
				ProductVariant:        it.ProductVariant,
				TotalAmount:           int64(it.TotalAmount),
				Quantity:              int32(it.Quantity),
				FixPricePerItem:       int64(it.FixPricePerItem),
				ProductImages:         imgs,
			})
		}

		orders = append(orders, &order.GetListOrder{
			Id:                  orderDetail.ID,
			OrderDate:           orderDetail.OrderDate,
			Status:              orderDetail.Status,
			TotalQuantity:       int32(orderDetail.TotalQuantity),
			TotalAmount:         int64(orderDetail.TotalAmount),
			ShippingNumber:      orderDetail.ShippingNumber,
			TotalShippingAmount: int64(orderDetail.TotalShippingAmount),
			CostName:            orderDetail.CostName,
			CostService:         orderDetail.CostService,
			VoucherId:           int64(orderDetail.VoucherID),
			VoucherDisc:         int64(orderDetail.VoucherDiscount),
			Notes:               orderDetail.Notes,
			SubTotal:            int64(orderDetail.SubTotal),
			Address: &order.Address{
				Id:                  int64(orderDetail.Address.ID),
				Province:            orderDetail.Address.Province,
				City:                orderDetail.Address.City,
				District:            orderDetail.Address.District,
				Subdistrict:         orderDetail.Address.SubDistrict,
				PostalCode:          orderDetail.Address.PostalCode,
				Address:             orderDetail.Address.Address,
				ReceivedName:        orderDetail.Address.ReceivedName,
				UserId:              orderDetail.Address.UserID,
				CityVendorId:        orderDetail.Address.CityVendorID,
				ProvinceVendorId:    orderDetail.Address.ProvinceVendorID,
				SubdistrictVendorId: orderDetail.Address.SubDistrictVendorID,
			},
			OrderItems: items,
			Payment: &order.Payment{
				RedirectUrl: orderDetail.Payment.RedirectURL,
			},
		})
	}

	return &order.GetListOrderResponse{
		Message:     "success",
		Orders:      orders,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		PageSize:    int64(res.PageSize),
		TotalData:   int64(res.TotalData),
	}, nil
}

func (api *OrderAPI) GetOrderById(ctx context.Context, req *order.GetOrderByIdRequest) (*order.GetOrderByIdResponse, error) {
	res, err := api.OrderService.GetOrderById(ctx, req.Id)
	if err != nil {
		if err.Error() == constants.ErrOrderNotFound {
			return &order.GetOrderByIdResponse{
				Message: constants.ErrOrderNotFound,
				Order:   nil,
			}, nil
		}

		log.Err(err).Msg("order::GetOrderByID - Failed to get order by id")
		return &order.GetOrderByIdResponse{
			Message: "failed to get order by id",
			Order:   nil,
		}, nil
	}

	return &order.GetOrderByIdResponse{
		Order: &order.OrderDetail{
			Id:                  res.ID,
			OrderDate:           res.OrderDate,
			Status:              res.Status,
			ShippingName:        res.ShippingName,
			ShippingAddress:     res.ShippingAddress,
			ShippingPhone:       res.ShippingPhone,
			ShippingNumber:      res.ShippingNumber,
			ShippingType:        res.ShippingType,
			TotalWeight:         int64(res.TotalWeight),
			TotalQuantity:       int64(res.TotalQuantity),
			TotalShippingCost:   res.TotalShippingCost,
			TotalProductAmount:  res.TotalProductAmount,
			TotalShippingAmount: res.TotalShippingAmount,
			TotalAmount:         res.TotalAmount,
			VoucherDiscount:     int64(res.VoucherDiscount),
			VoucherId:           res.VoucherID,
			CostName:            res.CostName,
			CostService:         res.CostService,
			AddressId:           int64(res.AddressID),
			UserId:              res.UserID,
			Notes:               res.Notes,
		},
		Message: "success",
	}, nil
}

func (api *OrderAPI) GetListOrderTransaction(ctx context.Context, req *order.GetListOrderRequest) (*order.GetListOrderResponse, error) {
	res, err := api.OrderService.GetListOrderTransaction(ctx, int(req.Page), int(req.Limit), req.Search, req.Status)
	if err != nil {
		log.Err(err).Msg("order::GetListOrderTransaction - Failed to get list order")
		return &order.GetListOrderResponse{
			Message: "failed to get list order",
			Orders:  nil,
		}, nil
	}

	var orders []*order.GetListOrder
	for _, orderDetail := range res.Orders {
		var items []*order.OrderItem
		for _, it := range orderDetail.OrderItems {
			var imgs []*order.ProductImage
			for _, pi := range it.ProductImages {
				imgs = append(imgs, &order.ProductImage{
					Id:        int64(pi.ID),
					ImageUrl:  pi.ImageURL,
					Position:  int64(pi.Position),
					ProductId: int64(pi.ProductID),
				})
			}

			items = append(items, &order.OrderItem{
				ProductId:             int64(it.ProductID),
				ProductName:           it.ProductName,
				ProductVariantSubName: it.ProductVariantSubName,
				ProductVariant:        it.ProductVariant,
				TotalAmount:           int64(it.TotalAmount),
				Quantity:              int32(it.Quantity),
				FixPricePerItem:       int64(it.FixPricePerItem),
				ProductImages:         imgs,
			})
		}

		orders = append(orders, &order.GetListOrder{
			Id:                  orderDetail.ID,
			OrderDate:           orderDetail.OrderDate,
			Status:              orderDetail.Status,
			TotalQuantity:       int32(orderDetail.TotalQuantity),
			TotalAmount:         int64(orderDetail.TotalAmount),
			ShippingNumber:      orderDetail.ShippingNumber,
			TotalShippingAmount: int64(orderDetail.TotalShippingAmount),
			CostName:            orderDetail.CostName,
			CostService:         orderDetail.CostService,
			VoucherId:           int64(orderDetail.VoucherID),
			VoucherDisc:         int64(orderDetail.VoucherDiscount),
			Notes:               orderDetail.Notes,
			SubTotal:            int64(orderDetail.SubTotal),
			Address: &order.Address{
				Id:                  int64(orderDetail.Address.ID),
				Province:            orderDetail.Address.Province,
				City:                orderDetail.Address.City,
				District:            orderDetail.Address.District,
				Subdistrict:         orderDetail.Address.SubDistrict,
				PostalCode:          orderDetail.Address.PostalCode,
				Address:             orderDetail.Address.Address,
				ReceivedName:        orderDetail.Address.ReceivedName,
				UserId:              orderDetail.Address.UserID,
				CityVendorId:        orderDetail.Address.CityVendorID,
				ProvinceVendorId:    orderDetail.Address.ProvinceVendorID,
				SubdistrictVendorId: orderDetail.Address.SubDistrictVendorID,
			},
			OrderItems: items,
			Payment: &order.Payment{
				RedirectUrl: orderDetail.Payment.RedirectURL,
			},
		})
	}

	return &order.GetListOrderResponse{
		Message:     "success",
		Orders:      orders,
		TotalPages:  int64(res.TotalPages),
		CurrentPage: int64(res.CurrentPage),
		PageSize:    int64(res.PageSize),
		TotalData:   int64(res.TotalData),
	}, nil
}

func (api *OrderAPI) UpdateOrderShippingNumber(ctx context.Context, req *order.UpdateOrderShippingNumberRequest) (*order.UpdateOrderShippingNumberResponse, error) {
	res, err := api.OrderService.UpdateOrderShippingNumber(ctx, &dto.UpdateOrderShippingNumberRequest{
		ShippingNumber: req.ShippingNumber,
	}, req.Id)
	if err != nil {
		if err.Error() == constants.ErrOrderNotFound {
			return &order.UpdateOrderShippingNumberResponse{
				Message: constants.ErrOrderNotFound,
			}, nil
		}

		if strings.Contains(err.Error(), constants.ErrShippingNumberAlreadyExists) {
			return &order.UpdateOrderShippingNumberResponse{
				Message: constants.ErrShippingNumberAlreadyExists,
			}, nil
		}

		log.Err(err).Msg("order::UpdateOrderShippingNumber - Failed to update order shipping number")
		return &order.UpdateOrderShippingNumberResponse{
			Message: "failed to update order shipping number",
		}, nil
	}

	var voucherID string
	if res.VoucherID != nil {
		voucherID = *res.VoucherID
	} else {
		voucherID = ""
	}

	return &order.UpdateOrderShippingNumberResponse{
		Id:                  res.ID,
		OrderDate:           res.OrderDate,
		Status:              res.Status,
		ShippingName:        res.ShippingName,
		ShippingAddress:     res.ShippingAddress,
		ShippingPhone:       res.ShippingPhone,
		ShippingNumber:      res.ShippingNumber,
		ShippingType:        res.ShippingType,
		TotalWeight:         int64(res.TotalWeight),
		TotalQuantity:       int64(res.TotalQuantity),
		TotalShippingCost:   res.TotalShippingCost,
		TotalProductAmount:  res.TotalProductAmount,
		TotalShippingAmount: res.TotalShippingAmount,
		TotalAmount:         res.TotalAmount,
		VoucherDiscount:     int64(res.VoucherDiscount),
		VoucherId:           voucherID,
		CostName:            res.CostName,
		CostService:         res.CostService,
		AddressId:           int64(res.AddressID),
		UserId:              res.UserID,
		Notes:               res.Notes,
		Message:             "success",
	}, nil
}

func (api *OrderAPI) UpdateOrderStatusTransaction(ctx context.Context, req *order.UpdateOrderStatusTransactionRequest) (*order.UpdateOrderStatusTransactionResponse, error) {
	err := api.OrderService.UpdateOrderStatusTransaction(ctx, &dto.UpdateOrderStatusTransactionRequest{
		Status: req.Status,
	}, req.OrderId)
	if err != nil {
		log.Err(err).Msg("order::UpdateOrderStatusTransaction - Failed to update order status transaction")
		return &order.UpdateOrderStatusTransactionResponse{
			Message: "failed to update order status transaction",
		}, nil
	}

	return &order.UpdateOrderStatusTransactionResponse{
		Message: "success",
	}, nil
}

func (api *OrderAPI) GetOrderItemsByOrderID(ctx context.Context, req *order.GetOrderItemsByOrderIDRequest) (*order.GetOrderItemsByOrderIDResponse, error) {
	orderItems, err := api.OrderService.GetOrderItemsByOrderID(ctx, req.OrderId)
	if err != nil {
		log.Error().Err(err).Msg("handler::GetOrderItemsByOrderID - Failed to get order items by order ID")
		return &order.GetOrderItemsByOrderIDResponse{
			Message:    "Failed to get order items by order ID",
			OrderItems: nil,
		}, nil
	}

	response := &order.GetOrderItemsByOrderIDResponse{
		Message:    "success",
		OrderItems: make([]*order.OrderDetailItem, len(orderItems)),
	}

	for i, orderItem := range orderItems {
		var productDiscount int64
		if orderItem.ProductDiscount != nil {
			productDiscount = int64(*orderItem.ProductDiscount)
		} else {
			productDiscount = 0
		}

		response.OrderItems[i] = &order.OrderDetailItem{
			Id:               int64(orderItem.ID),
			ProductId:        int64(orderItem.ProductID),
			OrderId:          orderItem.OrderID.String(),
			ProductName:      orderItem.ProductName,
			ProductVariant:   orderItem.ProductVariant,
			ProductDiscount:  productDiscount,
			Quantity:         int64(orderItem.Quantity),
			ProductVariantId: int64(orderItem.ProductVariantID),
		}
	}

	return response, nil
}

func (api *OrderAPI) CalculateTotalSummary(ctx context.Context, req *order.CalculateTotalSummaryRequest) (*order.CalculateTotalSummaryResponse, error) {
	res, err := api.OrderService.CalculateTotalSummary(ctx, req.StartDate, req.EndDate)
	if err != nil {
		log.Err(err).Msg("order::CalculateTotalSummary - Failed to calculate total summary")
		return &order.CalculateTotalSummaryResponse{
			Message: "failed to calculate total summary",
		}, nil
	}

	return &order.CalculateTotalSummaryResponse{
		Message:             "success",
		TotalAmount:         res.TotalAmount,
		TotalTransaction:    res.TotalTransaction,
		TotalSellingProduct: int64(res.TotalSellingProduct),
		TotalCapital:        res.TotalCapital,
		NetSales:            res.Netsales,
	}, nil
}
