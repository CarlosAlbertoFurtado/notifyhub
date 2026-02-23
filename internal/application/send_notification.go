package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CarlosAlbertoFurtado/notifyhub/internal/domain"
)

type SendNotificationInput struct {
	Channel   string            `json:"channel" binding:"required"`
	Recipient string            `json:"recipient" binding:"required"`
	Subject   string            `json:"subject"`
	Body      string            `json:"body" binding:"required"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type SendNotificationUseCase struct {
	repo    domain.NotificationRepository
	senders map[domain.Channel]domain.Sender
}

func NewSendNotificationUseCase(
	repo domain.NotificationRepository,
	senders map[domain.Channel]domain.Sender,
) *SendNotificationUseCase {
	return &SendNotificationUseCase{repo: repo, senders: senders}
}

func (uc *SendNotificationUseCase) Execute(ctx context.Context, input SendNotificationInput) (*domain.Notification, error) {
	n, err := domain.NewNotification(
		domain.Channel(input.Channel),
		input.Recipient,
		input.Subject,
		input.Body,
	)
	if err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	n.Metadata = input.Metadata

	if err := uc.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("persist notification: %w", err)
	}

	go uc.dispatch(n)

	return n, nil
}

func (uc *SendNotificationUseCase) dispatch(n *domain.Notification) {
	ctx := context.Background()

	sender, ok := uc.senders[n.Channel]
	if !ok {
		n.MarkFailed("no sender configured for channel: " + string(n.Channel))
		_ = uc.repo.Update(ctx, n)
		return
	}

	if err := sender.Send(ctx, n); err != nil {
		slog.Error("send failed", "channel", n.Channel, "id", n.ID, "error", err)
		n.MarkFailed(err.Error())
	} else {
		slog.Info("notification sent", "channel", n.Channel, "id", n.ID)
		n.MarkSent()
	}

	_ = uc.repo.Update(ctx, n)
}
