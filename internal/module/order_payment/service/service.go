package service

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/dto"
)

func (s *orderPaymentService) HandleMidtransCallback(ctx context.Context, req *dto.MidtransCallbackRequest) error {
	// pay, err := s.orderPaymentRepository.GetLatestByOrderID(ctx, req.OrderID)
	// if err != nil {
	// 	log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to get latest order payment")
	// 	return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	// }

	// if req.TransactionStatus == constants.MidtransSettlement {
	// }

	return nil
}
