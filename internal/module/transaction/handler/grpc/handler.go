package grpc

import (
	"context"
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/cmd/proto/transaction"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/transaction/dto"
	"github.com/rs/zerolog/log"
)

func (api *TransactionAPI) CreateTransaction(ctx context.Context, req *transaction.CreateTransactionRequest) (*transaction.CreateTransactionResponse, error) {
	var transactionRequests []dto.TransactionRequest
	for _, transaction := range req.TransactionRequests {
		transactionRequests = append(transactionRequests, dto.TransactionRequest{
			Quantity:  int(transaction.Quantity),
			ProductID: int(transaction.ProductId),
		})
	}

	tableNumber, _ := strconv.Atoi(req.TableNumber)
	totalMoney, _ := strconv.Atoi(req.TotalMoney)

	var transactionItems []dto.TransactionItem
	for _, item := range req.TransactionItems {
		transactionItems = append(transactionItems, dto.TransactionItem{
			ID:                      int(item.Id),
			Quantity:                item.Quantity,
			TotalAmount:             item.TotalAmount,
			ProductName:             item.ProductName,
			ProductPrice:            item.ProductPrice,
			ProductID:               int(item.ProductId),
			TotalAmountCapitalPrice: item.TotalAmountCapitalPrice,
			ProductCapitalPrice:     item.ProductCapitalPrice,
		})
	}

	transactionPayload := &dto.CreateTransactionRequest{
		Name:                     req.Name,
		Status:                   req.Status,
		Email:                    req.Email,
		IsMember:                 req.IsMember,
		PhoneNumber:              req.PhoneNumber,
		TransactionRequest:       transactionRequests,
		CallbackFinish:           req.CallbackFinish,
		TableNumber:              tableNumber,
		Notes:                    req.Notes,
		PaymentType:              req.PaymentType,
		TotalMoney:               totalMoney,
		TotalQuantity:            int(req.TotalQuantity),
		TotalProductAmount:       int(req.TotalProductAmount),
		TotalAmount:              int(req.TotalAmount),
		DiscountPercentage:       int(req.DiscountPercentage),
		ChangeMoney:              int(req.ChangeMoney),
		TaxAmount:                int(req.TaxAmount),
		TotalProductCapitalPrice: int(req.TotalProductCapitalPrice),
		TransactionItems:         transactionItems,
		TransactionID:            req.TransactionId,
	}

	res, err := api.TransactionService.CreateTransaction(ctx, transactionPayload)
	if err != nil {
		log.Err(err).Msg("order::CreateOrder - Failed to create order")
		return &transaction.CreateTransactionResponse{
			Message: "failed to create transaction",
		}, nil
	}

	return &transaction.CreateTransactionResponse{
		VTransactionId:      res.VTransactionID,
		VPaymentId:          res.VPaymentID,
		VPaymentRedirectUrl: res.VPaymentRedirectUrl,
		CreatedAt:           res.CreatedAt,
		TransactionItems:    req.TransactionItems,
		Message:             "success",
	}, nil
}
