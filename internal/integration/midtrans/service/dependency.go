package service

import (
	midtransPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/midtrans/ports"
	"github.com/midtrans/midtrans-go/snap"
)

var _ midtransPorts.MidtransService = &midtransService{}

type midtransService struct {
	snapClient *snap.Client
}

func NewMidtransService(snapClient *snap.Client) *midtransService {
	return &midtransService{
		snapClient: snapClient,
	}
}
