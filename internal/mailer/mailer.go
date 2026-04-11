package mailer

import (
	"fmt"
	"log"
	"os"

	"github.com/wneessen/go-mail"
)

type Mailer struct {
	dialer *mail.Client
	sender string
}

func New() (*Mailer, error) {
	dialer, err := mail.NewClient(
		os.Getenv("SMTP_HOST"),
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(os.Getenv("SMTP_USER")),
		mail.WithPassword(os.Getenv("SMTP_PASS")),
	)

	if err != nil {
		return nil, err
	}

	return &Mailer{
		dialer: dialer,
		sender: os.Getenv("SMTP_USER"),
	}, nil
}

func (m *Mailer) SendOtp(recipient, otp string) error {
	plain := fmt.Sprintf(`
Your OTP is: %s
This code is valid for 5 minutes.
chitchat
`, otp)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body style="font-family: Arial, sans-serif; line-height:1.5;">
		<p>Your OTP is:</p>
		<h2 style="letter-spacing:2px;">%s</h2>
		<p>This code is valid for 5 minutes.</p>
		<p>chitchat</p>
		</body>
		</html>
`, otp)

	return m.sendEmail(recipient, "Your OTP", plain, html)
}

func (m *Mailer) SendMagicLink(recipient, link string) error {
	plain := fmt.Sprintf(`
Click the link below to sign in to chitchat:
%s

This link will expire in 15 minutes.
If you didn't request this, please ignore this email.

chitchat
`, link)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				.container { font-family: Arial, sans-serif; line-height: 1.5; max-width: 500px; margin: 0 auto; padding: 20px; }
				.button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
				.footer { color: #666; font-size: 12px; margin-top: 30px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Sign in to chitchat</h2>
				<p>Click the button below to sign in. This link will expire in 15 minutes.</p>
				<a href="%s" class="button">Sign In</a>
				<div class="footer">
					<p>If you didn't request this email, please ignore it.</p>
					<p>chitchat</p>
				</div>
			</div>
		</body>
		</html>
`, link)

	return m.sendEmail(recipient, "Sign in to chitchat", plain, html)
}

func (m *Mailer) sendEmail(recipient, subject, plain, html string) error {
	message := mail.NewMsg()
	if err := message.From(os.Getenv("SMTP_FROM")); err != nil {
		return err
	}

	if err := message.To(recipient); err != nil {
		return err
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, plain)
	message.AddAlternativeString(mail.TypeTextHTML, html)

	go func(m *Mailer, message *mail.Msg) {
		if err := m.dialer.DialAndSend(message); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}(m, message)

	return nil
}
