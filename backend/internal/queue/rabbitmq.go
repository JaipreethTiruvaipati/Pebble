// Package queue defines RabbitMQ connectivity, event payloads, and publish/consume helpers
// used to decouple Pebble microservices (bill scoring → penalties → pool investment).
package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// RabbitMQ wraps a persistent AMQP connection and channel for the Pebble event bus.
// Declares pebble.events (topic) and pebble.dlx (direct) exchanges on Connect.
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Connect dials RabbitMQ and declares the Pebble exchanges and dead-letter topology.
//
// Parameters:
//   - url: AMQP URL (e.g. amqp://guest:guest@localhost:5672/)
//
// Returns:
//   - *RabbitMQ: ready client with an open channel
//   - error: dial, channel, or exchange declare failure
//
// How it works: opens connection + channel, declares pebble.dlx (direct) for failed messages
// and pebble.events (topic) for routing keys like bills.scored and investments.executed.
// Each microservice calls Connect once at startup.
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

// Close shuts down the AMQP channel and underlying connection.
//
// Parameters: none (receiver is *RabbitMQ)
//
// Returns: nothing; errors from Close are ignored.
//
// Pebble flow: deferred from service main() on graceful shutdown.
func (r *RabbitMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}

// Publish JSON-marshals body and sends a persistent message to pebble.events.
//
// Parameters:
//   - ctx: passed to PublishWithContext for cancellation
//   - routingKey: topic routing key (use Topic* constants from events.go)
//   - body: event struct (PenaltyQueuedEvent, BillsScoredEvent, etc.)
//
// Returns:
//   - error: marshal or publish failure
//
// Pebble flow: producers (bill-service, penalty-service, investment-service, scoring-service)
// emit domain events consumed by peer services bound to matching queue names.
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

// Consume declares a main queue with DLQ binding and processes deliveries in a goroutine.
//
// Parameters:
//   - queueName: durable queue name (e.g. "penalty-service.bills")
//   - routingKey: binding key on pebble.events (e.g. TopicBillsScored)
//   - handler: callback receiving raw JSON body; return error to Nack → DLQ
//
// Returns:
//   - error: queue declare, bind, or Consume setup failure
//
// How it works: creates {queueName}.dlq bound to pebble.dlx, declares main queue with
// x-dead-letter-exchange, binds to routingKey, then Ack's on success or Nack's to DLQ
// on handler error. Handler runs concurrently in a background goroutine.
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
