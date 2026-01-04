package email

import (
	"context"
	"fmt"
	"time"

	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog/log"
)

type ResendClient interface {
	Emails() ResendEmailsService
}

type ResendEmailsService interface {
	SendWithContext(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error)
}

type Service struct {
	client ResendClient
}

func (s *Service) SetClient(client ResendClient) {
	s.client = client
}

func (s *Service) Client() ResendClient {
	return s.client
}

func NewService(apiKey string) *Service {
	client := resend.NewClient(apiKey)
	return &Service{
		client: &resendClientWrapper{client: client},
	}
}

type resendClientWrapper struct {
	client *resend.Client
}

func (r *resendClientWrapper) Emails() ResendEmailsService {
	return &resendEmailsServiceWrapper{emails: r.client.Emails}
}

type resendEmailsServiceWrapper struct {
	emails resend.EmailsSvc
}

func (r *resendEmailsServiceWrapper) SendWithContext(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	return r.emails.SendWithContext(ctx, params)
}

func (s *Service) SendMagicLink(ctx context.Context, toEmail, toName, magicLinkURL string, isLogin bool) error {
	var subject, headerText, bodyText, buttonText, footerText string

	if isLogin {
		subject = "Sign in to Circa"
		headerText = "Sign in to Circa"
		bodyText = "Click the button below to sign in to your account:"
		buttonText = "Sign In"
		footerText = "If you didn't request this login link, you can safely ignore this email."
	} else {
		subject = "Verify your email for Circa"
		headerText = "Welcome to Circa!"
		bodyText = "Thank you for signing up! Please verify your email address by clicking the button below:"
		buttonText = "Verify Email"
		footerText = "If you didn't create an account, you can safely ignore this email."
	}

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
		</head>
		<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
			<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center; border-radius: 8px 8px 0 0;">
				<h1 style="color: white; margin: 0; font-size: 28px;">%s</h1>
			</div>
			<div style="background: #ffffff; padding: 40px; border-radius: 0 0 8px 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
				<p style="font-size: 16px; margin-bottom: 20px;">Hi %s,</p>
				<p style="font-size: 16px; margin-bottom: 20px;">%s</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 14px 28px; text-decoration: none; border-radius: 6px; display: inline-block; font-weight: 600; font-size: 16px;">%s</a>
				</div>
				<p style="font-size: 14px; color: #666; margin-top: 30px;">Or copy and paste this link into your browser:</p>
				<p style="font-size: 12px; color: #999; word-break: break-all; background: #f5f5f5; padding: 10px; border-radius: 4px;">%s</p>
				<p style="font-size: 14px; color: #666; margin-top: 30px;">This link will expire in 24 hours.</p>
				<p style="font-size: 14px; color: #666; margin-top: 20px;">%s</p>
			</div>
			<div style="text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee;">
				<p style="font-size: 12px; color: #999;">© %d Circa. All rights reserved.</p>
			</div>
		</body>
		</html>
	`, headerText, toName, bodyText, magicLinkURL, buttonText, magicLinkURL, footerText, time.Now().Year())

	var textBody string
	if isLogin {
		textBody = fmt.Sprintf(`
Hi %s,

Click the link below to sign in to your account:

%s

This link will expire in 24 hours.

If you didn't request this login link, you can safely ignore this email.
		`, toName, magicLinkURL)
	} else {
		textBody = fmt.Sprintf(`
Hi %s,

Thank you for signing up! Please verify your email address by visiting this link:

%s

This link will expire in 24 hours.

If you didn't create an account, you can safely ignore this email.
		`, toName, magicLinkURL)
	}

	params := &resend.SendEmailRequest{
		From:    "Circa <onboarding@resend.dev>",
		To:      []string{toEmail},
		Subject: subject,
		Html:    htmlBody,
		Text:    textBody,
	}

	sent, err := s.client.Emails().SendWithContext(ctx, params)
	if err != nil {
		log.Error().Err(err).Str("email", toEmail).Msg("Failed to send magic link email")
		return err
	}

	log.Info().
		Str("email", toEmail).
		Str("resend_id", sent.Id).
		Msg("Magic link email sent successfully")

	return nil
}
