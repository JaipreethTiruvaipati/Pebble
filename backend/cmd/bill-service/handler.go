// Package main (handler.go) implements the bill upload HTTP handler that bridges client
// uploads to object storage and the bills.uploaded event for the scoring pipeline.
package main

import (
	"io"
	"net/http"

	"github.com/jaipreeth/pebble/backend/internal/httputil"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/rs/zerolog/log"
)

// HandleBillUpload serves POST /upload: expects multipart fields transaction_id, user_id,
// and file receipt; uploads to S3; publishes bills.uploaded with {transaction_id, user_id, s3_key}.
func HandleBillUpload(uploader *S3Uploader, rmq *queue.RabbitMQ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "file too large or invalid form", "FILE_TOO_LARGE")
			return
		}

		transactionID := r.FormValue("transaction_id")
		userID := r.FormValue("user_id")
		if transactionID == "" || userID == "" {
			httputil.RespondError(w, http.StatusBadRequest, "transaction_id and user_id required", "MISSING_FIELDS")
			return
		}

		file, header, err := r.FormFile("receipt")
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "missing 'receipt' file", "MISSING_FILE")
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to read file", "READ_ERROR")
			return
		}

		s3Key, err := uploader.UploadImage(r.Context(), fileBytes, header.Header.Get("Content-Type"))
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to upload to S3", "S3_ERROR")
			return
		}

		event := map[string]string{
			"transaction_id": transactionID,
			"user_id":        userID,
			"s3_key":         s3Key,
		}
		if err := rmq.Publish(r.Context(), queue.TopicBillsUploaded, event); err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to queue bill", "QUEUE_ERROR")
			return
		}

		log.Info().Str("s3_key", s3Key).Str("transaction_id", transactionID).Msg("bill uploaded and queued")

		httputil.RespondJSON(w, http.StatusAccepted, map[string]string{
			"message": "Bill uploaded successfully and is being processed",
			"s3_key":  s3Key,
		})
	}
}
