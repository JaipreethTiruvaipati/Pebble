package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// RespondJSON sends a JSON response with the given status code.
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

// RespondError sends a standard error JSON response.
func RespondError(w http.ResponseWriter, status int, message, code string) {
	RespondJSON(w, status, map[string]string{
		"error": message,
		"code":  code,
	})
}
