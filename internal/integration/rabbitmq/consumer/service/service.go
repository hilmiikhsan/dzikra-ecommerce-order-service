package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/consumer/ports"
	orderEntity "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/entity"
	orderPaymentEntity "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/entity"
	orderStatusEntity "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/entity"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *consumerService) ExpireOrder(ctx context.Context, req *ports.ExpireOrderRequest) error {
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		log.Error().Err(err).Msg("service::ExpireOrder - Invalid order ID format")
		return fmt.Errorf("invalid orderID: %w", err)
	}

	paymentUUID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		log.Error().Err(err).Msg("service::ExpireOrder - Invalid payment ID format")
		return fmt.Errorf("invalid paymentID: %w", err)
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req.OrderID).Msg("service::CreateOrder - Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req.OrderID).Msg("service::CreateOrder - Failed to rollback transaction")
			}
		}
	}()

	if err = s.orderPaymentRepository.UpdateOrderPaymentStatus(ctx, tx, &orderPaymentEntity.OrderPayment{
		ID:                paymentUUID,
		TransactionStatus: constants.MidtransExpire,
		PaymentStatus:     constants.MidtransExpire,
	}); err != nil {
		log.Error().Err(err).Msg("service::ExpireOrder - Failed to update order payment status")
		return fmt.Errorf("failed to update order payment status: %w", err)
	}

	if err = s.orderRepository.UpdateStatus(ctx, tx, &orderEntity.Order{
		ID:     orderUUID,
		Status: constants.OrderPaymentStatusExpired,
	}); err != nil {
		log.Error().Err(err).Msg("service::ExpireOrder - Failed to update order status")
		return fmt.Errorf("failed to update order status: %w", err)
	}

	orderStatusHistoryPayload := &orderStatusEntity.OrderStatusHistory{
		OrderID:     orderUUID,
		Status:      constants.OrderPaymentStatusExpired,
		Description: "Pesanan telah kedaluwarsa karena pembayaran tidak diselesaikan.",
		ChangedBy:   "Sistem",
		ChangedAt:   time.Now(),
	}

	if err = s.orderStatusRepository.InsertNewOrderStatusHistory(ctx, tx, orderStatusHistoryPayload); err != nil {
		log.Error().Err(err).Msg("service::ExpireOrder - Failed to insert order status history")
		return fmt.Errorf("failed to insert order status history: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req.OrderID).Msg("service::CreateOrder - Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	go func() {
		_, err := s.externalNotification.SendFcmNotification(ctx, &notification.SendFcmNotificationRequest{
			Title:           "Pesanan kedaluwarsa",
			Body:            fmt.Sprintf("Pesanan dengan id %s telah kedaluwarsa karena pembayaran tidak diselesaikan.", req.OrderID),
			UserId:          req.UserID,
			IsStatusChanged: true,
		})
		if err != nil {
			log.Error().Err(err).Msg("ExpireOrder - failed to send FCM")
		}
	}()

	log.Info().Str("orderID", req.OrderID).Str("paymentID", req.PaymentID).Msg("service::ExpireOrder - Order expired successfully")

	return nil
}
