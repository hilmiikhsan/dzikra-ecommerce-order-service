package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rajaongkir/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *rajaongkirService) GetWaybill(ctx context.Context, waybill, courier string) (*dto.GetWaybillResponse, error) {
	values := url.Values{}
	values.Set("waybill", waybill)
	values.Set("courier", courier)

	endpoint := fmt.Sprintf("%s/waybill", config.Envs.RajaOngkir.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBufferString(values.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("service::GetWaybill - NewRequest failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("service::GetWaybill - HTTP request failed")
		return nil, err_msg.NewCustomErrors(http.StatusBadGateway, err_msg.WithMessage(constants.ErrExternalServiceUnavailable))
	}
	defer resp.Body.Close()

	var payload dto.RajaOngkirWaybillPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("service::GetWaybill - JSON decode failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	r := payload.Rajaongkir.Result

	out := &dto.GetWaybillResponse{
		Summary: dto.WaybillSummary{
			Resi:         r.Summary.WaybillNumber,
			ServiceCode:  r.Summary.ServiceCode,
			WaybillDate:  r.Summary.WaybillDate,
			ShipperName:  r.Summary.ShipperName,
			ReceiverName: r.Summary.ReceiverName,
			Origin:       r.Summary.Origin,
			Destination:  r.Summary.Destination,
			Status:       r.Summary.Status,
			CourierName:  r.Summary.CourierName,
		},
		Manifest: make([]dto.WaybillManifest, len(r.Manifest)),
		DeliveryStatus: dto.WaybillDeliveryStatus{
			Status:      r.DeliveryStatus.Status,
			PODReceiver: r.DeliveryStatus.PODReceiver,
			PODDate:     r.DeliveryStatus.PODDate,
			PODTime:     r.DeliveryStatus.PODTime,
		},
	}

	for i, m := range r.Manifest {
		out.Manifest[i] = dto.WaybillManifest{
			Description: m.ManifestDescription,
			Date:        m.ManifestDate,
			Time:        m.ManifestTime,
			CityName:    m.CityName,
		}
	}

	return out, nil
}
