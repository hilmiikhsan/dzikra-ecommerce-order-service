package address

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/external/proto/address"
)

type ExternalAddress interface {
	GetAddressesByIds(ctx context.Context, req *address.GetAddressesByIdsRequest) (*address.GetAddressesResponse, error)
}
