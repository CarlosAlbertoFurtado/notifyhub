package infra

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"

	"github.com/CarlosAlbertoFurtado/notifyhub/internal/domain"
)

// EmailSender sends notifications via SMTP.
type EmailSender struct {
	host string
	port int
	user string
	pass string
	from string
}

func NewEmailSender(host string, port int, user, pass, from string) *EmailSender {
	return &EmailSender{host: host, port: port, user: user, pass: pass, from: from}
}

func (s *EmailSender) Send(_ context.Context, n *domain.Notification) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.from, n.Recipient, n.Subject, n.Body,
	)

	var auth smtp.Auth
	if s.user != "" {
		auth = smtp.PlainAuth("", s.user, s.pass, s.host)
	}

	if err := smtp.SendMail(addr, auth, s.from, []string{n.Recipient}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp: %w", err)
	}

	slog.Info("email sent", "to", n.Recipient, "subject", n.Subject)
	return nil
}

// LogSender is a fallback sender that just logs (useful for dev/test).
type LogSender struct{}

func (s *LogSender) Send(_ context.Context, n *domain.Notification) error {
	slog.Info("notification dispatched",
		"channel", n.Channel,
		"recipient", n.Recipient,
		"subject", n.Subject,
	)
	return nil
}
