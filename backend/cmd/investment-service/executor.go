package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/cache"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/models"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/jaipreeth/pebble/backend/pkg/allocate"
	"github.com/jaipreeth/pebble/backend/pkg/broker"
	"github.com/rs/zerolog/log"
)

// PoolExecutor runs micro-batch investment execution and publishes completion events.
type PoolExecutor struct {
	db     *pgxpool.Pool
	redis  *cache.Client
	broker *broker.SmallcaseClient
	rmq    *queue.RabbitMQ
}

// NewPoolExecutor wires dependencies for background triggers and queue consumers.
func NewPoolExecutor(db *pgxpool.Pool, redis *cache.Client, b *broker.SmallcaseClient, rmq *queue.RabbitMQ) *PoolExecutor {
	return &PoolExecutor{db: db, redis: redis, broker: b, rmq: rmq}
}

// ExecutePool pulls pooled penalties, allocates, trades via Smallcase, persists investments, and notifies.
func (e *PoolExecutor) ExecutePool(ctx context.Context, triggerSource string) error {
	log.Info().Str("trigger", triggerSource).Msg("initiating investment pool micro-batch execution")

	totalPooled, err := queries.SumPooledAmount(ctx, e.db)
	if err != nil {
		return err
	}
	if totalPooled <= 0 {
		log.Info().Msg("no pooled funds to invest")
		return nil
	}

	var signals []models.MarketSignal
	if e.redis != nil {
		_, _ = e.redis.GetJSON(ctx, cache.KeyMarketSignals, &signals)
	}
	if len(signals) == 0 {
		signals, _ = defaultSignals()
	}

	alloc := allocate.ComputeAllocation(signals)
	splits := map[string]float64{
		"equity": totalPooled * (alloc.Equity / 100.0),
		"gold":   totalPooled * (alloc.Gold / 100.0),
		"bonds":  totalPooled * (alloc.Bonds / 100.0),
	}

	brokerRef := fmt.Sprintf("pebble-%s-%d", triggerSource, time.Now().Unix())
	for asset, amt := range splits {
		if amt <= 0 {
			continue
		}
		if err := e.broker.ExecuteTrade(asset, amt); err != nil {
			return fmt.Errorf("broker trade %s: %w", asset, err)
		}
	}

	investmentIDs, err := queries.MarkPoolInvested(ctx, e.db, triggerSource, brokerRef, splits)
	if err != nil {
		return err
	}

	event := queue.InvestmentExecutedEvent{
		TriggerType:   triggerSource,
		TotalAmount:   totalPooled,
		BrokerRef:     brokerRef,
		InvestmentIDs: investmentIDs,
		Allocation: []queue.InvestmentAllocation{
			{AssetClass: "equity", Amount: splits["equity"]},
			{AssetClass: "gold", Amount: splits["gold"]},
			{AssetClass: "bonds", Amount: splits["bonds"]},
		},
	}
	if err := e.rmq.Publish(ctx, queue.TopicInvestmentsExecuted, event); err != nil {
		return err
	}

	log.Info().
		Str("broker_ref", brokerRef).
		Float64("total", totalPooled).
		Int("investments", len(investmentIDs)).
		Msg("investment pool executed successfully")
	return nil
}

func defaultSignals() ([]models.MarketSignal, error) {
	return []models.MarketSignal{
		{AssetClass: "equity", Indicator: "baseline", Value: 50, Action: "HOLD", Timestamp: time.Now()},
	}, nil
}
