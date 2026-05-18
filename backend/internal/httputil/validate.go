// Package httputil provides shared HTTP response helpers and request validation for Pebble API handlers.
package httputil

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ErrEmptyField indicates a required string field was blank after trimming.
var ErrEmptyField = errors.New("field cannot be empty")

// ErrFieldTooLong indicates a string exceeded its maximum rune length.
var ErrFieldTooLong = errors.New("field exceeds maximum length")

// ErrInvalidAmount indicates a monetary amount was zero or negative.
var ErrInvalidAmount = errors.New("amount must be positive")

// ValidateNonEmpty trims and validates a required string field from an API request body.
//
// Parameters:
//   - field: logical name used in error messages (e.g. "merchant")
//   - value: raw input from JSON or form
//   - maxLen: maximum rune count; 0 disables length check
//
// Returns:
//   - string: trimmed value on success
//   - error: ErrEmptyField or ErrFieldTooLong wrapped with field name
func ValidateNonEmpty(field, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s: %w", field, ErrEmptyField)
	}
	if maxLen > 0 && utf8.RuneCountInString(value) > maxLen {
		return "", fmt.Errorf("%s: %w", field, ErrFieldTooLong)
	}
	return value, nil
}

// ValidateAmount ensures a monetary amount is positive and within an optional ceiling.
//
// Parameters:
//   - amount: INR value from request (bill total, top-up, etc.)
//   - max: upper bound; 0 means no maximum check
//
// Returns:
//   - nil when amount > 0 and (max == 0 or amount <= max)
//   - ErrInvalidAmount or a formatted exceed error
//
// Pebble flow: bill logging and wallet top-up handlers before hitting queries.CreateTransaction.
func ValidateAmount(amount float64, max float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if max > 0 && amount > max {
		return fmt.Errorf("amount exceeds maximum of %.2f", max)
	}
	return nil
}

// ValidateReferralCode normalises and validates a referral code before DB redemption.
//
// Parameters:
//   - code: raw code from POST /referrals/redeem
//
// Returns:
//   - string: trimmed uppercase code on success
//   - error: empty or too-long codes (max 32 runes)
//
// Pebble flow: api-gateway referral handler before queries.RedeemReferralCode.
func ValidateReferralCode(code string) (string, error) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return "", fmt.Errorf("code: %w", ErrEmptyField)
	}
	if utf8.RuneCountInString(code) > 32 {
		return "", fmt.Errorf("code: %w", ErrFieldTooLong)
	}
	return code, nil
}
