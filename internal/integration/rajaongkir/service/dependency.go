package service

import (
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/redis/ports"
	rajaongkirPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/integration/rajaongkir/ports"
)

var _ rajaongkirPorts.RajaongkirService = &rajaongkirService{}

type rajaongkirService struct {
	redisRepository redisPorts.RedisRepository
}

func NewRajaongkirService(
	redisRepository redisPorts.RedisRepository,
) *rajaongkirService {
	return &rajaongkirService{
		redisRepository: redisRepository,
	}
}
