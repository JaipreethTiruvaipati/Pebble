package main

import (
	"io"
	"net/http"

	"github.com/jaipreeth/pebble/backend/internal/httputil"
	"github.com/rs/zerolog/log"
)

// HandleBillUpload processes a multipart/form-data request containing a receipt image.
func HandleBillUpload(uploader *S3Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit upload size to 10MB
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "file too large or invalid form", "FILE_TOO_LARGE")
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

		// 1. Upload to S3
		s3Key, err := uploader.UploadImage(r.Context(), fileBytes, header.Header.Get("Content-Type"))
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to upload to S3", "S3_ERROR")
			return
		}

		log.Info().Str("s3_key", s3Key).Msg("receipt uploaded successfully")

		// 2. Publish "bills.uploaded" event to RabbitMQ
		// TODO: Use the internal/queue package to publish the event to the Scoring Service.
		
		httputil.RespondJSON(w, http.StatusAccepted, map[string]string{
			"message": "Bill uploaded successfully and is being processed",
			"s3_key":  s3Key,
		})
	}
}
