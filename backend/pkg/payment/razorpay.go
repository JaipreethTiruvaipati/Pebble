package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyRazorpaySignature validates the webhook payload against the secret.
// This is critical to ensure that top-up webhooks actually came from Razorpay and weren't forged.
func VerifyRazorpaySignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	
	// Use hmac.Equal to prevent timing attacks
	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}
