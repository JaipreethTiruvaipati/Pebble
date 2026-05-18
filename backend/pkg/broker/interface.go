// Package broker abstracts trade execution against Pebble's partner broker (Smallcase).
//
// investment-service wires a SmallcaseClient into PoolExecutor to place sandbox orders
// when pooled penalty cash is invested. A shared Broker interface will allow swapping
// implementations (mock for tests, live API in production) without changing callers.
package broker

// TODO: implement interface
