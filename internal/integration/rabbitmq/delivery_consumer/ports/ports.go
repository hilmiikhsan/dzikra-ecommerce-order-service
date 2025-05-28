package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type DeliveryMessage struct {
	OrderID        string    `json:"order_id"`
	ShippingNumber string    `json:"shipping_number"`
	Courier        string    `json:"courier"`
	NextCheckAt    time.Time `json:"next_check_at"`
	Attempt        int       `json:"attempt"`
	MaxAttempts    int       `json:"max_attempts"`
}

type DeliveryRequest struct {
	FcmToken       string `json:"fcm_token"`
	UserID         string `json:"user_id"`
	OrderID        string `json:"order_id"`
	ShippingNumber string `json:"shipping_number"`
	Courier        string `json:"courier"`
}

type DeliveryConsumerService interface {
	CheckDelivery(ctx context.Context, req *DeliveryRequest) (delivered bool, err error)
}

func StartDeliveryConsumer(
	conn *amqp.Connection,
	db *sqlx.DB,
	_ /*redis*/ interface{},
	svc DeliveryConsumerService,
) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("deliveryConsumer: failed to open a channel")
		return err
	}

	if err := ch.ExchangeDeclare(
		"delivery.check.ex",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{"x-delayed-type": "direct"},
	); err != nil {
		log.Error().Err(err).Msg("deliveryConsumer: declare exchange failed")
		return err
	}

	q, err := ch.QueueDeclare(
		"delivery_check_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("deliveryConsumer: declare queue failed")
		return err
	}
	if err := ch.QueueBind(
		q.Name,
		"delivery.check",
		"delivery.check.ex",
		false,
		nil,
	); err != nil {
		log.Error().Err(err).Msg("deliveryConsumer: bind queue failed")
		return err
	}

	if err := ch.Qos(1, 0, false); err != nil {
		log.Error().Err(err).Msg("deliveryConsumer: set QoS failed")
		return err
	}

	msgs, err := ch.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("deliveryConsumer: consume failed")
		return err
	}

	go func() {
		for d := range msgs {
			var m DeliveryMessage
			if err := json.Unmarshal(d.Body, &m); err != nil {
				log.Error().Err(err).Msg("deliveryConsumer: bad payload")
				d.Nack(false, false)
				continue
			}

			log.Info().
				Str("order", m.OrderID).
				Msgf("deliveryConsumer: attempt %d/%d", m.Attempt, m.MaxAttempts)

			delivered, err := svc.CheckDelivery(context.Background(), &DeliveryRequest{
				OrderID:        m.OrderID,
				ShippingNumber: m.ShippingNumber,
				Courier:        m.Courier,
			})
			if err != nil {
				log.Error().Err(err).Msg("deliveryConsumer: service error")
				d.Nack(false, true)
				continue
			}

			if delivered {
				log.Info().Str("order", m.OrderID).Msg("deliveryConsumer: delivered, done")
				d.Ack(false)
				continue
			}

			if m.Attempt+1 < m.MaxAttempts {
				m.Attempt++
				delay := constants.InitialDelay
				body, err := json.Marshal(m)
				if err != nil {
					log.Error().Err(err).Msg("deliveryConsumer: failed to marshal message")
					d.Nack(false, false)
					continue
				}

				if err := ch.Publish(
					"delivery.check.ex",
					"delivery.check",
					false, false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        body,
						Headers:     amqp.Table{"x-delay": int64(delay / time.Millisecond)},
					},
				); err != nil {
					log.Error().Err(err).Msg("deliveryConsumer: failed to re-publish for retry")
					d.Ack(false)
					continue
				}

				log.Info().Str("order", m.OrderID).Int("next_attempt", m.Attempt).Msg("deliveryConsumer: scheduled next check")
			}

			d.Ack(false)
		}
	}()

	log.Info().Msg("deliveryConsumer: waiting for messages")

	return nil
}
