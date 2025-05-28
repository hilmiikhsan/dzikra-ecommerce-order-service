package dto

type RajaOngkirWaybillPayload struct {
	Rajaongkir struct {
		Result struct {
			Summary struct {
				WaybillNumber string `json:"waybill_number"`
				ServiceCode   string `json:"service_code"`
				WaybillDate   string `json:"waybill_date"`
				ShipperName   string `json:"shipper_name"`
				ReceiverName  string `json:"receiver_name"`
				Origin        string `json:"origin"`
				Destination   string `json:"destination"`
				Status        string `json:"status"`
				CourierName   string `json:"courier_name"`
			} `json:"summary"`
			Manifest []struct {
				ManifestDescription string `json:"manifest_description"`
				ManifestDate        string `json:"manifest_date"`
				ManifestTime        string `json:"manifest_time"`
				CityName            string `json:"city_name"`
			} `json:"manifest"`
			DeliveryStatus struct {
				Status      string `json:"status"`
				PODReceiver string `json:"pod_receiver"`
				PODDate     string `json:"pod_date"`
				PODTime     string `json:"pod_time"`
			} `json:"delivery_status"`
		} `json:"result"`
	} `json:"rajaongkir"`
}

// DTO yang akan dikirim ke handler/API layer
type GetWaybillResponse struct {
	Summary        WaybillSummary        `json:"summary"`
	Manifest       []WaybillManifest     `json:"manifest"`
	DeliveryStatus WaybillDeliveryStatus `json:"delivery_status"`
}

type WaybillSummary struct {
	Resi         string `json:"resi"`
	ServiceCode  string `json:"service_code"`
	WaybillDate  string `json:"waybill_date"`
	ShipperName  string `json:"shipper_name"`
	ReceiverName string `json:"receiver_name"`
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	Status       string `json:"status"`
	CourierName  string `json:"courier_name"`
}

type WaybillManifest struct {
	Description string `json:"manifest_description"`
	Date        string `json:"manifest_date"`
	Time        string `json:"manifest_time"`
	CityName    string `json:"city_name"`
}

type WaybillDeliveryStatus struct {
	Status      string `json:"status"`
	PODReceiver string `json:"pod_receiver"`
	PODDate     string `json:"pod_date"`
	PODTime     string `json:"pod_time"`
}
