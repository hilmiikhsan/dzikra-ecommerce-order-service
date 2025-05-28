package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/notification"
	payment "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/payment_pos"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/user_fcm_token"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/dto"
	orderPaymentEntity "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/entity"
	orderStatusHistoryEntity "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_status_history/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *orderPaymentService) HandleMidtransCallback(ctx context.Context, req *dto.MidtransCallbackRequest) error {
	pay, err := s.orderPaymentRepository.GetLatestByOrderID(ctx, req.OrderID)
	if err != nil {
		log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to get latest order payment")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		log.Error().Err(err).Msg("service::HandleMidtransCallback - Invalid order ID format")
		return err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidOrderIDFormat))
	}

	orderResult, err := s.orderRepository.FindOrderByID(ctx, orderID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrOrderNotFound) {
			log.Warn().Err(err).Msg("service::HandleMidtransCallback - Order not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrOrderNotFound))
		}

		log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to find order by ID")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::HandleMidtransCallback - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::HandleMidtransCallback - Failed to rollback transaction")
			}
		}
	}()

	if pay.ApplicationID == constants.AppPOS {
		res, err := s.externalUserFcmToken.GetUserFcmTokenByUserID(ctx, &user_fcm_token.GetUserFcmTokenByUserIDRequest{
			UserId: orderResult.UserID.String(),
		})
		if err != nil {
			log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to get user fcm token")
			return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		_, err = s.externalPaymentCallback.PaymentCallback(ctx, &payment.PaymentCallbackRequest{
			ApplicationId: constants.AppPOS,
			PaymentId:     pay.ID.String(),
			TransactionId: req.TransactionID,
			Status:        req.TransactionStatus,
			UserFcmToken:  res.UserFcmToken,
			UserId:        res.UserId,
			FullName:      res.FullName,
			Email:         res.Email,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to call payment callback for POS application")
			return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	} else {
		if req.TransactionStatus == constants.MidtransSettlement {
			payloadUpdateStatus := &entity.Order{
				ID:     orderResult.ID,
				Status: constants.OrderStatusProcess,
			}

			if err := s.orderRepository.UpdateStatus(ctx, tx, payloadUpdateStatus); err != nil {
				log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to update order status to process")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			payloadStatusHistlory := &orderStatusHistoryEntity.OrderStatusHistory{
				OrderID:     orderResult.ID,
				Status:      constants.OrderStatusPaymentConfirmed,
				Description: "Pembayaran telah dikonfirmasi.",
				ChangedBy:   "Sistem",
			}

			secondPayloadStatusHistlory := &orderStatusHistoryEntity.OrderStatusHistory{
				OrderID:     orderResult.ID,
				Status:      constants.OrderStatusProcess,
				Description: "Pesanan sedang diproses.",
				ChangedBy:   "Sistem",
			}

			if err := s.orderStatusHistoryRepository.InsertNewOrderStatusHistory(ctx, tx, payloadStatusHistlory); err != nil {
				log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to insert new order status history")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			if err := s.orderStatusHistoryRepository.InsertNewOrderStatusHistory(ctx, tx, secondPayloadStatusHistlory); err != nil {
				log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to insert second order status history")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			res, err := s.externalUserFcmToken.GetUserFcmTokenByUserID(ctx, &user_fcm_token.GetUserFcmTokenByUserIDRequest{
				UserId: orderResult.UserID.String(),
			})
			if err != nil {
				log.Error().Err(err).Msg("service::HandleMidtransCallback - Failed to get user fcm token")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			go func() {
				notifFcmRequest := &notification.SendFcmNotificationRequest{
					FcmToken:        res.UserFcmToken,
					Title:           "Pembayaran telah dikonfirmasi",
					Body:            fmt.Sprintf("Pesanan dengan id %s akan segera diproses.", orderResult.ID),
					UserId:          orderResult.UserID.String(),
					IsStatusChanged: true,
				}

				if _, err := s.externalNotification.SendFcmNotification(ctx, notifFcmRequest); err != nil {
					log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to send FCM notification")
				}
			}()

			parsedTime, _ := time.Parse(time.RFC3339, req.TransactionTime)
			paymentStatus := constants.MidtransPending

			switch req.TransactionStatus {
			case constants.MidtransSettlement:
				paymentStatus = constants.MidtransSuccess
			case constants.MidtransExpire:
				paymentStatus = constants.MidtransExpire
			}

			transactionUUID, err := uuid.Parse(req.TransactionID)
			if err != nil {
				log.Error().Err(err).Msg("service::HandleMidtransCallback - Invalid transaction ID format")
				return err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidTransactionIDFormat))
			}

			err = s.orderPaymentRepository.UpdateOrderPayment(ctx, tx, &orderPaymentEntity.OrderPayment{
				ID:                pay.ID,
				TransactionID:     transactionUUID,
				TransactionStatus: req.TransactionStatus,
				PaymentType:       req.PaymentType,
				SignatureKey:      req.SignatureKey,
				TransactionTime:   parsedTime,
				PaymentStatus:     paymentStatus,
			})
			if err != nil {
				log.Error().Err(err).Msg("service::HandleMidtransCallback - failed update payment record")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::HandleMidtransCallback - Failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
