package queue_test

import (
	dbmocks "circa/internal/db/mocks"
	"circa/internal/queue"
	sqlc "circa/internal/db/sqlc/generated"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Enqueue(t *testing.T) {
	tests := []struct {
		name          string
		jobType       string
		payload       queue.JobPayload
		maxRetries    *int
		setupMocks    func(*dbmocks.MockStore)
		expectedError error
		validateJob   func(*testing.T, *sqlc.Job)
	}{
		{
			name:    "success - default retries",
			jobType: "test_job",
			payload: queue.JobPayload{
				"key": "value",
			},
			maxRetries: nil,
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("CreateJob", mock.Anything, mock.MatchedBy(func(params sqlc.CreateJobParams) bool {
					return params.Type == "test_job" &&
						params.MaxRetries == 3 &&
						len(params.Payload) > 0
				})).Return(createTestJob(), nil)
			},
			expectedError: nil,
			validateJob: func(t *testing.T, job *sqlc.Job) {
				assert.Equal(t, "test_job", job.Type)
				assert.Equal(t, int32(3), job.MaxRetries)
			},
		},
		{
			name:    "success - custom retries",
			jobType: "test_job",
			payload: queue.JobPayload{
				"key": "value",
			},
			maxRetries: intPtr(5),
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("CreateJob", mock.Anything, mock.MatchedBy(func(params sqlc.CreateJobParams) bool {
					return params.Type == "test_job" &&
						params.MaxRetries == 5
				})).Return(createTestJob(), nil)
			},
			expectedError: nil,
			validateJob: func(t *testing.T, job *sqlc.Job) {
				assert.Equal(t, "test_job", job.Type)
			},
		},
		{
			name:    "error - database error",
			jobType: "test_job",
			payload: queue.JobPayload{
				"key": "value",
			},
			maxRetries: nil,
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("CreateJob", mock.Anything, mock.Anything).
					Return(sqlc.Job{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateJob:   nil,
		},
		{
			name:    "error - invalid payload",
			jobType: "test_job",
			payload: queue.JobPayload{
				"key": make(chan int),
			},
			maxRetries:    nil,
			setupMocks:    func(ms *dbmocks.MockStore) {},
			expectedError: errors.New("json: unsupported type: chan int"),
			validateJob:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := dbmocks.NewMockStore(t)
			tt.setupMocks(mockStore)

			service := queue.NewService(mockStore)

			job, err := service.Enqueue(context.Background(), tt.jobType, tt.payload, tt.maxRetries)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, job)
			} else {
				require.NoError(t, err)
				require.NotNil(t, job)
				if tt.validateJob != nil {
					tt.validateJob(t, job)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_GetNextPendingJob(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*dbmocks.MockStore)
		expectedJob   *sqlc.Job
		expectedError error
	}{
		{
			name: "success - job found",
			setupMocks: func(ms *dbmocks.MockStore) {
				job := createTestJob()
				ms.On("GetNextPendingJob", mock.Anything).
					Return(job, nil)
			},
			expectedJob:   nil,
			expectedError: nil,
		},
		{
			name: "success - no job found",
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("GetNextPendingJob", mock.Anything).
					Return(sqlc.Job{}, pgx.ErrNoRows)
			},
			expectedJob:   nil,
			expectedError: nil,
		},
		{
			name: "error - database error",
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("GetNextPendingJob", mock.Anything).
					Return(sqlc.Job{}, errors.New("database error"))
			},
			expectedJob:   nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := dbmocks.NewMockStore(t)
			tt.setupMocks(mockStore)

			service := queue.NewService(mockStore)

			job, err := service.GetNextPendingJob(context.Background())

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, job)
			} else {
				require.NoError(t, err)
				if tt.expectedJob == nil {
					if job != nil {
						assert.NotEmpty(t, job.ID)
					}
				} else {
					require.NotNil(t, job)
					assert.Equal(t, tt.expectedJob.ID, job.ID)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_MarkJobCompleted(t *testing.T) {
	tests := []struct {
		name          string
		jobID         uuid.UUID
		setupMocks    func(*dbmocks.MockStore)
		expectedError error
	}{
		{
			name:  "success",
			jobID: uuid.New(),
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("UpdateJobStatus", mock.Anything, mock.MatchedBy(func(params sqlc.UpdateJobStatusParams) bool {
					return params.Status == "completed" &&
						params.ErrorMessage == nil
				})).Return(createTestJob(), nil)
			},
			expectedError: nil,
		},
		{
			name:  "error - database error",
			jobID: uuid.New(),
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("UpdateJobStatus", mock.Anything, mock.Anything).
					Return(sqlc.Job{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := dbmocks.NewMockStore(t)
			tt.setupMocks(mockStore)

			service := queue.NewService(mockStore)

			err := service.MarkJobCompleted(context.Background(), tt.jobID)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_MarkJobFailed(t *testing.T) {
	tests := []struct {
		name          string
		jobID         uuid.UUID
		errorMsg      string
		setupMocks    func(*dbmocks.MockStore)
		expectedError error
	}{
		{
			name:     "success",
			jobID:    uuid.New(),
			errorMsg: "test error",
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("UpdateJobStatus", mock.Anything, mock.MatchedBy(func(params sqlc.UpdateJobStatusParams) bool {
					return params.Status == "failed" &&
						params.ErrorMessage != nil &&
						*params.ErrorMessage == "test error"
				})).Return(createTestJob(), nil)
			},
			expectedError: nil,
		},
		{
			name:     "error - database error",
			jobID:    uuid.New(),
			errorMsg: "test error",
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("UpdateJobStatus", mock.Anything, mock.Anything).
					Return(sqlc.Job{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := dbmocks.NewMockStore(t)
			tt.setupMocks(mockStore)

			service := queue.NewService(mockStore)

			err := service.MarkJobFailed(context.Background(), tt.jobID, tt.errorMsg)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_RetryJob(t *testing.T) {
	tests := []struct {
		name          string
		jobID         uuid.UUID
		errorMsg      string
		setupMocks    func(*dbmocks.MockStore)
		expectedError error
	}{
		{
			name:     "success",
			jobID:    uuid.New(),
			errorMsg: "retry error",
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("IncrementJobRetry", mock.Anything, mock.MatchedBy(func(params sqlc.IncrementJobRetryParams) bool {
					return params.ErrorMessage != nil &&
						*params.ErrorMessage == "retry error"
				})).Return(createTestJob(), nil)
			},
			expectedError: nil,
		},
		{
			name:     "error - database error",
			jobID:    uuid.New(),
			errorMsg: "retry error",
			setupMocks: func(ms *dbmocks.MockStore) {
				ms.On("IncrementJobRetry", mock.Anything, mock.Anything).
					Return(sqlc.Job{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := dbmocks.NewMockStore(t)
			tt.setupMocks(mockStore)

			service := queue.NewService(mockStore)

			err := service.RetryJob(context.Background(), tt.jobID, tt.errorMsg)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func createTestJob() sqlc.Job {
	now := time.Now()
	payload, _ := json.Marshal(map[string]interface{}{"test": "data"})
	return sqlc.Job{
		ID:          uuid.New(),
		Type:        "test_job",
		Payload:     payload,
		Status:      "pending",
		RetryCount:  0,
		MaxRetries:  3,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func createTestJobPtr() *sqlc.Job {
	job := createTestJob()
	return &job
}

func createTestJobWithID(id uuid.UUID) sqlc.Job {
	job := createTestJob()
	job.ID = id
	return job
}

