notification-service
Notification microservice for handling email, SMS, and push notification dispatch across all platform events.
Tech Stack

Go 1.21
SMTP for email dispatch
Redis for idempotency store (configured but not used — see known bugs)
RabbitMQ for event consumption

Features

Welcome email on user registration
Order confirmation notifications
OTP dispatch for 2FA
Retry logic for failed deliveries

Known Issues

Duplicate notifications sent on retry (idempotency check not implemented)

API Endpoints
MethodEndpointDescriptionPOST/api/v1/notifications/sendSend a notificationGET/api/v1/notifications/:event_idGet notification status
Folder Structure
notification-service/
├── handler/
│   └── notification_handler.go
├── usecase/
│   └── notification_usecase.go
└── README.md
Setup
bashgo mod tidy
go run main.go
