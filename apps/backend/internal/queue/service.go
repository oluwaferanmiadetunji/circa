package queue

import (
	"circa/internal/db"
	sqlc "circa/internal/db/sqlc/generated"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

type JobPayload map[string]interface{}

func (s *Service) Enqueue(ctx context.Context, jobType string, payload JobPayload, maxRetries *int) (*sqlc.Job, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal job payload")
		return nil, err
	}

	retries := 1
	if maxRetries != nil {
		retries = *maxRetries
	}

	params := sqlc.CreateJobParams{
		Type:        jobType,
		Payload:     payloadJSON,
		MaxRetries:  int32(retries),
		ScheduledAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	job, err := s.store.CreateJob(ctx, params)
	if err != nil {
		log.Error().Err(err).Str("job_type", jobType).Msg("Failed to enqueue job")
		return nil, err
	}

	log.Info().
		Str("job_id", job.ID.String()).
		Str("job_type", jobType).
		Msg("Job enqueued successfully")

	return &job, nil
}

func (s *Service) GetNextPendingJob(ctx context.Context) (*sqlc.Job, error) {
	job, err := s.store.GetNextPendingJob(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (s *Service) MarkJobCompleted(ctx context.Context, jobID uuid.UUID) error {
	var errorMsg *string = nil
	_, err := s.store.UpdateJobStatus(ctx, sqlc.UpdateJobStatusParams{
		Status:       "completed",
		ErrorMessage: errorMsg,
		ID:           jobID,
	})
	return err
}

func (s *Service) MarkJobFailed(ctx context.Context, jobID uuid.UUID, errorMsg string) error {
	_, err := s.store.UpdateJobStatus(ctx, sqlc.UpdateJobStatusParams{
		Status:       "failed",
		ErrorMessage: &errorMsg,
		ID:           jobID,
	})
	return err
}

func (s *Service) RetryJob(ctx context.Context, jobID uuid.UUID, errorMsg string) error {
	_, err := s.store.IncrementJobRetry(ctx, sqlc.IncrementJobRetryParams{
		ErrorMessage: &errorMsg,
		ID:           jobID,
	})
	return err
}
