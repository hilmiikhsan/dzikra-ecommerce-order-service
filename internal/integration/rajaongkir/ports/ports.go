package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rajaongkir/dto"
)

type RajaongkirService interface {
	GetWaybill(ctx context.Context, waybill, courier string) (*dto.GetWaybillResponse, error)
}
