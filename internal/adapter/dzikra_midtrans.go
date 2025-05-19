package adapter

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog/log"
)

func WithDzikraMidtrans() Option {
	return func(a *Adapter) {
		cfg := config.Envs.Midtrans
		env := midtrans.Sandbox
		if cfg.IsProd {
			env = midtrans.Production
		}

		var client snap.Client
		client.New(cfg.ServerKey, env)

		a.DzikraMidtrans = &client
		log.Info().Msg("Midtrans Snap client initialized")
	}
}
