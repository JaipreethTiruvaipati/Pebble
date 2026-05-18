package allocate

import (
	"testing"

	"github.com/jaipreeth/pebble/backend/internal/models"
)

func TestComputeBrokerOrders_UsesNamedETFs(t *testing.T) {
	orders := ComputeBrokerOrders([]models.MarketSignal{}, 1000)
	if len(orders) != 3 {
		t.Fatalf("expected 3 orders, got %d", len(orders))
	}
	byInst := map[string]bool{}
	for _, o := range orders {
		byInst[o.Instrument] = true
	}
	for _, want := range []string{InstrumentNifty50ETF, InstrumentGoldETF, InstrumentLiquidFund} {
		if !byInst[want] {
			t.Fatalf("missing instrument %s", want)
		}
	}
}
