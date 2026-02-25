package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CarlosAlbertoFurtado/notifyhub/internal/application"
	"github.com/CarlosAlbertoFurtado/notifyhub/internal/config"
	"github.com/CarlosAlbertoFurtado/notifyhub/internal/domain"
	"github.com/CarlosAlbertoFurtado/notifyhub/internal/handler"
	"github.com/CarlosAlbertoFurtado/notifyhub/internal/infra"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("starting notifyhub", "port", cfg.Port)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(pool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	notifRepo := infra.NewPostgresNotificationRepo(pool)

	// LogSender por padrão; troca pra SMTP real se as credenciais estiverem configuradas
	senders := map[domain.Channel]domain.Sender{
		domain.ChannelEmail:   &infra.LogSender{},
		domain.ChannelSMS:     &infra.LogSender{},
		domain.ChannelWebhook: &infra.LogSender{},
	}

	if cfg.SMTPUser != "" {
		senders[domain.ChannelEmail] = infra.NewEmailSender(
			cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom,
		)
	}

	sendUC := application.NewSendNotificationUseCase(notifRepo, senders)

	nh := handler.NewNotificationHandler(sendUC, notifRepo)
	router := handler.SetupRouter(nh)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server listening", "addr", cfg.ListenAddr())
		if err := router.Run(cfg.ListenAddr()); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	slog.Info("shutting down...")
	time.Sleep(1 * time.Second)
}

func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS notifications (
			id         VARCHAR(36) PRIMARY KEY,
			channel    VARCHAR(20) NOT NULL,
			recipient  VARCHAR(255) NOT NULL,
			subject    VARCHAR(255),
			body       TEXT NOT NULL,
			status     VARCHAR(20) NOT NULL DEFAULT 'pending',
			error      TEXT,
			metadata   JSONB,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			sent_at    TIMESTAMPTZ
		);

		CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
		CREATE INDEX IF NOT EXISTS idx_notifications_channel ON notifications(channel);
		CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at DESC);
	`)
	if err != nil {
		return err
	}
	slog.Info("migrations applied")
	return nil
}
