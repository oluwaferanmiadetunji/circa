package email

import "context"

type EmailService interface {
	SendMagicLink(ctx context.Context, toEmail, toName, magicLinkURL string) error
}

