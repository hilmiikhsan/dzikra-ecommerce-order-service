package service

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/dto"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog/log"
)

func (s *midtransService) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*snap.Response, error) {
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
		log.Error().Err(err).Msg("midtrans::CreateTransaction - failed to create transaction")
		return nil, fmt.Errorf("midtrans::CreateTransaction - failed to create transaction: %w", err)
	}

	log.Info().Msgf("midtrans::CreateTransaction - transaction created successfully, orderID: %s", req.OrderID)

	return snapResp, nil
}
