// Package allocate implements signal-driven portfolio allocation for Pebble pooled investments.
//
// instruments.go maps percentage allocations to named Smallcase instruments and INR amounts
// for execution by investment-service via pkg/broker.
package allocate

import "github.com/jaipreeth/pebble/backend/internal/models"

// Broker instrument codes (Smallcase / AMC tickers).
const (
	// InstrumentNifty50ETF is the equity sleeve ticker executed via Smallcase sandbox.
	InstrumentNifty50ETF = "NIFTY50_ETF"
	// InstrumentGoldETF is the gold sleeve ticker executed via Smallcase sandbox.
	InstrumentGoldETF = "GOLD_ETF"
	// InstrumentLiquidFund is the bonds/cash sleeve ticker (liquid debt fund).
	InstrumentLiquidFund = "LIQUID_FUND"
)

// BrokerOrder is a single executable order: one asset class, one instrument code,
// and the INR amount and allocation percentage derived from pooled penalties.
type BrokerOrder struct {
	AssetClass string  `json:"asset_class"`
	Instrument string  `json:"instrument"`
	Amount     float64 `json:"amount"`
	Pct        float64 `json:"pct"`
}

// ComputeBrokerOrders turns market signals and total pooled INR into three broker orders.
//
// It calls ComputeAllocation, then assigns each sleeve to a fixed instrument
// (NIFTY50_ETF, GOLD_ETF, LIQUID_FUND) with amounts rounded to two decimal places.
// investment-service PoolExecutor iterates these orders and calls broker.SmallcaseClient.ExecuteTrade.
func ComputeBrokerOrders(signals []models.MarketSignal, totalINR float64) []BrokerOrder {
	alloc := ComputeAllocation(signals)
	orders := []BrokerOrder{
		{
			AssetClass: "equity",
			Instrument: InstrumentNifty50ETF,
			Amount:     round2(totalINR * alloc.Equity / 100.0),
			Pct:        alloc.Equity,
		},
		{
			AssetClass: "gold",
			Instrument: InstrumentGoldETF,
			Amount:     round2(totalINR * alloc.Gold / 100.0),
			Pct:        alloc.Gold,
		},
		{
			AssetClass: "bonds",
			Instrument: InstrumentLiquidFund,
			Amount:     round2(totalINR * alloc.Bonds / 100.0),
			Pct:        alloc.Bonds,
		},
	}
	return orders
}

// round2 rounds v to two decimal places using half-up integer scaling.
// Used internally so broker amounts and DB splits match displayed currency precision.
func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

// OrdersToAssetSplits aggregates broker orders into a map keyed by asset_class for persistence.
//
// Multiple orders for the same class are summed. Zero-amount orders are skipped.
// PoolExecutor passes the result to queries.MarkPoolInvested when recording investments.
func OrdersToAssetSplits(orders []BrokerOrder) map[string]float64 {
	out := make(map[string]float64)
	for _, o := range orders {
		if o.Amount > 0 {
			out[o.AssetClass] += o.Amount
		}
	}
	return out
}
