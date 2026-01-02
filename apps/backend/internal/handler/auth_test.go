package handler

import (
	"bytes"
	"circa/api"
	sqlc "circa/internal/db/sqlc/generated"
	circaerrors "circa/internal/errors"
	"circa/internal/service/auth"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) CreatePendingSignup(ctx context.Context, fullName, email string, displayName *string) (*auth.SignupResult, error) {
	args := m.Called(ctx, fullName, email, displayName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.SignupResult), args.Error(1)
}

func (m *mockAuthService) GenerateNonce(address string) (*auth.NonceResult, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.NonceResult), args.Error(1)
}

func TestHandler_AuthSignup(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mockAuthService)
		expectedStatus int
		expectedBody   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - valid request with display name",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"email":        "john@example.com",
				"display_name": "johndoe",
			},
			setupMocks: func(m *mockAuthService) {
				result := &auth.SignupResult{
					PendingSignup: createTestPendingSignup(),
					MagicLink:     createTestMagicLink(),
				}
				m.On("CreatePendingSignup", mock.Anything, "John Doe", "john@example.com", mock.MatchedBy(func(d *string) bool {
					return d != nil && *d == "johndoe"
				})).Return(result, nil)
			},
			expectedStatus: 200,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.AuthSignupResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Signup request created successfully", response.Message)
			},
		},
		{
			name: "success - valid request without display name",
			requestBody: map[string]interface{}{
				"full_name": "Jane Smith",
				"email":     "jane@example.com",
			},
			setupMocks: func(m *mockAuthService) {
				result := &auth.SignupResult{
					PendingSignup: createTestPendingSignup(),
					MagicLink:     createTestMagicLink(),
				}
				m.On("CreatePendingSignup", mock.Anything, "Jane Smith", "jane@example.com", (*string)(nil)).
					Return(result, nil)
			},
			expectedStatus: 200,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.AuthSignupResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Signup request created successfully", response.Message)
			},
		},
		{
			name: "error - invalid JSON body",
			requestBody: func() interface{} {
				return "invalid json"
			}(),
			setupMocks:     func(m *mockAuthService) {},
			expectedStatus: 400,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.ErrorBadRequest
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, 400, response.Code)
				assert.Equal(t, "Invalid request body", response.Message)
			},
		},
		{
			name: "error - missing full name",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			setupMocks:     func(m *mockAuthService) {},
			expectedStatus: 400,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.ErrorBadRequest
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, 400, response.Code)
				assert.Equal(t, "Full name is required", response.Message)
			},
		},
		{
			name: "error - empty full name",
			requestBody: map[string]interface{}{
				"full_name": "",
				"email":     "test@example.com",
			},
			setupMocks:     func(m *mockAuthService) {},
			expectedStatus: 400,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.ErrorBadRequest
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, 400, response.Code)
				assert.Equal(t, "Full name is required", response.Message)
			},
		},
		{
			name: "error - missing email",
			requestBody: map[string]interface{}{
				"full_name": "Test User",
			},
			setupMocks:     func(m *mockAuthService) {},
			expectedStatus: 400,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.ErrorBadRequest
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, 400, response.Code)
				assert.Equal(t, "Email is required", response.Message)
			},
		},
		{
			name: "error - empty email",
			requestBody: map[string]interface{}{
				"full_name": "Test User",
				"email":     "",
			},
			setupMocks:     func(m *mockAuthService) {},
			expectedStatus: 400,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.ErrorBadRequest
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, 400, response.Code)
				assert.Equal(t, "Invalid request body", response.Message)
			},
		},
		{
			name: "error - email already exists",
			requestBody: map[string]interface{}{
				"full_name": "John Doe",
				"email":     "existing@example.com",
			},
			setupMocks: func(m *mockAuthService) {
				m.On("CreatePendingSignup", mock.Anything, "John Doe", "existing@example.com", (*string)(nil)).
					Return(nil, circaerrors.ErrEmailAlreadyExists)
			},
			expectedStatus: 400,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.ErrorBadRequest
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, 400, response.Code)
				assert.Equal(t, "Email already exists", response.Message)
			},
		},
		{
			name: "error - service returns generic error",
			requestBody: map[string]interface{}{
				"full_name": "John Doe",
				"email":     "test@example.com",
			},
			setupMocks: func(m *mockAuthService) {
				m.On("CreatePendingSignup", mock.Anything, "John Doe", "test@example.com", (*string)(nil)).
					Return(nil, errors.New("database connection error"))
			},
			expectedStatus: 500,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.ErrorInternalServerError
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, 500, response.Code)
				assert.Equal(t, "Internal server error", response.Message)
			},
		},
		{
			name: "success - empty display name string",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"email":        "john@example.com",
				"display_name": "",
			},
			setupMocks: func(m *mockAuthService) {
				result := &auth.SignupResult{
					PendingSignup: createTestPendingSignup(),
					MagicLink:     createTestMagicLink(),
				}
				emptyStr := ""
				m.On("CreatePendingSignup", mock.Anything, "John Doe", "john@example.com", &emptyStr).
					Return(result, nil)
			},
			expectedStatus: 200,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response api.AuthSignupResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Signup request created successfully", response.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			var reqBody []byte
			var err error
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockAuth := &mockAuthService{}
			mockAuth.Test(t)
			tt.setupMocks(mockAuth)

			handler := &Handler{
				authService: mockAuth,
			}

			err = handler.AuthSignup(c)

			if tt.expectedStatus >= 400 {
				if err != nil {
					httpErr, ok := err.(*echo.HTTPError)
					if ok {
						assert.Equal(t, tt.expectedStatus, httpErr.Code)
					}
				}
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != nil {
				tt.expectedBody(t, rec)
			}

			mockAuth.AssertExpectations(t)
		})
	}
}

func createTestPendingSignup() sqlc.PendingSignup {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	return sqlc.PendingSignup{
		ID:          uuid.New(),
		FullName:    pgtype.Text{String: "John Doe", Valid: true},
		Email:       pgtype.Text{String: "john@example.com", Valid: true},
		DisplayName: stringPtr("johndoe"),
		Status:      "pending",
		ExpiresAt:   pgtype.Timestamp{Time: expiresAt, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: now, Valid: true},
	}
}

func createTestMagicLink() sqlc.MagicLink {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	return sqlc.MagicLink{
		ID:              uuid.New(),
		PendingSignupID: uuid.New(),
		TokenHash:       "test_token_hash",
		ExpiresAt:       pgtype.Timestamp{Time: expiresAt, Valid: true},
		CreatedAt:       pgtype.Timestamp{Time: now, Valid: true},
		UpdatedAt:       pgtype.Timestamp{Time: now, Valid: true},
	}
}

func stringPtr(s string) *string {
	return &s
}
