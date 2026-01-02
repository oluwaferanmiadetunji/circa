package email_test

import (
	"circa/internal/email"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/resend/resend-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockResendClient struct {
	sendFunc func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error)
}

func (m *mockResendClient) Emails() email.ResendEmailsService {
	return &mockEmailsService{sendFunc: m.sendFunc}
}

type mockEmailsService struct {
	sendFunc func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error)
}

func (m *mockEmailsService) SendWithContext(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	return m.sendFunc(ctx, params)
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		fromEmail string
		fromName  string
	}{
		{
			name:      "success - creates service with all fields",
			apiKey:    "test-api-key",
			fromEmail: "test@example.com",
			fromName:  "Test Name",
		},
		{
			name:      "success - empty api key",
			apiKey:    "",
			fromEmail: "test@example.com",
			fromName:  "Test Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := email.NewService(tt.apiKey, tt.fromEmail, tt.fromName)
			require.NotNil(t, service)
		})
	}
}

func TestService_SendMagicLink(t *testing.T) {
	tests := []struct {
		name          string
		toEmail       string
		toName        string
		magicLinkURL  string
		fromEmail     string
		fromName      string
		setupMock     func() *mockResendClient
		expectedError error
		validateEmail func(*testing.T, *resend.SendEmailRequest)
	}{
		{
			name:         "success - sends email with correct parameters",
			toEmail:      "user@example.com",
			toName:       "John Doe",
			magicLinkURL: "https://example.com/verify?token=abc123",
			fromEmail:    "noreply@circa.com",
			fromName:     "Circa",
			setupMock: func() *mockResendClient {
				return &mockResendClient{
					sendFunc: func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
						return &resend.SendEmailResponse{Id: "test-id"}, nil
					},
				}
			},
			expectedError: nil,
			validateEmail: func(t *testing.T, params *resend.SendEmailRequest) {
				assert.Equal(t, "Circa <noreply@circa.com>", params.From)
				assert.Equal(t, []string{"user@example.com"}, params.To)
				assert.Equal(t, "Verify your email for Circa", params.Subject)
				assert.Contains(t, params.Html, "John Doe")
				assert.Contains(t, params.Html, "https://example.com/verify?token=abc123")
				assert.Contains(t, params.Text, "John Doe")
				assert.Contains(t, params.Text, "https://example.com/verify?token=abc123")
			},
		},
		{
			name:         "success - email contains current year",
			toEmail:      "user@example.com",
			toName:       "Jane Smith",
			magicLinkURL: "https://example.com/verify?token=xyz789",
			fromEmail:    "noreply@circa.com",
			fromName:     "Circa",
			setupMock: func() *mockResendClient {
				return &mockResendClient{
					sendFunc: func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
						return &resend.SendEmailResponse{Id: "test-id"}, nil
					},
				}
			},
			expectedError: nil,
			validateEmail: func(t *testing.T, params *resend.SendEmailRequest) {
				currentYear := time.Now().Year()
				assert.Contains(t, params.Html, fmt.Sprintf("%d", currentYear))
			},
		},
		{
			name:         "error - resend client returns error",
			toEmail:      "user@example.com",
			toName:       "John Doe",
			magicLinkURL: "https://example.com/verify?token=abc123",
			fromEmail:    "noreply@circa.com",
			fromName:     "Circa",
			setupMock: func() *mockResendClient {
				return &mockResendClient{
					sendFunc: func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
						return nil, errors.New("resend API error")
					},
				}
			},
			expectedError: errors.New("resend API error"),
			validateEmail: nil,
		},
		{
			name:         "success - email contains magic link URL in both HTML and text",
			toEmail:      "user@example.com",
			toName:       "Test User",
			magicLinkURL: "https://circa.com/auth/verify?token=special-token-123",
			fromEmail:    "support@circa.com",
			fromName:     "Circa Support",
			setupMock: func() *mockResendClient {
				return &mockResendClient{
					sendFunc: func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
						return &resend.SendEmailResponse{Id: "test-id"}, nil
					},
				}
			},
			expectedError: nil,
			validateEmail: func(t *testing.T, params *resend.SendEmailRequest) {
				assert.Contains(t, params.Html, "https://circa.com/auth/verify?token=special-token-123")
				assert.Contains(t, params.Text, "https://circa.com/auth/verify?token=special-token-123")
				assert.Contains(t, params.Html, "Test User")
				assert.Contains(t, params.Text, "Test User")
			},
		},
		{
			name:         "success - email contains proper HTML structure",
			toEmail:      "user@example.com",
			toName:       "John Doe",
			magicLinkURL: "https://example.com/verify?token=abc123",
			fromEmail:    "noreply@circa.com",
			fromName:     "Circa",
			setupMock: func() *mockResendClient {
				return &mockResendClient{
					sendFunc: func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
						return &resend.SendEmailResponse{Id: "test-id"}, nil
					},
				}
			},
			expectedError: nil,
			validateEmail: func(t *testing.T, params *resend.SendEmailRequest) {
				assert.Contains(t, params.Html, "<!DOCTYPE html>")
				assert.Contains(t, params.Html, "<html>")
				assert.Contains(t, params.Html, "<body")
				assert.Contains(t, params.Html, "Verify Email")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedParams *resend.SendEmailRequest
			mockClient := &mockResendClient{
				sendFunc: func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
					capturedParams = params
					return tt.setupMock().sendFunc(ctx, params)
				},
			}

			service := email.NewService("test-api-key", tt.fromEmail, tt.fromName)
			originalClient := service.Client()
			service.SetClient(mockClient)

			err := service.SendMagicLink(context.Background(), tt.toEmail, tt.toName, tt.magicLinkURL)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				if tt.validateEmail != nil {
					tt.validateEmail(t, capturedParams)
				}
			}

			service.SetClient(originalClient)
		})
	}
}

func TestService_SendMagicLink_EmailContent(t *testing.T) {
	service := email.NewService("test-api-key", "noreply@circa.com", "Circa")
	
	var capturedParams *resend.SendEmailRequest
	mockClient := &mockResendClient{
		sendFunc: func(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
			capturedParams = params
			return &resend.SendEmailResponse{Id: "test-id"}, nil
		},
	}

	originalClient := service.Client()
	service.SetClient(mockClient)

	err := service.SendMagicLink(context.Background(), "user@example.com", "John Doe", "https://example.com/verify?token=abc123")
	require.NoError(t, err)

	html := capturedParams.Html
	text := capturedParams.Text

	assert.Contains(t, html, "Welcome to Circa!")
	assert.Contains(t, html, "John Doe")
	assert.Contains(t, html, "https://example.com/verify?token=abc123")
	assert.Contains(t, html, "24 hours")
	assert.Contains(t, html, "Verify Email")

	assert.Contains(t, text, "John Doe")
	assert.Contains(t, text, "https://example.com/verify?token=abc123")
	assert.Contains(t, text, "24 hours")

	assert.True(t, strings.Contains(html, fmt.Sprintf("%d", time.Now().Year())))

	service.SetClient(originalClient)
}
