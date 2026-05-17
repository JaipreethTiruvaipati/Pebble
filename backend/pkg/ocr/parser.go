package ocr

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jaipreeth/pebble/backend/internal/models"
	"github.com/rs/zerolog/log"
)

// ParseIndianReceipt applies Regex and heuristics to extract line items and the merchant from raw OCR text.
func ParseIndianReceipt(rawText string) (string, []models.ScoredItem) {
	log.Info().Msg("parsing raw OCR text into line items")
	// TODO: Phase 1 - Implement robust Regex parsing for Indian receipts (GST, FSSAI, Swiggy/Zomato formats)
	
	// Very naive stub implementation
	merchant := "Unknown Merchant"
	lines := strings.Split(rawText, "\n")
	if len(lines) > 0 {
		merchant = lines[0]
	}

	items := []models.ScoredItem{}
	
	// Example regex for finding prices like "400" or "400.00"
	priceRegex := regexp.MustCompile(`(\d+(\.\d{1,2})?)`)
	
	for _, line := range lines[1:] {
		if strings.Contains(strings.ToLower(line), "total") {
			continue
		}
		
		match := priceRegex.FindString(line)
		if match != "" {
			amount, _ := strconv.ParseFloat(match, 64)
			items = append(items, models.ScoredItem{
				Name:   strings.ReplaceAll(line, match, ""), // roughly the item name
				Amount: amount,
			})
		}
	}

	return merchant, items
}
