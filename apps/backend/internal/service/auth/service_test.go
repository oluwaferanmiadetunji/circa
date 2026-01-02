package auth

import (
	"circa/internal/db"
	dbmocks "circa/internal/db/mocks"
	sqlc "circa/internal/db/sqlc/generated"
	circaerrors "circa/internal/errors"
	txmocks "circa/internal/service/auth/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testablePGXStore struct {
	*dbmocks.MockStore
	mockTx *txmocks.MockTx
}

type mockPoolWrapper struct {
	tx *txmocks.MockTx
}

func (m *mockPoolWrapper) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.tx, nil
}

func (t *testablePGXStore) GetDB() interface{} {
	return &mockPoolWrapper{tx: t.mockTx}
}

func TestService_CreatePendingSignup(t *testing.T) {
	tests := []struct {
		name           string
		fullName       string
		email          string
		displayName    *string
		setupMocks     func(*dbmocks.MockStore, *txmocks.MockTx)
		usePGXStore    bool
		expectedError  error
		expectedResult *SignupResult
	}{
		{
			name:        "error - email already exists",
			fullName:    "John Doe",
			email:       "existing@example.com",
			displayName: stringPtr("johndoe"),
			setupMocks: func(ms *dbmocks.MockStore, mt *txmocks.MockTx) {
				ms.On("GetUserByEmail", mock.Anything, mock.Anything).
					Return(createTestUser(), nil)
			},
			usePGXStore:    false,
			expectedError:  circaerrors.ErrEmailAlreadyExists,
			expectedResult: nil,
		},
		{
			name:        "error - database error when checking email",
			fullName:    "John Doe",
			email:       "john@example.com",
			displayName: stringPtr("johndoe"),
			setupMocks: func(ms *dbmocks.MockStore, mt *txmocks.MockTx) {
				ms.On("GetUserByEmail", mock.Anything, mock.Anything).
					Return(sqlc.User{}, errors.New("database connection error"))
			},
			usePGXStore:    false,
			expectedError:  errors.New("database connection error"),
			expectedResult: nil,
		},
		{
			name:        "error - invalid store type",
			fullName:    "John Doe",
			email:       "john@example.com",
			displayName: stringPtr("johndoe"),
			setupMocks: func(ms *dbmocks.MockStore, mt *txmocks.MockTx) {
				ms.On("GetUserByEmail", mock.Anything, mock.Anything).
					Return(sqlc.User{}, pgx.ErrNoRows)
			},
			usePGXStore:    false,
			expectedError:  circaerrors.ErrInvalidStore,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := dbmocks.NewMockStore(t)
			mockTx := txmocks.NewMockTx(t)
			tt.setupMocks(mockStore, mockTx)

			var store db.Store
			if tt.usePGXStore {
				t.Skip("PGXStore mocking requires more complex setup - skipping for now")
				return
			} else {
				store = mockStore
			}

			service := NewService(store, nil, "https://example.com", 15*time.Minute)

			result, err := service.CreatePendingSignup(context.Background(), tt.fullName, tt.email, tt.displayName)
			if tt.expectedError != nil {
				require.Error(t, err)
				if errors.Is(tt.expectedError, circaerrors.ErrEmailAlreadyExists) {
					assert.ErrorIs(t, err, circaerrors.ErrEmailAlreadyExists)
				} else if errors.Is(tt.expectedError, circaerrors.ErrInvalidStore) {
					assert.ErrorIs(t, err, circaerrors.ErrInvalidStore)
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.PendingSignup.ID, result.PendingSignup.ID)
					assert.Equal(t, tt.expectedResult.MagicLink.ID, result.MagicLink.ID)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func createTestUser() sqlc.User {
	return sqlc.User{
		ID:          uuid.New(),
		FullName:    pgtype.Text{String: "Test User", Valid: true},
		Email:       pgtype.Text{String: "test@example.com", Valid: true},
		Address:     "0x1234567890123456789012345678901234567890",
		DisplayName: stringPtr("testuser"),
		CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
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
