package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Channel   Channel   `json:"channel"`
	Subject   string    `json:"subject,omitempty"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTemplate(name string, channel Channel, subject, body string) (*Template, error) {
	if name == "" || len(name) < 2 {
		return nil, errors.New("template name must be at least 2 characters")
	}
	if body == "" {
		return nil, errors.New("template body is required")
	}

	now := time.Now().UTC()
	return &Template{
		ID:        uuid.New().String(),
		Name:      name,
		Channel:   channel,
		Subject:   subject,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
