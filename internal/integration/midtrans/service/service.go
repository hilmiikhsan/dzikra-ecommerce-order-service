package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/dto"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog/log"
)

func (s *midtransService) CreateTransactionEcommerce(ctx context.Context, req *dto.CreateTransactionEcommerceRequest) (*snap.Response, error) {
	var items []midtrans.ItemDetails
	for _, item := range req.ItemDetails {
		items = append(items, midtrans.ItemDetails{
			ID:       item.ID,
			Price:    item.Price,
			Name:     item.Name,
			Brand:    item.Brand,
			Category: item.Category,
			Qty:      int32(item.Quantity),
		})
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.TotalAmount,
		},
		CreditCard: &snap.CreditCardDetails{Secure: true},
		Items:      &items,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.ShippingName,
			Email: req.Email,
			Phone: req.ShippingPhone,
			BillAddr: &midtrans.CustomerAddress{
				Address:     req.ShippingAddress,
				CountryCode: "IDN",
			},
			ShipAddr: &midtrans.CustomerAddress{
				Address:     req.ShippingAddress,
				CountryCode: "IDN",
			},
		},
	}

	snapResp, err := s.snapClient.CreateTransaction(snapReq)
	if err != nil {
		log.Error().Err(err).Msg("midtrans::CreateTransactionEcommerce - failed to create transaction")
		return nil, fmt.Errorf("midtrans::CreateTransactionEcommerce - failed to create transaction: %w", err)
	}

	log.Info().Msgf("midtrans::CreateTransactionEcommerce - transaction created successfully, orderID: %s", req.OrderID)

	return snapResp, nil
}

func (s *midtransService) CreateTransactionPOS(ctx context.Context, req *dto.CreateTransactionPOSRequest) (*dto.CreateTransactionPOSResponse, error) {
	var items []midtrans.ItemDetails
	for _, item := range req.ItemDetails {
		items = append(items, midtrans.ItemDetails{
			ID:    fmt.Sprintf("%v", item["id"]),
			Name:  fmt.Sprintf("%v", item["name"]),
			Price: parseAmount(item["price"]),
			Qty:   int32(parseAmount(item["quantity"])),
		})
	}

	custMap := req.CustomerDetails
	custDetail := &midtrans.CustomerDetails{
		FName: fmt.Sprintf("%v", custMap["first_name"]),
		Email: fmt.Sprintf("%v", custMap["email"]),
		Phone: fmt.Sprintf("%v", custMap["phone"]),
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: parseAmount(req.Amount),
		},
		Items:           &items,
		CustomerDetail:  custDetail,
		EnabledPayments: convertToSnapPayments(req.EnabledPayments),
		Callbacks: &snap.Callbacks{
			Finish: req.CallbackFinish,
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     fmt.Sprintf("%v", req.Expiry["unit"]),
			Duration: parseAmount(req.Expiry["duration"]),
		},
	}

	snapResp, err := s.snapClient.CreateTransaction(snapReq)
	if err != nil {
		log.Error().Err(err).Msg("midtrans::CreateTransactionPOS - failed to create transaction")
		return nil, fmt.Errorf("midtrans::CreateTransactionPOS - failed to create transaction: %w", err)
	}

	log.Info().Msgf("midtrans::CreateTransactionPOS - transaction created successfully, orderID: %s", req.OrderID)

	return &dto.CreateTransactionPOSResponse{
		Transaction: map[string]interface{}{
			"token":        snapResp.Token,
			"redirect_url": snapResp.RedirectURL,
		},
		PaymentRecord: map[string]interface{}{
			"transaction_id": snapResp.Token,
			"id":             snapResp.Token,
		},
	}, nil
}

func parseAmount(v interface{}) int64 {
	switch val := v.(type) {
	case string:
		n, _ := strconv.ParseInt(val, 10, 64)
		return n
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	default:
		return 0
	}
}

func convertToSnapPayments(payments []string) []snap.SnapPaymentType {
	result := []snap.SnapPaymentType{}
	for _, p := range payments {
		result = append(result, snap.SnapPaymentType(p))
	}
	return result
}
