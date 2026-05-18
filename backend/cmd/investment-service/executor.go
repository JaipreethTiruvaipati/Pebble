// Package main (executor.go) executes pooled penalty investments: reads Redis signals,
// allocates across instruments, places Smallcase trades, persists investments, and
// publishes investments.executed for notification-service.
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

// PoolExecutor runs micro-batch investment execution and publishes investments.executed.
type PoolExecutor struct {
	db     *pgxpool.Pool
	redis  *cache.Client
	broker *broker.SmallcaseClient
	rmq    *queue.RabbitMQ
}

// NewPoolExecutor wires PostgreSQL, Redis (market signals), Smallcase broker, and RabbitMQ.
func NewPoolExecutor(db *pgxpool.Pool, redis *cache.Client, b *broker.SmallcaseClient, rmq *queue.RabbitMQ) *PoolExecutor {
	return &PoolExecutor{db: db, redis: redis, broker: b, rmq: rmq}
}

// ExecutePool sums un-invested pool contributions, allocates via market signals, executes
// broker trades, marks rows invested with triggerSource, and publishes investments.executed.
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

	orders := allocate.ComputeBrokerOrders(signals, totalPooled)
	splits := allocate.OrdersToAssetSplits(orders)

	brokerRef := fmt.Sprintf("pebble-%s-%d", triggerSource, time.Now().Unix())
	for _, order := range orders {
		if order.Amount <= 0 {
			continue
		}
		if err := e.broker.ExecuteTrade(order.Instrument, order.Amount); err != nil {
			return fmt.Errorf("broker trade %s: %w", order.Instrument, err)
		}
	}

	investmentIDs, err := queries.MarkPoolInvested(ctx, e.db, triggerSource, brokerRef, splits)
	if err != nil {
		return err
	}

	allocation := make([]queue.InvestmentAllocation, 0, len(orders))
	for _, o := range orders {
		if o.Amount > 0 {
			allocation = append(allocation, queue.InvestmentAllocation{
				AssetClass: o.Instrument,
				Amount:     o.Amount,
			})
		}
	}
	event := queue.InvestmentExecutedEvent{
		TriggerType:   triggerSource,
		TotalAmount:   totalPooled,
		BrokerRef:     brokerRef,
		InvestmentIDs: investmentIDs,
		Allocation:    allocation,
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

// defaultSignals returns a neutral HOLD baseline when Redis has no market-poller cache.
func defaultSignals() ([]models.MarketSignal, error) {
	return []models.MarketSignal{
		{AssetClass: "equity", Indicator: "baseline", Value: 50, Action: "HOLD", Timestamp: time.Now()},
	}, nil
}
