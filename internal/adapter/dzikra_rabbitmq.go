package adapter

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

func WithRabbitMQ() Option {
	return func(a *Adapter) {
		uri := config.Envs.RabbitMQ.URI

		conn, err := amqp.Dial(uri)
		if err != nil {
			log.Fatal().Err(err).Msg("adapter::WithRabbitMQ - Failed to connect to RabbitMQ")
		}

		a.RabbitMQConn = conn
		log.Info().Msg("adapter::WithRabbitMQ - RabbitMQ connected")
	}
}
