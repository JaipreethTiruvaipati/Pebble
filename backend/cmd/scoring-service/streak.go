// Package main (streak.go) evaluates low-impulse weekly streaks after each bill is scored
// and publishes streak.updated for notification-service and penalty-rate discounts.
package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/rs/zerolog/log"
)

const (
	lowImpulseAvgScoreMax = 45.0 // avg line-item score below this = disciplined week
	minWeekTransactions   = 1
	streakCooldownDays    = 7
)

// EvaluateWeeklyStreak checks whether the user's rolling week average impulse score is at or
// below lowImpulseAvgScoreMax; if so increments streak_count and publishes streak.updated.
func EvaluateWeeklyStreak(ctx context.Context, db *pgxpool.Pool, rmq *queue.RabbitMQ, userID uuid.UUID) error {
	user, err := queries.GetUserByID(ctx, db, userID)
	if err != nil || user == nil {
		return err
	}

	if user.StreakLastUpdated != nil {
		if time.Since(*user.StreakLastUpdated) < streakCooldownDays*24*time.Hour {
			return nil // already evaluated this week
		}
	}

	stats, err := queries.GetWeekImpulseStats(ctx, db, userID, user.PenaltyThreshold)
	if err != nil {
		return err
	}
	if stats.TxCount < minWeekTransactions {
		return nil
	}
	if stats.AvgScore > lowImpulseAvgScoreMax {
		return nil
	}

	newCount, err := queries.IncrementStreak(ctx, db, userID)
	if err != nil {
		return err
	}

	log.Info().
		Int("streak", newCount).
		Float64("week_avg", stats.AvgScore).
		Str("user_id", userID.String()).
		Msg("low-impulse week — streak updated")

	return rmq.Publish(ctx, queue.TopicStreakUpdated, queue.StreakUpdatedEvent{
		UserID:       userID,
		StreakCount:  newCount,
		WeekAvgScore: stats.AvgScore,
	})
}
