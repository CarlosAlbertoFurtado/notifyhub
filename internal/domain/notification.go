package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Channel string

const (
	ChannelEmail   Channel = "email"
	ChannelSMS     Channel = "sms"
	ChannelWebhook Channel = "webhook"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusDelivered Status = "delivered"
)

type Notification struct {
	ID        string    `json:"id"`
	Channel   Channel   `json:"channel"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject,omitempty"`
	Body      string    `json:"body"`
	Status    Status    `json:"status"`
	Error     string    `json:"error,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
}

func NewNotification(channel Channel, recipient, subject, body string) (*Notification, error) {
	if recipient == "" {
		return nil, errors.New("recipient is required")
	}
	if body == "" {
		return nil, errors.New("body is required")
	}
	if channel == ChannelEmail && subject == "" {
		return nil, errors.New("subject is required for email notifications")
	}
	if !isValidChannel(channel) {
		return nil, errors.New("invalid channel: must be email, sms, or webhook")
	}

	return &Notification{
		ID:        uuid.New().String(),
		Channel:   channel,
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
		Status:    StatusPending,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (n *Notification) MarkSent() {
	now := time.Now().UTC()
	n.Status = StatusSent
	n.SentAt = &now
}

func (n *Notification) MarkFailed(err string) {
	n.Status = StatusFailed
	n.Error = err
}

func (n *Notification) MarkDelivered() {
	n.Status = StatusDelivered
}

func isValidChannel(c Channel) bool {
	switch c {
	case ChannelEmail, ChannelSMS, ChannelWebhook:
		return true
	}
	return false
}
