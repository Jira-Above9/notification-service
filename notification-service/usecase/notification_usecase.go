package usecase

import (
	"context"
	"fmt"
)

type NotificationUsecase interface {
	SendOrderConfirmation(ctx context.Context, orderID, email string) error
}

type notificationUsecase struct {
	emailClient EmailClient
}

type EmailClient interface {
	Send(ctx context.Context, to, subject, body string) error
}

func NewNotificationUsecase(client EmailClient) NotificationUsecase {
	return &notificationUsecase{emailClient: client}
}

func (u *notificationUsecase) SendOrderConfirmation(ctx context.Context, orderID, email string) error {
	subject := "Order Confirmation"
	body := fmt.Sprintf("Your order %s has been confirmed!", orderID)
	return u.emailClient.Send(ctx, email, subject, body)
}
