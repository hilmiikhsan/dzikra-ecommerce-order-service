package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type ExpireMessage struct {
	FcmToken  string    `json:"fcm_token"`
	UserID    string    `json:"user_id"`
	PaymentID string    `json:"payment_id"`
	OrderID   string    `json:"order_id"`
	ExpiredAt time.Time `json:"expired_at"`
}

type ExpireOrderRequest struct {
	FcmToken  string    `json:"fcm_token"`
	UserID    string    `json:"user_id"`
	PaymentID string    `json:"payment_id"`
	OrderID   string    `json:"order_id"`
	ExpiredAt time.Time `json:"expired_at"`
}

type ConsumerService interface {
	ExpireOrder(ctx context.Context, req *ExpireOrderRequest) error
}

func StartExpireConsumer(
	conn *amqp.Connection,
	db *sqlx.DB,
	_ /*redis*/ interface{},
	consumerSerice ConsumerService,
) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("expireConsumer: failed to open a channel")
		return err
	}

	if err := ch.ExchangeDeclare(
		"payment.expire.ex", // nama exchange
		"x-delayed-message", // type
		true,                // durable
		false,
		false,
		false,
		amqp.Table{"x-delayed-type": "direct"},
	); err != nil {
		log.Error().Err(err).Msg("expireConsumer: failed to declare exchange")
		return err
	}

	const queueName = "payment_expire_queue"
	q, err := ch.QueueDeclare(
		queueName, // nama queue
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("expireConsumer: failed to declare queue")
		return err
	}

	if err := ch.QueueBind(
		q.Name,
		"payment.expire",    // routing key
		"payment.expire.ex", // exchange
		false,
		nil,
	); err != nil {
		log.Error().Err(err).Msg("expireConsumer: failed to bind queue to exchange")
		return err
	}

	if err := ch.Qos(1, 0, false); err != nil {
		log.Error().Err(err).Msg("expireConsumer: failed to set QoS")
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("expireConsumer: failed to register consumer")
		return err
	}

	go func() {
		for d := range msgs {
			var m ExpireMessage
			if err := json.Unmarshal(d.Body, &m); err != nil {
				log.Error().Err(err).Msg("expireConsumer: invalid payload")
				d.Nack(false, false)
				continue
			}

			log.Info().
				Str("order", m.OrderID).
				Str("payment", m.PaymentID).
				Time("expiredAt", m.ExpiredAt).
				Msg("expireConsumer: received expire message")

			ctx := context.Background()
			if err := consumerSerice.ExpireOrder(ctx, &ExpireOrderRequest{
				FcmToken:  m.FcmToken,
				UserID:    m.UserID,
				PaymentID: m.PaymentID,
				OrderID:   m.OrderID,
				ExpiredAt: m.ExpiredAt,
			}); err != nil {
				log.Error().Err(err).Msg("expireConsumer: failed ExpireOrder")
				d.Nack(false, true) // requeue
				continue
			}

			d.Ack(false)
		}
	}()

	log.Info().Msg("expireConsumer: waiting for messages")

	return nil
}
