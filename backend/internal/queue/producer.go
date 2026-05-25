// Package queue (producer.go) provides publishing helpers with delivery confirmation
// and logging for production message reliability.
package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/jaipreeth/pebble/backend/pkg/retry"
	"github.com/rs/zerolog/log"
)

// PublishWithRetry publishes a message to the event bus with exponential backoff retry.
// Used for critical events (bills.scored, investments.executed) where message loss
// would break the pipeline.
func (r *RabbitMQ) PublishWithRetry(ctx context.Context, routingKey string, body interface{}) error {
	return retry.Do(ctx, retry.Default(), fmt.Sprintf("publish-%s", routingKey), func() error {
		return r.Publish(ctx, routingKey, body)
	})
}

// PublishDelayed logs publish latency for observability. In production, this can be
// extended with Prometheus histogram recording.
func (r *RabbitMQ) PublishDelayed(ctx context.Context, routingKey string, body interface{}) error {
	start := time.Now()
	err := r.Publish(ctx, routingKey, body)
	duration := time.Since(start)

	if err != nil {
		log.Error().
			Err(err).
			Str("routing_key", routingKey).
			Dur("duration", duration).
			Msg("publish failed")
		return err
	}

	log.Debug().
		Str("routing_key", routingKey).
		Dur("duration", duration).
		Msg("message published with latency tracking")
	return nil
}
