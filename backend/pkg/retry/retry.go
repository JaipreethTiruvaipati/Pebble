// Package retry provides exponential-backoff retry helpers for external API calls
// (Gemini, Razorpay, Smallcase) to handle transient failures in production.
package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

// Config controls retry behavior.
type Config struct {
	MaxAttempts int           // total attempts (1 = no retry)
	BaseDelay  time.Duration // initial backoff delay (doubles each attempt)
	MaxDelay   time.Duration // upper bound on any single backoff wait
	Jitter     bool          // add randomized jitter to prevent thundering herd
}

// Default returns a production-safe retry config for external APIs.
func Default() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:  500 * time.Millisecond,
		MaxDelay:   5 * time.Second,
		Jitter:     true,
	}
}

// Aggressive returns a retry config for critical-path operations (e.g. broker trades).
func Aggressive() Config {
	return Config{
		MaxAttempts: 5,
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
		Jitter:     true,
	}
}

// Do executes fn with exponential backoff. Returns the first nil error or the last
// error after all attempts are exhausted. Respects context cancellation.
func Do(ctx context.Context, cfg Config, operation string, fn func() error) error {
	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			if attempt > 1 {
				log.Info().Str("operation", operation).Int("attempt", attempt).Msg("retry succeeded")
			}
			return nil
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		// Check context before sleeping
		if ctx.Err() != nil {
			return fmt.Errorf("%s: context cancelled after %d attempts: %w", operation, attempt, ctx.Err())
		}

		delay := backoffDelay(attempt, cfg)
		log.Warn().
			Err(lastErr).
			Str("operation", operation).
			Int("attempt", attempt).
			Dur("retry_in", delay).
			Msg("retrying after transient failure")

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return fmt.Errorf("%s: context cancelled during backoff: %w", operation, ctx.Err())
		}
	}

	return fmt.Errorf("%s: all %d attempts failed: %w", operation, cfg.MaxAttempts, lastErr)
}

// DoWithResult executes fn returning (T, error) with exponential backoff.
func DoWithResult[T any](ctx context.Context, cfg Config, operation string, fn func() (T, error)) (T, error) {
	var lastResult T
	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		lastResult, lastErr = fn()
		if lastErr == nil {
			return lastResult, nil
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		if ctx.Err() != nil {
			return lastResult, fmt.Errorf("%s: context cancelled after %d attempts: %w", operation, attempt, ctx.Err())
		}

		delay := backoffDelay(attempt, cfg)
		log.Warn().
			Err(lastErr).
			Str("operation", operation).
			Int("attempt", attempt).
			Dur("retry_in", delay).
			Msg("retrying after transient failure")

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return lastResult, fmt.Errorf("%s: context cancelled during backoff: %w", operation, ctx.Err())
		}
	}

	return lastResult, fmt.Errorf("%s: all %d attempts failed: %w", operation, cfg.MaxAttempts, lastErr)
}

// backoffDelay computes exponential delay with optional jitter.
func backoffDelay(attempt int, cfg Config) time.Duration {
	exp := math.Pow(2, float64(attempt-1))
	delay := time.Duration(float64(cfg.BaseDelay) * exp)
	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}
	if cfg.Jitter {
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay += jitter
	}
	return delay
}
