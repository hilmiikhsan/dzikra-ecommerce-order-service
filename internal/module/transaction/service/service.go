package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	midtransDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *transactionService) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	startTime := time.Now().In(loc).Format("2006-01-02 15:04:05 -07:00")
	expiryUnit := "MINUTES"
	expiryDuration := 60

	midtransReq := &midtransDto.CreateTransactionPOSRequest{
		OrderID: req.TransactionID,
		Amount:  fmt.Sprintf("%d", req.TotalAmount),
		CustomerDetails: map[string]interface{}{
			"first_name": req.Name,
			"email":      req.Email,
			"phone":      req.PhoneNumber,
		},
		ItemDetails:     make([]map[string]interface{}, len(req.TransactionItems)),
		EnabledPayments: []string{"other_qris"},
		ApplicationID:   "POS",
		Expiry: map[string]interface{}{
			"start_time": startTime,
			"unit":       expiryUnit,
			"duration":   expiryDuration,
		},
		CallbackFinish: req.CallbackFinish,
	}

	for i, it := range req.TransactionItems {
		midtransReq.ItemDetails[i] = map[string]interface{}{
			"id":       fmt.Sprintf("%d", it.ProductID),
			"price":    it.ProductPrice,
			"quantity": it.Quantity,
			"name":     it.ProductName,
		}
	}

	midtransReq.ItemDetails = append(midtransReq.ItemDetails, map[string]interface{}{
		"id":       "tax",
		"price":    fmt.Sprintf("%d", req.TaxAmount),
		"quantity": 1,
		"name":     fmt.Sprintf("Pajak (%d%%)", req.TaxAmount),
	})

	if req.DiscountPercentage > 0 {
		discountBI := big.NewInt(0).Mul(big.NewInt(int64(req.TotalProductAmount)), big.NewInt(int64(req.DiscountPercentage)))
		discountBI.Div(discountBI, big.NewInt(100))

		midtransReq.ItemDetails = append(midtransReq.ItemDetails, map[string]interface{}{
			"id":       "discount",
			"price":    fmt.Sprintf("%d", discountBI.Neg(discountBI).Int64()),
			"quantity": 1,
			"name":     "Member Discount",
		})
	}

	midtransResp, err := s.midtransService.CreateTransactionPOS(ctx, midtransReq)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - failed create transaction POS")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	redirectURL := midtransResp.Transaction["redirect_url"].(string)
	paymentID := midtransResp.PaymentRecord["id"].(string)
	transactionID := midtransResp.PaymentRecord["transaction_id"].(string)

	respItems := make([]dto.TransactionItem, len(req.TransactionItems))
	for i, it := range req.TransactionItems {
		respItems[i] = dto.TransactionItem{
			ID:                      it.ProductID,
			Quantity:                it.Quantity,
			TotalAmount:             it.TotalAmount,
			ProductName:             it.ProductName,
			ProductPrice:            it.ProductPrice,
			TransactionID:           req.TransactionID,
			ProductID:               it.ProductID,
			TotalAmountCapitalPrice: it.TotalAmountCapitalPrice,
			ProductCapitalPrice:     it.ProductCapitalPrice,
		}
	}

	return &dto.CreateTransactionResponse{
		VTransactionID:      transactionID,
		VPaymentID:          paymentID,
		VPaymentRedirectUrl: redirectURL,
		TransactionItems:    respItems,
	}, nil
}
