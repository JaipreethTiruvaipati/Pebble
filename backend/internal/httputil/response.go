// Package httputil provides shared HTTP response helpers and request validation for Pebble API handlers.
// api-gateway and other services use these to return consistent JSON error shapes to the mobile client.
package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// RespondJSON marshals payload to JSON and writes it with the given HTTP status.
//
// Parameters:
//   - w: http.ResponseWriter for the active request
//   - status: HTTP status code (e.g. http.StatusOK, http.StatusCreated)
//   - payload: any JSON-serializable value (struct, map, slice)
//
// Returns: nothing; on marshal failure writes 500 with a fixed INTERNAL_ERROR body.
//
// Side effects: sets Content-Type application/json, WriteHeader(status), then body bytes.
// Used across api-gateway handlers for success responses (portfolio, referrals, bills).
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal JSON response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error", "code": "INTERNAL_ERROR"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// RespondError writes a standard Pebble API error envelope as JSON.
//
// Parameters:
//   - w: http.ResponseWriter for the active request
//   - status: HTTP status (4xx/5xx)
//   - message: human-readable error string for the client
//   - code: machine-readable code (e.g. "INVALID_INPUT", "UNAUTHORIZED")
//
// Returns: nothing; delegates to RespondJSON with {"error", "code"} map.
//
// Pebble flow: handlers return validation failures and auth errors without duplicating JSON shape.
func RespondError(w http.ResponseWriter, status int, message, code string) {
	RespondJSON(w, status, map[string]string{
		"error": message,
		"code":  code,
	})
}
