// Package market ingests Indian market data and produces allocation signals for Pebble.
//
// market-poller calls fetchers on a schedule, then ComputeOpportunitySignals; results are
// cached in Redis (cache.KeyMarketSignals) for investment-service PoolExecutor.
package market

// FetchAMFIData downloads and parses AMFI daily NAV files for mutual fund benchmarks.
//
// Bond sleeve allocation and liquid-fund context will use this data once Phase 2 parsing
// is implemented. market-poller invokes this alongside NSE, MCX, and CCIL fetchers.
func FetchAMFIData() error {
	// TODO: Fetch and parse the AMFI daily NAV text file
	return nil
}
