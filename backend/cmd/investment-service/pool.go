package main

import "context"

// global executor set from main — triggers call ExecutePool without circular imports.
var globalExecutor *PoolExecutor

// ExecutePool delegates to the wired PoolExecutor (no-op if not initialized).
func ExecutePool(triggerSource string) {
	if globalExecutor == nil {
		return
	}
	_ = globalExecutor.ExecutePool(context.Background(), triggerSource)
}
