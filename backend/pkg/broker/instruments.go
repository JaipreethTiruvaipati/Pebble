// Package broker abstracts trade execution against Pebble's partner broker (Smallcase).
//
// instruments.go defines the allowlist of tickers pkg/allocate emits and SmallcaseClient accepts.
package broker

// SupportedInstruments lists Pebble broker codes that may be passed to ExecuteTrade.
// Must stay in sync with allocate.Instrument* constants and Smallcase sandbox catalog.
var SupportedInstruments = []string{
	"NIFTY50_ETF",
	"GOLD_ETF",
	"LIQUID_FUND",
}

// IsValidInstrument reports whether code is a known Pebble/Smallcase instrument.
// ExecuteTrade uses this to reject typos and default invalid codes to LIQUID_FUND.
func IsValidInstrument(code string) bool {
	for _, s := range SupportedInstruments {
		if s == code {
			return true
		}
	}
	return false
}
