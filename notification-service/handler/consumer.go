package handler

import (
	"context"
	"encoding/json"
	"log"
	"notification-service/usecase"
	"github.com/streadway/amqp"
)

type NotificationConsumer struct {
	notificationUsecase usecase.NotificationUsecase
	channel             *amqp.Channel
}

func NewNotificationConsumer(u usecase.NotificationUsecase, ch *amqp.Channel) *NotificationConsumer {
	return &NotificationConsumer{
		notificationUsecase: u,
		channel:             ch,
	}
}

// ConsumeOrderEvents consumes order events from the message queue
// BUG: Messages are not being acknowledged properly, causing the queue to
// redeliver the same message multiple times. The ack call happens BEFORE
// processing completes, and on processing failure the message is auto-requeued
// causing duplicate email delivery.
func (c *NotificationConsumer) ConsumeOrderEvents(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		"order.created",
		"",
		true,   // BUG: auto-ack is true, messages are acked before processing
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var event OrderCreatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("failed to unmarshal event: %v", err)
			continue
		}

		// Process the email sending
		// BUG: No deduplication check here - if consumer crashes after sending email
		// but before network ack, the message will be redelivered and email sent again.
		// Missing: Store processed event IDs in cache/database to prevent duplicate sending.
		if err := c.notificationUsecase.SendOrderConfirmation(ctx, event.OrderID, event.Email); err != nil {
			log.Printf("failed to send email: %v", err)
			// Message is already auto-acked so it's lost, but retry still causes duplicates
		}
	}
	return nil
}

type OrderCreatedEvent struct {
	EventID string `json:"event_id"`
	OrderID string `json:"order_id"`
	Email   string `json:"email"`
}
