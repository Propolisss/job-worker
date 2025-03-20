package repository

import (
  "context"
  "encoding/json"
  "fmt"
  "time"

  errs "flussonic_tz/internal/errors"

  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"

  "flussonic_tz/internal/service"
  "flussonic_tz/models"

  "github.com/go-redis/redis/v8"
)

const (
  StatusPending    = "pending"
  StatusInProgress = "in_progress"
  StatusCompleted  = "completed"
  StatusFailed     = "failed"
)

type RedisRepository struct {
  client    *redis.Client
  queueName string
}

func NewRedisRepository(client *redis.Client, queueName string) service.JobRepository {
  return &RedisRepository{
    client:    client,
    queueName: queueName,
  }
}

func (r *RedisRepository) AddJob(ctx context.Context, job *models.Job) error {
  jsonMsg, err := json.Marshal(job)
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrMarshalJob)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return wrapped
  }

  err = r.client.ZAdd(ctx, r.queueName, &redis.Z{
    Score:  job.Score,
    Member: jsonMsg,
  }).Err()
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrAddJob)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return wrapped
  }

  status := map[string]interface{}{
    "status":     StatusPending,
    "score":      job.Score,
    "created_at": time.Now().Format(time.RFC3339),
  }

  return r.client.HSet(ctx, fmt.Sprintf("task:%s", job.ID), status).Err()
}

func (r *RedisRepository) GetJob(ctx context.Context) (*models.Job, error) {
  result, err := r.client.ZPopMin(ctx, r.queueName, 1).Result()
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrGetJob)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return nil, wrapped
  }

  if len(result) == 0 {
    return nil, errors.New("Job not found")
  }

  jsonJob, ok := result[0].Member.(string)
  if !ok {
    wrapped := errors.Wrap(err, errs.ErrCastError)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return nil, wrapped
  }

  var job *models.Job
  if err = json.Unmarshal([]byte(jsonJob), &job); err != nil {
    wrapped := errors.Wrap(err, errs.ErrUnmarshalJob)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return nil, wrapped
  }

  err = r.client.HSet(ctx, fmt.Sprintf("task:%s", job.ID), map[string]interface{}{
    "status":     StatusInProgress,
    "started_at": time.Now().Format(time.RFC3339),
  }).Err()
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrUpdateJob)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return nil, wrapped
  }

  return job, nil
}

func (r *RedisRepository) CompleteJob(ctx context.Context, jobID string) error {
  return r.client.HSet(ctx, fmt.Sprintf("task:%s", jobID), map[string]interface{}{
    "status":      StatusCompleted,
    "finished_at": time.Now().Format(time.RFC3339),
  }).Err()
}

func (r *RedisRepository) FailJob(ctx context.Context, jobID string) error {
  return r.client.HSet(ctx, fmt.Sprintf("task:%s", jobID), map[string]interface{}{
    "status":      StatusFailed,
    "finished_at": time.Now().Format(time.RFC3339),
  }).Err()
}

func (r *RedisRepository) GetJobStatus(ctx context.Context, jobID string) (string, error) {
  res, err := r.client.HGetAll(ctx, fmt.Sprintf("task:%s", jobID)).Result()
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrGetJobStatus)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return "", wrapped
  }

  jsonResp, err := json.Marshal(res)
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrUnmarshalJobStatus)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return "", wrapped
  }

  return string(jsonResp), nil
}
