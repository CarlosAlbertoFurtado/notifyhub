package domain

import "context"

// NotificationRepository persists and queries notifications.
type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	FindAll(ctx context.Context, params ListParams) (*PaginatedResult, error)
	Update(ctx context.Context, n *Notification) error
	Stats(ctx context.Context) (*NotificationStats, error)
}

// TemplateRepository persists notification templates.
type TemplateRepository interface {
	Create(ctx context.Context, t *Template) error
	FindByID(ctx context.Context, id string) (*Template, error)
	FindByName(ctx context.Context, name string) (*Template, error)
	FindAll(ctx context.Context) ([]*Template, error)
	Update(ctx context.Context, t *Template) error
	Delete(ctx context.Context, id string) error
}

// Sender dispatches a notification through a specific channel.
type Sender interface {
	Send(ctx context.Context, n *Notification) error
}

type ListParams struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Channel string `json:"channel,omitempty"`
	Status  string `json:"status,omitempty"`
}

type PaginatedResult struct {
	Data       []*Notification `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

type NotificationStats struct {
	Total     int64 `json:"total"`
	Sent      int64 `json:"sent"`
	Failed    int64 `json:"failed"`
	Pending   int64 `json:"pending"`
	Delivered int64 `json:"delivered"`
}
