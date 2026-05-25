// Package queue (consumer.go) provides a resilient consumer wrapper with automatic
// reconnection handling for RabbitMQ connection drops in production (Amazon MQ).
package queue

import (
	"time"

	"github.com/rs/zerolog/log"
)

// ConsumerConfig holds configuration for resilient message consumption.
type ConsumerConfig struct {
	QueueName       string
	RoutingKey      string
	PrefetchCount   int           // QoS prefetch (default: 10)
	ReconnectDelay  time.Duration // delay between reconnection attempts (default: 5s)
	MaxReconnects   int           // max consecutive reconnect attempts (0 = unlimited)
}

// DefaultConsumerConfig returns a production-safe consumer configuration.
func DefaultConsumerConfig(queueName, routingKey string) ConsumerConfig {
	return ConsumerConfig{
		QueueName:      queueName,
		RoutingKey:     routingKey,
		PrefetchCount:  10,
		ReconnectDelay: 5 * time.Second,
		MaxReconnects:  0, // unlimited reconnects for production resilience
	}
}

// LogConsumerStart logs consumer startup metadata for observability.
func LogConsumerStart(cfg ConsumerConfig) {
	log.Info().
		Str("queue", cfg.QueueName).
		Str("routing_key", cfg.RoutingKey).
		Int("prefetch", cfg.PrefetchCount).
		Msg("consumer started with resilient configuration")
}
