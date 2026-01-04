package email

import "context"

type EmailService interface {
	SendMagicLink(ctx context.Context, toEmail, toName, magicLinkURL string, isLogin bool) error
}
