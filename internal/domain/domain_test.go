package domain_test

import (
	"testing"

	"github.com/CarlosAlbertoFurtado/notifyhub/internal/domain"
)

func TestNewNotification_Valid(t *testing.T) {
	n, err := domain.NewNotification(domain.ChannelEmail, "user@example.com", "Hello", "Welcome!")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.ID == "" {
		t.Error("expected ID to be generated")
	}
	if n.Status != domain.StatusPending {
		t.Errorf("expected status pending, got %s", n.Status)
	}
}

func TestNewNotification_EmptyRecipient(t *testing.T) {
	_, err := domain.NewNotification(domain.ChannelEmail, "", "Hello", "body")
	if err == nil {
		t.Error("expected error for empty recipient")
	}
}

func TestNewNotification_EmptyBody(t *testing.T) {
	_, err := domain.NewNotification(domain.ChannelEmail, "user@example.com", "Hello", "")
	if err == nil {
		t.Error("expected error for empty body")
	}
}

func TestNewNotification_EmailRequiresSubject(t *testing.T) {
	_, err := domain.NewNotification(domain.ChannelEmail, "user@example.com", "", "body")
	if err == nil {
		t.Error("expected error for email without subject")
	}
}

func TestNewNotification_SMSNoSubjectRequired(t *testing.T) {
	n, err := domain.NewNotification(domain.ChannelSMS, "+5519999999999", "", "Your code is 1234")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Channel != domain.ChannelSMS {
		t.Errorf("expected channel sms, got %s", n.Channel)
	}
}

func TestNewNotification_InvalidChannel(t *testing.T) {
	_, err := domain.NewNotification(domain.Channel("pigeon"), "addr", "s", "b")
	if err == nil {
		t.Error("expected error for invalid channel")
	}
}

func TestNotification_MarkSent(t *testing.T) {
	n, _ := domain.NewNotification(domain.ChannelWebhook, "https://hook.site/abc", "", "payload")
	n.MarkSent()
	if n.Status != domain.StatusSent {
		t.Errorf("expected status sent, got %s", n.Status)
	}
	if n.SentAt == nil {
		t.Error("expected SentAt to be set")
	}
}

func TestNotification_MarkFailed(t *testing.T) {
	n, _ := domain.NewNotification(domain.ChannelEmail, "a@b.com", "s", "b")
	n.MarkFailed("timeout")
	if n.Status != domain.StatusFailed {
		t.Errorf("expected status failed, got %s", n.Status)
	}
	if n.Error != "timeout" {
		t.Errorf("expected error 'timeout', got %s", n.Error)
	}
}

func TestNewTemplate_Valid(t *testing.T) {
	tmpl, err := domain.NewTemplate("welcome", domain.ChannelEmail, "Welcome!", "<h1>Hi</h1>")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tmpl.Name != "welcome" {
		t.Errorf("expected name 'welcome', got %s", tmpl.Name)
	}
}

func TestNewTemplate_ShortName(t *testing.T) {
	_, err := domain.NewTemplate("a", domain.ChannelEmail, "s", "b")
	if err == nil {
		t.Error("expected error for short name")
	}
}

func TestNewTemplate_EmptyBody(t *testing.T) {
	_, err := domain.NewTemplate("valid-name", domain.ChannelEmail, "s", "")
	if err == nil {
		t.Error("expected error for empty body")
	}
}

func TestNotification_MarkFailedTwice(t *testing.T) {
	// segunda falha sobrescreve a mensagem de erro anterior
	n, _ := domain.NewNotification(domain.ChannelEmail, "a@b.com", "s", "b")
	n.MarkFailed("first error")
	n.MarkFailed("second error")
	if n.Error != "second error" {
		t.Errorf("expected 'second error', got %s", n.Error)
	}
}

func TestNotification_SentAtNilBeforeSend(t *testing.T) {
	n, _ := domain.NewNotification(domain.ChannelSMS, "+5511999990000", "", "otp: 4832")
	if n.SentAt != nil {
		t.Error("SentAt should be nil before sending")
	}
}
