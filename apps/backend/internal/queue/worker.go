package queue

import (
	sqlc "circa/internal/db/sqlc/generated"
	"circa/internal/email"
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

type Worker struct {
	queueService *Service
	emailService email.EmailService
	stopChan     chan struct{}
}

func NewWorker(queueService *Service, emailService email.EmailService) *Worker {
	return &Worker{
		queueService: queueService,
		emailService: emailService,
		stopChan:     make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Info().Msg("Queue worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Queue worker stopping")
			return
		case <-w.stopChan:
			log.Info().Msg("Queue worker stopped")
			return
		case <-ticker.C:
			w.ProcessJobs(ctx)
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopChan)
}

func (w *Worker) ProcessJobs(ctx context.Context) {
	job, err := w.queueService.GetNextPendingJob(ctx)
	if err != nil {
		return
	}

	if job == nil {
		return
	}

	log.Info().
		Str("job_id", job.ID.String()).
		Str("job_type", job.Type).
		Msg("Processing job")

	switch job.Type {
	case "send_magic_link_email":
		w.handleSendMagicLinkEmail(ctx, job)
	default:
		log.Warn().
			Str("job_type", job.Type).
			Msg("Unknown job type")
		w.queueService.MarkJobFailed(ctx, job.ID, "Unknown job type")
	}
}

func (w *Worker) handleSendMagicLinkEmail(ctx context.Context, job *sqlc.Job) {
	var payload struct {
		Email        string `json:"email"`
		Name         string `json:"name"`
		MagicLinkURL string `json:"magic_link_url"`
	}

	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to unmarshal job payload")
		w.queueService.MarkJobFailed(ctx, job.ID, "Invalid payload format")
		return
	}

	if w.emailService == nil {
		log.Error().Str("job_id", job.ID.String()).Msg("Email service not available")
		w.queueService.MarkJobFailed(ctx, job.ID, "Email service not configured")
		return
	}

	if err := w.emailService.SendMagicLink(ctx, payload.Email, payload.Name, payload.MagicLinkURL); err != nil {
		log.Error().
			Err(err).
			Str("job_id", job.ID.String()).
			Int32("retry_count", job.RetryCount).
			Int32("max_retries", job.MaxRetries).
			Msg("Failed to send magic link email")

		if job.RetryCount < job.MaxRetries {
			if retryErr := w.queueService.RetryJob(ctx, job.ID, err.Error()); retryErr != nil {
				log.Error().Err(retryErr).Str("job_id", job.ID.String()).Msg("Failed to retry job")
			} else {
				log.Info().
					Str("job_id", job.ID.String()).
					Int32("retry_count", job.RetryCount+1).
					Msg("Job scheduled for retry")
			}
		} else {
			w.queueService.MarkJobFailed(ctx, job.ID, err.Error())
		}
		return
	}

	if err := w.queueService.MarkJobCompleted(ctx, job.ID); err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to mark job as completed")
		return
	}

	log.Info().
		Str("job_id", job.ID.String()).
		Str("email", payload.Email).
		Msg("Magic link email sent successfully")
}

