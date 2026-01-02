package queue_test

import (
	dbmocks "circa/internal/db/mocks"
	"circa/internal/queue"
	"circa/internal/queue/mocks"
	sqlc "circa/internal/db/sqlc/generated"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

func TestWorker_ProcessJobs(t *testing.T) {
	tests := []struct {
		name               string
		setupMocks         func(*dbmocks.MockStore, *mocks.MockEmailService)
		expectedCalls      func(*testing.T, *dbmocks.MockStore, *mocks.MockEmailService)
		useNilEmailService bool
	}{
		{
			name: "success - no pending jobs",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.On("GetNextPendingJob", mock.Anything).
					Return(sqlc.Job{}, pgx.ErrNoRows).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
			},
		},
		{
			name: "success - process send_magic_link_email job",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("send_magic_link_email", map[string]interface{}{
					"email":         "test@example.com",
					"name":          "Test User",
					"magic_link_url": "https://example.com/verify?token=abc123",
				})
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				es.On("SendMagicLink", mock.Anything, "test@example.com", "Test User", "https://example.com/verify?token=abc123").
					Return(nil).Once()
				ms.On("UpdateJobStatus", mock.Anything, mock.MatchedBy(func(params sqlc.UpdateJobStatusParams) bool {
					return params.Status == "completed"
				})).Return(job, nil).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
				es.AssertExpectations(t)
			},
		},
		{
			name: "error - unknown job type",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("unknown_job", map[string]interface{}{})
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				ms.On("UpdateJobStatus", mock.Anything, mock.MatchedBy(func(params sqlc.UpdateJobStatusParams) bool {
					return params.Status == "failed" &&
						params.ErrorMessage != nil &&
						*params.ErrorMessage == "Unknown job type"
				})).Return(job, nil).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
			},
		},
		{
			name: "error - invalid payload",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("send_magic_link_email", map[string]interface{}{})
				job.Payload = []byte("invalid json")
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				ms.On("UpdateJobStatus", mock.Anything, mock.MatchedBy(func(params sqlc.UpdateJobStatusParams) bool {
					return params.Status == "failed" &&
						params.ErrorMessage != nil &&
						*params.ErrorMessage == "Invalid payload format"
				})).Return(job, nil).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
			},
		},
		{
			name: "error - email service nil",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("send_magic_link_email", map[string]interface{}{
					"email":         "test@example.com",
					"name":          "Test User",
					"magic_link_url": "https://example.com/verify?token=abc123",
				})
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				ms.On("UpdateJobStatus", mock.Anything, mock.MatchedBy(func(params sqlc.UpdateJobStatusParams) bool {
					return params.Status == "failed" &&
						params.ErrorMessage != nil &&
						*params.ErrorMessage == "Email service not configured"
				})).Return(job, nil).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
			},
			useNilEmailService: true,
		},
		{
			name: "error - email send fails, retry available",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("send_magic_link_email", map[string]interface{}{
					"email":         "test@example.com",
					"name":          "Test User",
					"magic_link_url": "https://example.com/verify?token=abc123",
				})
				job.RetryCount = 0
				job.MaxRetries = 3
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				es.On("SendMagicLink", mock.Anything, "test@example.com", "Test User", "https://example.com/verify?token=abc123").
					Return(errors.New("email service error")).Once()
				ms.On("IncrementJobRetry", mock.Anything, mock.MatchedBy(func(params sqlc.IncrementJobRetryParams) bool {
					return params.ErrorMessage != nil &&
						*params.ErrorMessage == "email service error"
				})).Return(job, nil).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
				es.AssertExpectations(t)
			},
		},
		{
			name: "error - email send fails, max retries exceeded",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("send_magic_link_email", map[string]interface{}{
					"email":         "test@example.com",
					"name":          "Test User",
					"magic_link_url": "https://example.com/verify?token=abc123",
				})
				job.RetryCount = 3
				job.MaxRetries = 3
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				es.On("SendMagicLink", mock.Anything, "test@example.com", "Test User", "https://example.com/verify?token=abc123").
					Return(errors.New("email service error")).Once()
				ms.On("UpdateJobStatus", mock.Anything, mock.MatchedBy(func(params sqlc.UpdateJobStatusParams) bool {
					return params.Status == "failed" &&
						params.ErrorMessage != nil &&
						*params.ErrorMessage == "email service error"
				})).Return(job, nil).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
				es.AssertExpectations(t)
			},
		},
		{
			name: "error - mark completed fails",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("send_magic_link_email", map[string]interface{}{
					"email":         "test@example.com",
					"name":          "Test User",
					"magic_link_url": "https://example.com/verify?token=abc123",
				})
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				es.On("SendMagicLink", mock.Anything, "test@example.com", "Test User", "https://example.com/verify?token=abc123").
					Return(nil).Once()
				ms.On("UpdateJobStatus", mock.Anything, mock.Anything).
					Return(sqlc.Job{}, errors.New("database error")).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
				es.AssertExpectations(t)
			},
		},
		{
			name: "error - retry job fails",
			setupMocks: func(ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				job := createTestJobWithType("send_magic_link_email", map[string]interface{}{
					"email":         "test@example.com",
					"name":          "Test User",
					"magic_link_url": "https://example.com/verify?token=abc123",
				})
				job.RetryCount = 1
				job.MaxRetries = 3
				ms.On("GetNextPendingJob", mock.Anything).Return(job, nil).Once()
				es.On("SendMagicLink", mock.Anything, "test@example.com", "Test User", "https://example.com/verify?token=abc123").
					Return(errors.New("email service error")).Once()
				ms.On("IncrementJobRetry", mock.Anything, mock.Anything).
					Return(sqlc.Job{}, errors.New("database error")).Once()
			},
			expectedCalls: func(t *testing.T, ms *dbmocks.MockStore, es *mocks.MockEmailService) {
				ms.AssertExpectations(t)
				es.AssertExpectations(t)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := dbmocks.NewMockStore(t)
			var mockEmailService *mocks.MockEmailService
			if !tt.useNilEmailService {
				mockEmailService = mocks.NewMockEmailService(t)
			}
			tt.setupMocks(mockStore, mockEmailService)

			queueService := queue.NewService(mockStore)
			var worker *queue.Worker
			if tt.useNilEmailService {
				worker = queue.NewWorker(queueService, nil)
			} else {
				worker = queue.NewWorker(queueService, mockEmailService)
			}

			ctx := context.Background()
			worker.ProcessJobs(ctx)

			tt.expectedCalls(t, mockStore, mockEmailService)
		})
	}
}

