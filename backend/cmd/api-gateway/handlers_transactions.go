package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
)

func handleListTransactions(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		items, err := queries.ListTransactions(r.Context(), db, userID, limit)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to list transactions", "INTERNAL_ERROR")
			return
		}
		if items == nil {
			items = []queries.TransactionSummary{}
		}
		httputil.RespondJSON(w, http.StatusOK, items)
	}
}

func handleGetTransaction(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		txID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid transaction id", "INVALID_ID")
			return
		}
		detail, err := queries.GetTransactionDetail(r.Context(), db, userID, txID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load transaction", "INTERNAL_ERROR")
			return
		}
		if detail == nil {
			httputil.RespondError(w, http.StatusNotFound, "transaction not found", "NOT_FOUND")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, detail)
	}
}

func handleConfirmTransaction(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		txID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid transaction id", "INVALID_ID")
			return
		}
		detail, err := queries.GetTransactionDetail(r.Context(), db, userID, txID)
		if err != nil || detail == nil {
			httputil.RespondError(w, http.StatusNotFound, "transaction not found", "NOT_FOUND")
			return
		}
		var totalPenalty float64
		penalties, _ := queries.ListPenalties(r.Context(), db, userID, "pending")
		for _, p := range penalties {
			for _, li := range detail.LineItems {
				if li.ID == p.LineItemID {
					totalPenalty += p.Amount
				}
			}
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"transaction_id":       txID,
			"penalties_created":    len(detail.LineItems),
			"total_penalty_queued": totalPenalty,
			"status":               "confirmed",
		})
	}
}

func handleOverrideScore(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		lineID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid line item id", "INVALID_ID")
			return
		}
		var req struct {
			OverrideScore int `json:"override_score"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid body", "INVALID_REQUEST")
			return
		}
		if req.OverrideScore < 0 || req.OverrideScore > 100 {
			httputil.RespondError(w, http.StatusBadRequest, "score must be 0-100", "INVALID_REQUEST")
			return
		}
		if err := queries.OverrideLineItemScore(r.Context(), db, userID, lineID, req.OverrideScore); err != nil {
			if err == pgx.ErrNoRows {
				httputil.RespondError(w, http.StatusNotFound, "line item not found", "NOT_FOUND")
				return
			}
			httputil.RespondError(w, http.StatusInternalServerError, "failed to update score", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"line_item_id":    lineID,
			"override_score":  req.OverrideScore,
			"user_overridden": true,
		})
	}
}

func handleUploadBill(cfg *config.Config, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid multipart form", "INVALID_REQUEST")
			return
		}
		merchant := r.FormValue("merchant")
		totalStr := r.FormValue("total_amount")
		if merchant == "" {
			httputil.RespondError(w, http.StatusBadRequest, "merchant required", "INVALID_REQUEST")
			return
		}
		total, _ := strconv.ParseFloat(totalStr, 64)
		if total <= 0 {
			total = 1
		}
		txID, err := queries.CreateTransaction(r.Context(), db, userID, merchant, total)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to create transaction", "INTERNAL_ERROR")
			return
		}

		file, header, err := r.FormFile("receipt")
		if err != nil {
			httputil.RespondJSON(w, http.StatusAccepted, map[string]interface{}{
				"transaction_id": txID,
				"message":        "transaction created; upload receipt to score",
			})
			return
		}
		defer file.Close()

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_ = writer.WriteField("transaction_id", txID.String())
		_ = writer.WriteField("user_id", userID.String())
		part, _ := writer.CreateFormFile("receipt", header.Filename)
		_, _ = io.Copy(part, file)
		writer.Close()

		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, cfg.BillServiceURL+"/upload", &body)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to forward bill", "INTERNAL_ERROR")
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			httputil.RespondJSON(w, http.StatusAccepted, map[string]interface{}{
				"transaction_id": txID,
				"message":        "transaction created; bill-service unreachable — start bill-service on :8081",
			})
			return
		}
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			httputil.RespondJSON(w, http.StatusAccepted, map[string]interface{}{
				"transaction_id": txID,
				"message":        string(respBody),
			})
			return
		}
		var billResp map[string]interface{}
		_ = json.Unmarshal(respBody, &billResp)
		billResp["transaction_id"] = txID
		httputil.RespondJSON(w, http.StatusAccepted, billResp)
	}
}
