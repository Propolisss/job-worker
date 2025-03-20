package service

import (
  "context"
  "time"

  "flussonic_tz/models"
  "flussonic_tz/pkg/generator"

  "github.com/rs/zerolog/log"
)

type JobRepository interface {
  AddJob(ctx context.Context, job *models.Job) error
  GetJob(ctx context.Context) (*models.Job, error)
  CompleteJob(ctx context.Context, jobID string) error
  FailJob(ctx context.Context, jobID string) error
  GetJobStatus(ctx context.Context, jobID string) (string, error)
}

type JobService struct {
  repo JobRepository
}

func NewJobService(repo JobRepository) *JobService {
  return &JobService{
    repo: repo,
  }
}

func (svc *JobService) CreateJob(ctx context.Context, name string, score float64) (string, error) {
  id, err := generator.GenerateID(32)
  if err != nil {
    log.Error().Err(err).Msg(err.Error())
    return "", err
  }

  job := &models.Job{
    ID:        id,
    Name:      name,
    Score:     score,
    Status:    "pending",
    CreatedAt: time.Now(),
  }

  return id, svc.repo.AddJob(ctx, job)
}

func (svc *JobService) GetJob(ctx context.Context) (*models.Job, error) {
  job, err := svc.repo.GetJob(ctx)
  if err != nil {
    log.Error().Err(err).Msg(err.Error())
    return nil, err
  }

  return job, nil
}

func (svc *JobService) GetJobStatus(ctx context.Context, jobID string) (string, error) {
  status, err := svc.repo.GetJobStatus(ctx, jobID)
  if err != nil {
    log.Error().Err(err).Msg(err.Error())
    return "", err
  }

  return status, nil
}
