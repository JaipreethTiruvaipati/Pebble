// Package main (pool.go) exposes a package-level ExecutePool entry for background triggers
// (threshold, time, opportunity) without importing main's initialization cycle.
package main

import "context"

// globalExecutor is set from main; triggers call ExecutePool without circular imports.
var globalExecutor *PoolExecutor

// ExecutePool delegates to PoolExecutor.ExecutePool with triggerSource ("threshold", "time", "opportunity").
func ExecutePool(triggerSource string) {
	if globalExecutor == nil {
		return
	}
	_ = globalExecutor.ExecutePool(context.Background(), triggerSource)
}
