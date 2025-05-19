package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/address"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/product_image"
	midtransDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/entity"
	orderItemHistory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_item/entity"
	orderPayment "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/entity"
	orderStatusHistory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *orderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateOrder - Failed to rollback transaction")
			}
		}
	}()

	id, _ := uuid.Parse(req.ID)
	userID, _ := uuid.Parse(req.UserID)

	// Insert new order
	orderResult, err := s.orderRepository.InsertNewOrder(ctx, tx, &entity.Order{
		ID:                  id,
		UserID:              userID,
		Status:              req.Status,
		ShippingName:        req.ShippingName,
		ShippingAddress:     req.ShippingAddress,
		ShippingPhone:       req.ShippingPhone,
		ShippingNumber:      req.ShippingNumber,
		ShippingType:        strings.ToUpper(req.ShippingType),
		TotalQuantity:       req.TotalQuantity,
		TotalWeight:         req.TotalWeight,
		VoucherDiscount:     req.VoucherDiscount,
		AddressID:           req.AddressID,
		CostName:            req.CostName,
		CostService:         req.CostService,
		VoucherID:           req.VoucherID,
		TotalProductAmount:  req.TotalProductAmount,
		TotalShippingCost:   req.TotalShippingCost,
		TotalShippingAmount: req.TotalShippingAmount,
		TotalAmount:         req.TotalAmount,
		Notes:               req.Notes,
		OrderDate:           req.OrderDate,
		CreatedAt:           req.CreatedAt,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to insert new order")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	statusList := []struct{ Code, Desc string }{
		{constants.StatusOrderCreated, "Pesanan telah dibuat."},
		{constants.OrderStatusUnpaid, "Pesanan menunggu pembayaran."},
	}

	for _, status := range statusList {
		err = s.orderStatusHistoryRepository.InsertNewOrderStatusHistory(ctx, tx, &orderStatusHistory.OrderStatusHistory{
			OrderID:     orderResult.ID,
			Status:      status.Code,
			Description: status.Desc,
			ChangedBy:   orderResult.UserID,
		})
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to insert new order status history")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	var itemDetails []midtransDto.ItemDetails
	for _, ci := range req.CartItems {
		qty := float64(ci.Quantity)
		var unitPrice, discPrice, weight float64
		if ci.ProductVariantID != 0 {
			vDisc, _ := strconv.ParseFloat(ci.ProductVariantDiscountPrice, 64)
			vReal, _ := strconv.ParseFloat(ci.ProductVariantRealPrice, 64)
			if vDisc > 0 {
				unitPrice = vDisc
			} else {
				unitPrice = vReal
			}
			discPrice = vDisc
			weight = ci.ProductVariantWeight * qty
		} else {
			pDisc, _ := strconv.ParseFloat(ci.ProductDiscountPrice, 64)
			pReal, _ := strconv.ParseFloat(ci.ProductRealPrice, 64)
			if pDisc > 0 {
				unitPrice = pDisc
			} else {
				unitPrice = pReal
			}
			discPrice = pDisc
			weight = ci.ProductWeight * qty
		}
		totalAmt := unitPrice * qty

		if err = s.orderItemRepository.InsertNewOrderItem(ctx, tx, &orderItemHistory.OrderItem{
			OrderID:              orderResult.ID,
			ProductID:            ci.ProductID,
			ProductName:          ci.ProductName,
			ProductVariant:       ci.ProductVariantSubName,
			ProductDiscount:      nil,
			Quantity:             ci.Quantity,
			ProductWeight:        weight,
			ProductPrice:         int(unitPrice),
			ProductDiscountPrice: int(discPrice),
			TotalAmount:          int(totalAmt),
			ProductVariantID:     ci.ProductVariantID,
		}); err != nil {
			log.Error().Err(err).Msg("service::CreateOrder - insert order_item failed")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		name := ci.ProductName
		if ci.ProductVariantSubName != "" {
			name = fmt.Sprintf("%s - %s", ci.ProductName, ci.ProductVariantSubName)
		}

		itemDetails = append(itemDetails, midtransDto.ItemDetails{
			ID:       fmt.Sprintf("%d", ci.ID),
			Price:    int64(unitPrice),
			Quantity: int64(ci.Quantity),
			Name:     name,
		})
	}

	orderItems, err := s.orderItemRepository.GetOrderItemsByOrderID(ctx, tx, orderResult.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrder - Failed to get order items for midtrans")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	for _, orderItem := range orderItems {
		name := orderItem.ProductName
		if orderItem.ProductVariant != "" {
			name = fmt.Sprintf("%s - %s", orderItem.ProductName, orderItem.ProductVariant)
		}

		itemDetails = append(itemDetails, midtransDto.ItemDetails{
			ID:       fmt.Sprintf("%d", orderItem.ID),
			Price:    int64(orderItem.ProductPrice),
			Quantity: int64(orderItem.Quantity),
			Name:     name,
		})
	}

	itemDetails = append(itemDetails, midtransDto.ItemDetails{
		ID:       "shipping_cost",
		Price:    int64(req.TotalShippingAmount),
		Quantity: 1,
		Name:     "Ongkos Kirim",
	})

	if req.VoucherDiscount > 0 {
		itemDetails = append(itemDetails, midtransDto.ItemDetails{
			ID:       "voucher_discount",
			Price:    -int64(req.VoucherDiscount),
			Quantity: 1,
			Name:     "Voucher Discount",
		})
	}

	// mapping transaction on Midtrans
	midtransReq := &midtransDto.CreateTransactionRequest{
		OrderID:         orderResult.ID.String(),
		TotalAmount:     orderResult.TotalAmount,
		ItemDetails:     itemDetails,
		ShippingName:    orderResult.ShippingName,
		Email:           req.Email,
		ShippingPhone:   req.ShippingPhone,
		ShippingAddress: req.ShippingAddress,
	}

	// create transaction on Midtrans
	midtransResp, err := s.midtransService.CreateTransaction(ctx, midtransReq)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to create transaction on Midtrans")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// generate payment ID
	paymentID, err := utils.GenerateUUIDv7String()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to generate payment ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// generate transaction id
	transactionID, err := utils.GenerateUUIDv7String()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to generate transaction ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get current time and add 24 hours for expiredAt
	transactionTime := utils.FormatTimeJakarta()
	expiredAt := transactionTime.Add(24 * time.Hour)

	// marshal midtrans response
	midtransRespJSON, err := json.Marshal(midtransResp)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrder - Failed to marshal midtrans response")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// marshal midtrans request
	transReqJSON, err := json.Marshal(midtransReq)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrder - Failed to marshal midtrans request")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// insert new order payment
	callbackJSON := json.RawMessage("{}")
	err = s.orderPaymentRepository.InsertNewOrderPayment(ctx, tx, &orderPayment.OrderPayment{
		ID:                 paymentID,
		OrderID:            orderResult.ID,
		PaymentMethod:      constants.PaymentMethodMidtrans,
		PaymentStatus:      constants.PaymentStatusPending,
		TransactionID:      transactionID,
		GrossAmount:        orderResult.TotalAmount,
		TransactionTime:    transactionTime,
		TransactionStatus:  constants.PaymentStatusPending,
		ExpiredAt:          expiredAt,
		MidtransResponse:   midtransRespJSON,
		CallbackResponse:   callbackJSON,
		TransactionRequest: transReqJSON,
		ApplicationID:      constants.AppECommerce,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to insert new order payment")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrderResponse{
		Order: dto.OrderDetail{
			ID:                  orderResult.ID.String(),
			OrderDate:           utils.FormatTime(orderResult.OrderDate),
			Status:              orderResult.Status,
			ShippingName:        orderResult.ShippingName,
			ShippingAddress:     orderResult.ShippingAddress,
			ShippingPhone:       orderResult.ShippingPhone,
			ShippingNumber:      orderResult.ShippingNumber,
			ShippingType:        orderResult.ShippingType,
			TotalWeight:         int(orderResult.TotalWeight),
			TotalQuantity:       int(orderResult.TotalQuantity),
			TotalShippingCost:   fmt.Sprintf("%d", orderResult.TotalShippingCost),
			TotalProductAmount:  fmt.Sprintf("%d", orderResult.TotalProductAmount),
			TotalShippingAmount: fmt.Sprintf("%d", orderResult.TotalShippingAmount),
			TotalAmount:         fmt.Sprintf("%d", orderResult.TotalAmount),
			VoucherDiscount:     int(orderResult.VoucherDiscount),
			VoucherID:           fmt.Sprintf("%d", orderResult.VoucherID),
			CostName:            orderResult.CostName,
			CostService:         orderResult.CostService,
			AddressID:           int(orderResult.AddressID),
			UserID:              orderResult.UserID.String(),
			Notes:               orderResult.Notes,
		},
		MidtransRedirectUrl: midtransResp.RedirectURL,
	}, nil
}

func (s *orderService) GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*dto.GetListOrderResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrder - Failed to parse user ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	total, err := s.orderRepository.CountByFilter(ctx, userUUID, search, status)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrder - Failed to count order")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	currentPage, perPage, offset := utils.Paginate(page, limit)
	totalPages := utils.CalculateTotalPages(total, perPage)

	rows, err := s.orderRepository.FindByFilter(ctx, userUUID, offset, perPage, search, status)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrder - Failed to find order")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	orderIDs := make([]uuid.UUID, 0, len(rows))
	addrIDs := make([]int64, 0, len(rows))
	for _, r := range rows {
		orderIDs = append(orderIDs, r.ID)
		addrIDs = append(addrIDs, int64(r.AddressID))
	}

	items, err := s.orderItemRepository.GetByOrderIDs(ctx, orderIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrder - Failed to get order items")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	itemsByOrder := make(map[uuid.UUID][]orderItemHistory.OrderItem, len(items))
	prodIDsSet := make(map[int64]struct{}, len(items))
	for _, it := range items {
		itemsByOrder[it.OrderID] = append(itemsByOrder[it.OrderID], it)
		prodIDsSet[int64(it.ProductID)] = struct{}{}
	}

	prodIDs := make([]int64, 0, len(prodIDsSet))
	for pid := range prodIDsSet {
		prodIDs = append(prodIDs, pid)
	}

	imgResp, err := s.externalProductImage.GetImagesByProductIds(ctx, &product_image.GetImagesRequest{
		ProductIds: prodIDs,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrder - Failed to get product images")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	imagesByProd := make(map[int64][]dto.ProductImage, len(imgResp.Images))
	for _, img := range imgResp.Images {
		imagesByProd[img.ProductId] = append(imagesByProd[img.ProductId], dto.ProductImage{
			ID:        int(img.Id),
			ImageURL:  img.ImageUrl,
			Position:  int(img.Position),
			ProductID: int(img.ProductId),
		})
	}

	addrResp, err := s.externalAddress.GetAddressesByIds(ctx, &address.GetAddressesByIdsRequest{
		Ids: addrIDs,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrder - Failed to get addresses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	addrsByID := make(map[int32]dto.Address, len(addrResp.Addresses))
	for _, a := range addrResp.Addresses {
		addrsByID[a.Id] = dto.Address{
			ID:                  int(a.Id),
			Province:            a.Province,
			ProvinceVendorID:    a.ProvinceVendorId,
			City:                a.City,
			CityVendorID:        a.CityVendorId,
			SubDistrict:         a.Subdistrict,
			SubDistrictVendorID: a.SubdistrictVendorId,
			PostalCode:          a.PostalCode,
			Address:             a.Address,
			ReceivedName:        a.ReceivedName,
			UserID:              a.UserId,
		}
	}

	out := make([]dto.GetListOrder, 0, len(rows))
	for _, r := range rows {
		oiDtos := make([]dto.OrderItem, 0, len(itemsByOrder[r.ID]))
		for _, it := range itemsByOrder[r.ID] {
			var subName, variantName string
			if it.ProductVariant != "" {
				subName = it.ProductVariant
				variantName = it.ProductVariant
			}

			prodImgs := imagesByProd[int64(it.ProductID)]

			oiDtos = append(oiDtos, dto.OrderItem{
				ProductID:             int(it.ProductID),
				ProductName:           it.ProductName,
				ProductVariantSubName: subName,
				ProductVariant:        variantName,
				TotalAmount:           it.TotalAmount,
				ProductDisc:           it.ProductDiscount,
				Quantity:              it.Quantity,
				FixPricePerItem:       it.ProductPrice,
				ProductImages:         prodImgs,
			})
		}

		addr := addrsByID[int32(r.AddressID)]

		out = append(out, dto.GetListOrder{
			ID:                  r.ID.String(),
			OrderDate:           utils.FormatTime(r.OrderDate),
			Status:              r.Status,
			TotalQuantity:       r.TotalQuantity,
			TotalAmount:         int(r.TotalAmount),
			ShippingNumber:      r.ShippingNumber,
			TotalShippingAmount: int(r.TotalShippingAmount),
			VoucherID:           r.VoucherID,
			VoucherDiscount:     r.VoucherDiscount,
			UserID:              r.UserID.String(),
			Notes:               r.Notes,
			SubTotal:            int(r.TotalProductAmount),
			Address:             addr,
			OrderItems:          oiDtos,
			Payment:             dto.Payment{},
		})
	}

	for i := range out {
		pay, err := s.orderPaymentRepository.GetLatestByOrderID(ctx, out[i].ID)
		if err != nil && err != sql.ErrNoRows {
			log.Error().Err(err).Msg("service::GetListOrder - Failed to get latest order payment")
			return nil, err
		}
		if pay != nil {
			var resp struct {
				RedirectURL string `json:"redirect_url"`
			}

			err = json.Unmarshal(pay.MidtransResponse, &resp)
			if err != nil {
				log.Error().Err(err).Msg("service::GetListOrder - Failed to unmarshal midtrans response")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			out[i].Payment.RedirectURL = resp.RedirectURL
		}
	}

	return &dto.GetListOrderResponse{
		Orders:      out,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}, nil
}

func (s *orderService) GetOrderById(ctx context.Context, id string) (*dto.OrderDetail, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Msg("service::GetOrderById - Failed to parse order ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	order, err := s.orderRepository.FindOrderByID(ctx, orderID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrOrderNotFound) {
			log.Error().Err(err).Msg("service::GetOrderById - Order not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrOrderNotFound))
		}

		log.Error().Err(err).Msg("service::GetOrderById - Failed to get order by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.OrderDetail{
		ID:                  order.ID.String(),
		OrderDate:           utils.FormatTime(order.OrderDate),
		Status:              order.Status,
		ShippingName:        order.ShippingName,
		ShippingAddress:     order.ShippingAddress,
		ShippingPhone:       order.ShippingPhone,
		ShippingNumber:      order.ShippingNumber,
		ShippingType:        order.ShippingType,
		TotalWeight:         int(order.TotalWeight),
		TotalQuantity:       int(order.TotalQuantity),
		TotalShippingCost:   fmt.Sprintf("%d", order.TotalShippingCost),
		TotalProductAmount:  fmt.Sprintf("%d", order.TotalProductAmount),
		TotalShippingAmount: fmt.Sprintf("%d", order.TotalShippingAmount),
		TotalAmount:         fmt.Sprintf("%d", order.TotalAmount),
		VoucherDiscount:     int(order.VoucherDiscount),
		VoucherID:           fmt.Sprintf("%d", order.VoucherID),
		CostName:            order.CostName,
		CostService:         order.CostService,
		AddressID:           int(order.AddressID),
		UserID:              order.UserID.String(),
		Notes:               order.Notes,
	}, nil
}
