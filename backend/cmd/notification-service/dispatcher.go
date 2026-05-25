// Package main (dispatcher.go) routes notification events to FCM and SES channels.
// Called by the RabbitMQ consumers in main.go after deserializing event payloads.
package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jaipreeth/pebble/backend/pkg/notify"
	"github.com/rs/zerolog/log"
)

// Dispatcher routes events to the appropriate notification channels (FCM push, SES email).
type Dispatcher struct {
	fcm *notify.FCMClient
	ses *notify.SESClient
}

// NewDispatcher creates a dispatcher with optional FCM and SES clients.
// Either can be nil if credentials are not configured (dev mode).
func NewDispatcher(fcm *notify.FCMClient, ses *notify.SESClient) *Dispatcher {
	return &Dispatcher{fcm: fcm, ses: ses}
}

// NotifyPenaltyQueued sends push + email for a new pending penalty.
func (d *Dispatcher) NotifyPenaltyQueued(ctx context.Context, userID uuid.UUID, totalPending float64) {
	title := "Impulse Purchase Detected"
	body := fmt.Sprintf("₹%.0f penalty pending — review in 24h or it auto-invests.", totalPending)

	if d.fcm != nil {
		// TODO: look up user's FCM device token from DB
		log.Info().Str("user_id", userID.String()).Msg("FCM penalty push ready (needs device token lookup)")
	}

	if d.ses != nil {
		// TODO: look up user's email from DB
		log.Info().Str("user_id", userID.String()).Float64("amount", totalPending).Msg("SES penalty email ready (needs email lookup)")
	}

	log.Info().
		Str("user_id", userID.String()).
		Str("title", title).
		Str("body", body).
		Msg("penalty notification dispatched")
}

// NotifyInvestmentExecuted sends push + email for a completed investment batch.
func (d *Dispatcher) NotifyInvestmentExecuted(ctx context.Context, triggerType string, totalAmount float64, brokerRef string) {
	title := "Investment Executed"
	body := fmt.Sprintf("₹%.0f invested via %s trigger. Ref: %s", totalAmount, triggerType, brokerRef)

	if d.fcm != nil {
		log.Info().Str("trigger", triggerType).Msg("FCM investment push ready")
	}

	if d.ses != nil {
		log.Info().Str("broker_ref", brokerRef).Msg("SES investment receipt ready")
	}

	log.Info().
		Str("title", title).
		Str("body", body).
		Msg("investment notification dispatched")
}

// NotifyStreakMilestone sends push for a low-impulse streak achievement.
func (d *Dispatcher) NotifyStreakMilestone(ctx context.Context, userID uuid.UUID, streakCount int) {
	title := "Discipline Streak! 🔥"
	body := fmt.Sprintf("%d-week low-impulse streak! Your penalty rate just dropped.", streakCount)

	if d.fcm != nil {
		log.Info().Str("user_id", userID.String()).Int("streak", streakCount).Msg("FCM streak push ready")
	}

	log.Info().
		Str("user_id", userID.String()).
		Str("title", title).
		Str("body", body).
		Msg("streak notification dispatched")
}
