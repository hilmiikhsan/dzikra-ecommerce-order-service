package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/cmd/proto/order"
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
