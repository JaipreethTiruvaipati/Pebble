// Package notify delivers outbound user notifications for Pebble microservices.
//
// fcm.go implements Firebase Cloud Messaging push delivery. notification-service
// calls it when consuming queue events (bills scored, investments executed, streaks).
package notify

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

// FCMClient wraps Firebase Cloud Messaging for push notifications.
type FCMClient struct {
	client *messaging.Client
}

// NewFCMClient initializes Firebase Admin SDK from the service account JSON at credPath
// and returns a ready-to-use FCM messaging client.
func NewFCMClient(ctx context.Context, credPath string) (*FCMClient, error) {
	if credPath == "" {
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS_PATH is not set")
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credPath))
	if err != nil {
		return nil, fmt.Errorf("firebase init failed: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase messaging client failed: %w", err)
	}

	return &FCMClient{client: client}, nil
}

// SendPush delivers a push notification to a specific device registration token.
func (f *FCMClient) SendPush(ctx context.Context, token, title, body string, data map[string]string) error {
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "pebble_alerts",
				Sound:     "default",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound:    "default",
					Category: "pebble_alert",
				},
			},
		},
	}

	messageID, err := f.client.Send(ctx, msg)
	if err != nil {
		log.Error().Err(err).Str("token", token[:10]+"...").Msg("FCM send failed")
		return fmt.Errorf("fcm send failed: %w", err)
	}

	log.Info().Str("message_id", messageID).Str("title", title).Msg("FCM push sent")
	return nil
}

// SendToTopic delivers a push notification to all devices subscribed to a topic.
func (f *FCMClient) SendToTopic(ctx context.Context, topic, title, body string, data map[string]string) error {
	msg := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	messageID, err := f.client.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("fcm topic send failed: %w", err)
	}

	log.Info().Str("message_id", messageID).Str("topic", topic).Msg("FCM topic push sent")
	return nil
}
