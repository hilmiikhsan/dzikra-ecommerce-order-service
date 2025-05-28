package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order_payment/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/response"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *orderPaymentHandler) midtransCallback(c *fiber.Ctx) error {
	var (
		req = new(dto.MidtransCallbackRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::midtransCallback - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::midtransCallback - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if !utils.VerifyMidtransSignature(req.OrderID, req.StatusCode, req.GrossAmount, req.SignatureKey, config.Envs.Midtrans.ServerKey) {
		log.Warn().Msg("handler::midtransCallback - Invalid signature key")
		return c.Status(fiber.StatusForbidden).JSON(response.Error(constants.InvalidSignatureKey))
	}

	if err := h.service.HandleMidtransCallback(ctx, req); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::midtransCallback - Failed to midtrans callback")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	log.Info().Msg("handler::midtransCallback - Success to midtrans callback")

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
