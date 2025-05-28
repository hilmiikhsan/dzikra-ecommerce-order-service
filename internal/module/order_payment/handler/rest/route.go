package rest

import "github.com/gofiber/fiber/v2"

func (h *orderPaymentHandler) OrderPaymentRoute(publicRouter fiber.Router) {
	publicRouter.Post("/payment/callback", h.midtransCallback)
}