func TestWorker_Start(t *testing.T) {
	mockStore := dbmocks.NewMockStore(t)
	mockEmailService := mocks.NewMockEmailService(t)

	queueService := queue.NewService(mockStore)
	worker := queue.NewWorker(queueService, mockEmailService)

	ctx, cancel := context.WithCancel(context.Background())

	mockStore.On("GetNextPendingJob", mock.Anything).
		Return(sqlc.Job{}, pgx.ErrNoRows).
		Maybe()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	worker.Start(ctx)

	mockStore.AssertExpectations(t)
}

func TestWorker_Stop(t *testing.T) {
	mockStore := dbmocks.NewMockStore(t)
	mockEmailService := mocks.NewMockEmailService(t)

	queueService := queue.NewService(mockStore)
	worker := queue.NewWorker(queueService, mockEmailService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockStore.On("GetNextPendingJob", mock.Anything).
		Return(sqlc.Job{}, pgx.ErrNoRows).
		Maybe()

	done := make(chan bool)
	go func() {
		worker.Start(ctx)
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	worker.Stop()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Worker did not stop")
	}
}

func createTestJobWithType(jobType string, payload map[string]interface{}) sqlc.Job {
	now := time.Now()
	payloadJSON, _ := json.Marshal(payload)
	return sqlc.Job{
		ID:          uuid.New(),
		Type:        jobType,
		Payload:     payloadJSON,
		Status:      "pending",
		RetryCount:  0,
		MaxRetries:  3,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

