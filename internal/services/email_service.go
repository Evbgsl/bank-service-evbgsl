package services

import (
	"crypto/tls"
	"fmt"
	"html"
	"net/mail"

	"github.com/evbgsl/bank-service-evbgsl/internal/config"
	gomail "gopkg.in/gomail.v2"
)

type EmailService struct {
	cfg config.SMTPConfig
}

func NewEmailService(cfg config.SMTPConfig) *EmailService {
	return &EmailService{
		cfg: cfg,
	}
}

func (s *EmailService) SendTestEmail(to string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return ErrInvalidEmail
	}

	subject := "Bank Service test email"
	body := `
		<h1>Bank Service</h1>
		<p>SMTP integration works correctly.</p>
	`

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendCreditPaymentEmail(
	to string,
	amount float64,
	status string,
) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return ErrInvalidEmail
	}

	subject := "Credit payment notification"

	body := fmt.Sprintf(`
		<h1>Credit payment notification</h1>
		<p>Status: <strong>%s</strong></p>
		<p>Amount: <strong>%.2f RUB</strong></p>
		<small>This is an automatic notification from Bank Service.</small>
	`, html.EscapeString(status), amount)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendEmail(to string, subject string, body string) error {
	if !s.cfg.Enabled {
		return nil
	}

	if s.cfg.Host == "" || s.cfg.User == "" || s.cfg.Password == "" || s.cfg.From == "" {
		return fmt.Errorf("smtp is enabled but configuration is incomplete")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.cfg.From)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		s.cfg.Host,
		s.cfg.Port,
		s.cfg.User,
		s.cfg.Password,
	)

	dialer.TLSConfig = &tls.Config{
		ServerName:         s.cfg.Host,
		InsecureSkipVerify: false,
	}

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}

	return nil
}
