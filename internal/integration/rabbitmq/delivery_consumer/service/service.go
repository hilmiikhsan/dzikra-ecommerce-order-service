package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rabbitmq/delivery_consumer/ports"
	orderEntity "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/entity"
	orderStatusEntity "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/entity"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *deliveryConsumerService) CheckDelivery(ctx context.Context, req *ports.DeliveryRequest) (delivered bool, err error) {
	res, err := s.rajaOngkirService.GetWaybill(ctx, req.ShippingNumber, strings.ToLower(req.Courier))
	if err != nil {
		log.Error().Err(err).Msg("CheckDelivery - Failed to get waybill from RajaOngkir")
		return false, fmt.Errorf("failed to get waybill: %w", err)
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req.OrderID).Msg("service::CheckDelivery - Failed to begin transaction")
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req.OrderID).Msg("service::CheckDelivery - Failed to rollback transaction")
			}
		}
	}()

	if res.DeliveryStatus.Status == "DELIVERED" || res.DeliveryStatus.PODReceiver != "" {
		orderID, err := uuid.Parse(req.OrderID)
		if err != nil {
			log.Error().Err(err).Msg("CheckDelivery - invalid orderID")
			return false, fmt.Errorf("invalid order id: %w", err)
		}

		if err := s.orderRepository.UpdateStatus(ctx, tx, &orderEntity.Order{
			ID:     orderID,
			Status: constants.OrderStatusDelivered,
		}); err != nil {
			log.Error().Err(err).Msg("CheckDelivery - failed update order status")
			return false, fmt.Errorf("failed to update order status: %w", err)
		}

		orderStatusHistoryPayload := &orderStatusEntity.OrderStatusHistory{
			OrderID:     orderID,
			Status:      constants.OrderStatusDelivered,
			ChangedBy:   "Sistem",
			Description: "Pesanan telah diterima oleh pelanggan.",
		}

		if err := s.orderStatusHistoryRepository.InsertNewOrderStatusHistory(ctx, tx, orderStatusHistoryPayload); err != nil {
			log.Error().Err(err).Msg("CheckDelivery - failed insert history")
			return false, fmt.Errorf("failed to insert order status history: %w", err)
		}

		if err = s.orderItemHistoryRepository.DuplicateOrderItemHistory(ctx, tx, req.OrderID); err != nil {
			log.Error().Err(err).Msg("CheckDelivery - failed duplicate order item history")
			return false, fmt.Errorf("failed to duplicate order item history: %w", err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Any("payload", req.OrderID).Msg("service::CheckDelivery - Failed to commit transaction")
			return false, fmt.Errorf("failed to commit transaction: %w", err)
		}

		go func() {
			_, err := s.externalNotification.SendFcmNotification(ctx, &notification.SendFcmNotificationRequest{
				Title:           "Pesanan telah diterima",
				Body:            fmt.Sprintf("Pesanan %s sudah diterima di tujuan.", req.OrderID),
				UserId:          req.UserID,
				IsStatusChanged: true,
			})
			if err != nil {
				log.Error().Err(err).Msg("CheckDelivery - failed send FCM")
			}
		}()

		log.Info().Str("order", req.OrderID).Msg("CheckDelivery - Order delivered successfully")

		return true, nil
	}

	log.Info().Str("order", req.OrderID).Msg("CheckDelivery - Order not delivered yet")

	return false, nil
}
