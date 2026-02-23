package infra

import (
	"context"
	"fmt"
	"math"

	"github.com/CarlosAlbertoFurtado/notifyhub/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresNotificationRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresNotificationRepo(pool *pgxpool.Pool) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{pool: pool}
}

func (r *PostgresNotificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO notifications (id, channel, recipient, subject, body, status, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		n.ID, n.Channel, n.Recipient, n.Subject, n.Body, n.Status, n.Metadata, n.CreatedAt,
	)
	return err
}

func (r *PostgresNotificationRepo) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	n := &domain.Notification{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, channel, recipient, subject, body, status, error, metadata, created_at, sent_at
		FROM notifications WHERE id = $1`, id,
	).Scan(&n.ID, &n.Channel, &n.Recipient, &n.Subject, &n.Body, &n.Status, &n.Error, &n.Metadata, &n.CreatedAt, &n.SentAt)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *PostgresNotificationRepo) FindAll(ctx context.Context, params domain.ListParams) (*domain.PaginatedResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	offset := (params.Page - 1) * params.Limit

	// count
	where, args := buildWhere(params)
	var total int64
	countQ := "SELECT COUNT(*) FROM notifications" + where
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, err
	}

	// data
	dataQ := fmt.Sprintf(
		"SELECT id, channel, recipient, subject, body, status, error, metadata, created_at, sent_at FROM notifications%s ORDER BY created_at DESC LIMIT %d OFFSET %d",
		where, params.Limit, offset,
	)
	rows, err := r.pool.Query(ctx, dataQ, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		if err := rows.Scan(&n.ID, &n.Channel, &n.Recipient, &n.Subject, &n.Body, &n.Status, &n.Error, &n.Metadata, &n.CreatedAt, &n.SentAt); err != nil {
			return nil, err
		}
		data = append(data, n)
	}

	return &domain.PaginatedResult{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}, nil
}

func (r *PostgresNotificationRepo) Update(ctx context.Context, n *domain.Notification) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE notifications SET status = $1, error = $2, sent_at = $3 WHERE id = $4`,
		n.Status, n.Error, n.SentAt, n.ID,
	)
	return err
}

func (r *PostgresNotificationRepo) Stats(ctx context.Context) (*domain.NotificationStats, error) {
	s := &domain.NotificationStats{}
	err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'sent'),
			COUNT(*) FILTER (WHERE status = 'failed'),
			COUNT(*) FILTER (WHERE status = 'pending'),
			COUNT(*) FILTER (WHERE status = 'delivered')
		FROM notifications`,
	).Scan(&s.Total, &s.Sent, &s.Failed, &s.Pending, &s.Delivered)
	return s, err
}

func buildWhere(p domain.ListParams) (string, []any) {
	where := ""
	var args []any
	i := 1

	if p.Channel != "" {
		where += fmt.Sprintf(" WHERE channel = $%d", i)
		args = append(args, p.Channel)
		i++
	}
	if p.Status != "" {
		if where == "" {
			where += " WHERE "
		} else {
			where += " AND "
		}
		where += fmt.Sprintf("status = $%d", i)
		args = append(args, p.Status)
	}
	return where, args
}
