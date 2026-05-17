package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Connect initializes the RabbitMQ connection and sets up DLQs.
func Connect(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare Dead Letter Exchange
	err = ch.ExchangeDeclare("pebble.dlx", "direct", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare dlx: %w", err)
	}

	// Declare main exchange
	err = ch.ExchangeDeclare("pebble.events", "topic", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare main exchange: %w", err)
	}

	return &RabbitMQ{conn: conn, ch: ch}, nil
}

// Close cleanly shuts down the channel and connection.
func (r *RabbitMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}

// Publish sends a JSON message to a routing key.
func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = r.ch.PublishWithContext(ctx, "pebble.events", routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         b,
	})
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}
	
	log.Info().Str("routing_key", routingKey).Msg("message published to queue")
	return nil
}

// Consume starts a consumer on a specific queue with a routing key, handling DLQ routing.
func (r *RabbitMQ) Consume(queueName, routingKey string, handler func(body []byte) error) error {
	// 1. Declare DLQ
	dlqName := queueName + ".dlq"
	_, err := r.ch.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	r.ch.QueueBind(dlqName, queueName+".dead", "pebble.dlx", false, nil)

	// 2. Declare Main Queue with DLX settings
	args := amqp.Table{
		"x-dead-letter-exchange":    "pebble.dlx",
		"x-dead-letter-routing-key": queueName + ".dead",
	}
	q, err := r.ch.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		return err
	}

	// 3. Bind Main Queue to exchange
	err = r.ch.QueueBind(q.Name, routingKey, "pebble.events", false, nil)
	if err != nil {
		return err
	}

	// 4. Start consuming
	msgs, err := r.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Error().Err(err).Msg("failed to process message, nacking to DLQ")
				d.Nack(false, false) // Reject and send to DLQ
			} else {
				d.Ack(false) // Success
			}
		}
	}()

	log.Info().Str("queue", queueName).Msg("started consuming")
	return nil
}
