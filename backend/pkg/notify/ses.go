// Package notify delivers outbound user notifications for Pebble microservices.
//
// ses.go implements AWS Simple Email Service for transactional mail (investment
// receipts, penalty summaries). notification-service uses SES_FROM_EMAIL from config.
package notify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/rs/zerolog/log"
)

// SESClient wraps AWS SES for transactional email delivery.
type SESClient struct {
	client    *ses.Client
	fromEmail string
}

// NewSESClient creates an SES client using the default AWS credential chain
// (IAM roles in ECS, env vars in dev). fromEmail is the verified sender address.
func NewSESClient(ctx context.Context, region, fromEmail string) (*SESClient, error) {
	if fromEmail == "" {
		return nil, fmt.Errorf("SES_FROM_EMAIL is not set")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("aws config failed: %w", err)
	}

	return &SESClient{
		client:    ses.NewFromConfig(cfg),
		fromEmail: fromEmail,
	}, nil
}

// SendEmail sends a transactional email via AWS SES.
func (s *SESClient) SendEmail(ctx context.Context, to, subject, htmlBody, textBody string) error {
	charset := "UTF-8"

	input := &ses.SendEmailInput{
		Source: &s.fromEmail,
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: &charset,
				Data:    &subject,
			},
			Body: &types.Body{
				Html: &types.Content{
					Charset: &charset,
					Data:    &htmlBody,
				},
				Text: &types.Content{
					Charset: &charset,
					Data:    &textBody,
				},
			},
		},
	}

	result, err := s.client.SendEmail(ctx, input)
	if err != nil {
		log.Error().Err(err).Str("to", to).Str("subject", subject).Msg("SES send failed")
		return fmt.Errorf("ses send failed: %w", err)
	}

	log.Info().
		Str("message_id", *result.MessageId).
		Str("to", to).
		Str("subject", subject).
		Msg("SES email sent")
	return nil
}

// SendPenaltyAlert sends a penalty consent notification email.
func (s *SESClient) SendPenaltyAlert(ctx context.Context, to string, totalPending float64) error {
	subject := fmt.Sprintf("Pebble: ₹%.0f penalty pending your review", totalPending)
	htmlBody := fmt.Sprintf(`
		<h2>Impulse Purchase Detected</h2>
		<p>Your recent transaction has been scored and a penalty of <strong>₹%.0f</strong> is pending.</p>
		<p>You have 24 hours to review and optionally cancel this penalty in the Pebble app.</p>
		<p>If no action is taken, the amount will be automatically invested on your behalf.</p>
		<br>
		<p style="color: #888;">— Pebble, your financial discipline partner</p>
	`, totalPending)
	textBody := fmt.Sprintf("Pebble: A penalty of Rs %.0f is pending your review. Open the app within 24 hours to review.", totalPending)

	return s.SendEmail(ctx, to, subject, htmlBody, textBody)
}

// SendInvestmentReceipt sends an investment confirmation email.
func (s *SESClient) SendInvestmentReceipt(ctx context.Context, to string, totalAmount float64, brokerRef string) error {
	subject := fmt.Sprintf("Pebble: ₹%.0f invested successfully", totalAmount)
	htmlBody := fmt.Sprintf(`
		<h2>Investment Confirmation</h2>
		<p>₹%.0f has been invested across your portfolio pools.</p>
		<p><strong>Reference:</strong> %s</p>
		<p>Check your Pebble app for the full allocation breakdown.</p>
		<br>
		<p style="color: #888;">— Pebble, your financial discipline partner</p>
	`, totalAmount, brokerRef)
	textBody := fmt.Sprintf("Pebble: Rs %.0f invested. Ref: %s. Check the app for details.", totalAmount, brokerRef)

	return s.SendEmail(ctx, to, subject, htmlBody, textBody)
}
